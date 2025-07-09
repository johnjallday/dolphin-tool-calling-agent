package agent

import (
   "context"
   "fmt"
   "path/filepath"
   "plugin"

   "github.com/BurntSushi/toml"
   "github.com/openai/openai-go"
   "github.com/johnjallday/dolphin-tool-calling-agent/registry"
   "github.com/johnjallday/dolphin-tool-calling-agent/tools"
)

// Agent defines the methods any agent must implement. 
type Agent interface { 
	// SendMessage sends a user message and processes the conversation. 
	SendMessage(ctx context.Context, userMessage string) error 
}

// AgentConfig represents settings for creating an agent from a TOML file.
type AgentConfig struct {
   Name        string   `toml:"name"`
   Model       string   `toml:"model"`
   ToolPaths   []string `toml:"tool_path"`
   PluginPaths []string `toml:"plugin_paths"`
}

// DefaultAgent is a concrete implementation of Agent. 
type DefaultAgent struct { 
	client *openai.Client 
	params openai.ChatCompletionNewParams 
}

// NewAgent creates and returns a new DefaultAgent. 
func NewAgent(client *openai.Client, model string) Agent { 
	params := openai.ChatCompletionNewParams{ 
		Messages: []openai.ChatCompletionMessageParamUnion{}, 
		Model: model, 
		Temperature: openai.Float(0), 
		Seed: openai.Int(0), 
	}

	registry.Initialize(&params)

	return &DefaultAgent{
		client: client,
		params: params,
	}
}
// NewAgentFromConfig loads a TOML config, registers tools, and returns a configured Agent.
func NewAgentFromConfig(client *openai.Client, configPath string) (Agent, error) {
   var cfg AgentConfig
   absPath, err := filepath.Abs(configPath)
   if err != nil {
       return nil, fmt.Errorf("unable to resolve config path: %w", err)
   }
   if _, err := toml.DecodeFile(absPath, &cfg); err != nil {
       return nil, fmt.Errorf("failed to decode config file: %w", err)
   }
   params := openai.ChatCompletionNewParams{
       Messages:    []openai.ChatCompletionMessageParamUnion{},
       Model:       cfg.Model,
       Temperature: openai.Float(0),
       Seed:        openai.Int(0),
   }

	for _, pp := range cfg.PluginPaths {
			absP, err := filepath.Abs(pp)
			if err != nil {
					return nil, fmt.Errorf("resolve plugin path %q: %w", pp, err)
			}
			plug, err := plugin.Open(absP)
			if err != nil {
					return nil, fmt.Errorf("open plugin %q: %w", absP, err)
			}

			// Try PluginPackage first
			if symPkg, err := plug.Lookup("PluginPackage"); err == nil {
					pkgFunc, ok := symPkg.(func() tools.ToolPackage)
					if !ok {
							return nil, fmt.Errorf("invalid PluginPackage signature in %q", absP)
					}
					pkg := pkgFunc()
					fmt.Printf("Loaded tool package: name=%s version=%s link=%s\n", pkg.Name, pkg.Version, pkg.Link)
					for _, spec := range pkg.Specs {
							registry.RegisterSpec(spec)
					}
					continue
			}

			// Fallback to old PluginSpecs
			symSpecs, err := plug.Lookup("PluginSpecs")
			if err != nil {
					return nil, fmt.Errorf("no PluginPackage or PluginSpecs in %q", absP)
			}
			specsFunc, ok := symSpecs.(func() []tools.ToolSpec)
			if !ok {
					return nil, fmt.Errorf("invalid PluginSpecs signature in %q", absP)
			}
			for _, spec := range specsFunc() {
					registry.RegisterSpec(spec)
			}
	}

   // Dynamically load Go plugins for additional tools
   registry.Initialize(&params)
   return &DefaultAgent{client: client, params: params}, nil
}

// SendMessage appends the user message, processes the chat response, 
// dispatches tool calls if any, and appends the final response. 
func (a *DefaultAgent) SendMessage(ctx context.Context, userMessage string) error { 
	// Append the user’s message to the conversation. 
	a.params.Messages = append(a.params.Messages, openai.UserMessage(userMessage))

	// Get the assistant's response.
	cmp, err := a.client.Chat.Completions.New(ctx, a.params)
	if err != nil {
		return err
	}

	assistantMsg := cmp.Choices[0].Message
	a.params.Messages = append(a.params.Messages, assistantMsg.ToParam())

   // If there are no tool calls, print the assistant’s response and exit.
   if len(assistantMsg.ToolCalls) == 0 {
       fmt.Println(assistantMsg.Content)
       return nil
   }
	// Dispatch tool calls and update the conversation with their responses.
	dispatchTools(assistantMsg.ToolCalls, &a.params)
	// After tool execution, get the final assistant response and print it.
	finalResp, err := a.client.Chat.Completions.New(ctx, a.params)
	if err != nil {
		return err
	}
	finalMsg := finalResp.Choices[0].Message
	fmt.Println(finalMsg.Content)
	// Append final assistant message to the conversation history.
	a.params.Messages = append(a.params.Messages, finalMsg.ToParam())
	return nil
}

// SendPromptAndReceiveToolCalls sends a prompt and returns the message along with any tool calls. 
func SendPromptAndReceiveToolCalls(ctx context.Context, client *openai.Client, params *openai.ChatCompletionNewParams) (openai.ChatCompletionMessage, []openai.ChatCompletionMessageToolCall) { 
	cmp, err := client.Chat.Completions.New(ctx, *params) 
	if err != nil { 
		panic(err) 
	} 
	msg := cmp.Choices[0].Message 
	return msg, msg.ToolCalls 
}

// dispatchTools processes any tool calls by dispatching them to the registered handlers. 
func dispatchTools(toolCalls []openai.ChatCompletionMessageToolCall, params *openai.ChatCompletionNewParams) { 
	handlers := registry.Handlers() 
	for _, tc := range toolCalls { 
		//fmt.Println(tc)
		if h, ok := handlers[tc.Function.Name]; ok { 
			h(tc, params) 
		} 
	} 
}

// runFinalChat executes the final chat call and prints the assistant’s response. 
func runFinalChat(ctx context.Context, client *openai.Client, params *openai.ChatCompletionNewParams) { 
	final, err := client.Chat.Completions.New(ctx, *params) 
	if err != nil { 
		panic(err) 
	} 
	fmt.Println(final.Choices[0].Message.Content) 
}
