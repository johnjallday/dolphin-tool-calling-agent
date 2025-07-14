package agent

import (
    "context"
    "fmt"
    "path/filepath"
    "plugin"

    "github.com/BurntSushi/toml"
    "github.com/openai/openai-go"
    "github.com/johnjallday/dolphin-tool-calling-agent/internal/registry"
    "github.com/johnjallday/dolphin-tool-calling-agent/pkg/tools"
)

// Agent holds config and runtime state.
type Agent struct {
    Name     string   `toml:"name"`
    Model    string   `toml:"model"`
		ToolPaths []string `toml:"tool_path"`
    Registry *registry.ToolRegistry
    client   *openai.Client
    params   openai.ChatCompletionNewParams
}

// NewAgentFromConfig loads a TOML config, registers tools, and returns an Agent.
func NewAgentFromConfig(client *openai.Client, configPath string) (*Agent, error) {
    var a Agent
    absConfig, err := filepath.Abs(configPath)
    if err != nil {
        return nil, fmt.Errorf("resolve config path: %w", err)
    }
    if _, err := toml.DecodeFile(absConfig, &a); err != nil {
        return nil, fmt.Errorf("decode config: %w", err)
    }

    a.client = client
    a.params = openai.ChatCompletionNewParams{
        Messages:    []openai.ChatCompletionMessageParamUnion{},
        Model:       a.Model,
        Temperature: openai.Float(0),
        Seed:        openai.Int(0),
    }

    a.Registry = registry.NewToolRegistry()
    for _, tp := range a.ToolPaths {
        absP, err := filepath.Abs(tp)
        if err != nil {
            return nil, fmt.Errorf("resolve plugin path %q: %w", tp, err)
        }
        plug, err := plugin.Open(absP)
        if err != nil {
            return nil, fmt.Errorf("open plugin %q: %w", absP, err)
        }

        if symPkg, err := plug.Lookup("PluginPackage"); err == nil {
            pkgFunc, ok := symPkg.(func() tools.ToolPackage)
            if !ok {
                return nil, fmt.Errorf("invalid PluginPackage signature in %q", absP)
            }
            for _, t := range pkgFunc().Tools {
                a.Registry.Register(t)
            }
            continue
        }
        if symSpecs, err := plug.Lookup("PluginSpecs"); err == nil {
            if ptr, ok := symSpecs.(*[]tools.Tool); ok {
                for _, t := range *ptr {
                    a.Registry.Register(t)
                }
                continue
            }
            if fn, ok := symSpecs.(func() []tools.Tool); ok {
                for _, t := range fn() {
                    a.Registry.Register(t)
                }
                continue
            }
            return nil, fmt.Errorf("invalid PluginSpecs in %q", absP)
        }
        return nil, fmt.Errorf("no PluginPackage or PluginSpecs in %q", absP)
    }

    a.Registry.Initialize(&a.params)
    return &a, nil
}

// SendMessage sends a user message, dispatches tool calls, and prints responses.
func (a *Agent) SendMessage(ctx context.Context, userMessage string) error {
    a.params.Messages = append(a.params.Messages, openai.UserMessage(userMessage))
    cmp, err := a.client.Chat.Completions.New(ctx, a.params)
    if err != nil {
        return err
    }
    assistant := cmp.Choices[0].Message
    a.params.Messages = append(a.params.Messages, assistant.ToParam())

    if len(assistant.ToolCalls) == 0 {
        fmt.Println(assistant.Content)
        return nil
    }
    //a.dispatchTools(assistant.ToolCalls, &a.params)


		a.dispatchTools(assistant.ToolCalls)
  
  

    finalResp, err := a.client.Chat.Completions.New(ctx, a.params)
    if err != nil {
        return err
    }
    finalMsg := finalResp.Choices[0].Message
    fmt.Println(finalMsg.Content)
    a.params.Messages = append(a.params.Messages, finalMsg.ToParam())
    return nil
}

func (a *Agent) dispatchTools(toolCalls []openai.ChatCompletionMessageToolCall) {
     for _, tc := range toolCalls {
       if h, ok := a.Registry.Handlers()[tc.Function.Name]; ok {
           h(tc, &a.params)
       }
     }
}

func (a *Agent) PrintTools() {
    a.Registry.PrintTools()
}
