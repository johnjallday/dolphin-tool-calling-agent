package agent

import (
	"os"
  "context"
  "fmt"
  "path/filepath"
  "plugin"
	"io/fs"
	"errors"
	"encoding/json"

  "github.com/openai/openai-go"
  "github.com/johnjallday/dolphin-tool-calling-agent/internal/registry"
  "github.com/johnjallday/dolphin-tool-calling-agent/pkg/tools"
)

type Agent struct {
  Name     string
  Model    string
  Registry *registry.ToolRegistry
  client   openai.Client
	history []ChatMessage
  params   openai.ChatCompletionNewParams
	systemPrompt openai.ChatCompletionMessageParamUnion
}

type ChatMessage struct {
  Role    string // "user" or "assistant"
  Content string
}

func NewAgent(name, model string, pluginNames []string) (*Agent, error) {
  client := openai.NewClient()

  // define your system prompt once, up front
	const sysText = `You are only allowed to respond by invoking one of the available functions.
	You must never return plain text directly.
	If you can't call any tools just say you don't have the necessary tools to execute.`

	sys := openai.SystemMessage(sysText)
  // seed params.Messages with the system prompt
  params := openai.ChatCompletionNewParams{
    Messages:    []openai.ChatCompletionMessageParamUnion{sys},
    Model:       model,
    Temperature: openai.Float(0),
    Seed:        openai.Int(0),
  }

  a := &Agent{
    Name:         name,
    Model:        model,
		history:      []ChatMessage{
			{"system", sysText},
		},
    client:       client,
    Registry:     registry.NewToolRegistry(),
    params:       params,
    systemPrompt: sys,
  }


  // load plugins (recursive search logic omitted for brevity)
  cwd, _ := os.Getwd()
  pluginDir := filepath.Join(cwd, "plugins")
  for _, pname := range pluginNames {
    soPath, err := locatePlugin(pluginDir, pname+".so")
    if err != nil {
      return nil, err
    }
    plug, err := plugin.Open(soPath)
    if err != nil {
      return nil, fmt.Errorf("open plugin %q: %w", pname, err)
    }
    sym, err := plug.Lookup("PluginPackage")
    if err != nil {
      return nil, fmt.Errorf("lookup PluginPackage in %q: %w", pname, err)
    }
    pkgFunc, ok := sym.(func() tools.ToolPackage)
    if !ok {
      return nil, fmt.Errorf("invalid PluginPackage signature in %q", pname)
    }
    pkg := pkgFunc()
    for _, t := range pkg.Tools {
      a.Registry.Register(t)
    }
  }

  a.Registry.Initialize(&a.params)
  return a, nil
}

// findPluginFile walks pluginDir looking for a file named pluginFileName.
// It returns the first match or an error if none are found.
func locatePlugin(pluginDir, pluginFileName string) (string, error) {
  var (
    foundPath    string
    sentinelErr  = errors.New("found")
  )

  err := filepath.WalkDir(pluginDir, func(path string, d fs.DirEntry, err error) error {
    if err != nil {
      return err
    }
    if !d.IsDir() && d.Name() == pluginFileName {
      foundPath = path
      return sentinelErr
    }
    return nil
  })

  // If we bailed out early with sentinelErr, clear it
  if err == sentinelErr {
    err = nil
  }
  if err != nil {
    return "", fmt.Errorf("walking %q: %w", pluginDir, err)
  }
  if foundPath == "" {
    return "", fmt.Errorf("plugin %q not found under %s", pluginFileName, pluginDir)
  }
  return foundPath, nil
}

func (a *Agent) SendMessage(ctx context.Context, userMessage string) (reply string, err error) {
  // 1) append the user message
  a.params.Messages = append(a.params.Messages, openai.UserMessage(userMessage))
  a.history = append(a.history, ChatMessage{"user", userMessage})

  // 2) first LLM call
  cmp, err := a.client.Chat.Completions.New(ctx, a.params)
  if err != nil {
    return "", err
  }
  assistant := cmp.Choices[0].Message

  // 3) record assistant’s reply (and any tooling)
  a.params.Messages = append(a.params.Messages, assistant.ToParam())
  if assistant.Content != "" {
    a.history = append(a.history, ChatMessage{"assistant", assistant.Content})
  }

  // 4) if there are no tool calls, just return the content
  if len(assistant.ToolCalls) == 0 {
    return assistant.Content, nil
  }

  // 5) otherwise perform the tool calls
	a.dispatchTools(assistant.ToolCalls)

  // 6) final LLM call after tools
  finalResp, err := a.client.Chat.Completions.New(ctx, a.params)
  if err != nil {
    return "", err
  }
  finalMsg := finalResp.Choices[0].Message

  // 7) record & return final content
  a.params.Messages = append(a.params.Messages, finalMsg.ToParam())
  if finalMsg.Content != "" {
    a.history = append(a.history, ChatMessage{"assistant", finalMsg.Content})
  }
  return finalMsg.Content, nil
}

func (a *Agent) dispatchTools(toolCalls []openai.ChatCompletionMessageToolCall) {
  for _, tc := range toolCalls {
    if h, ok := a.Registry.Handlers()[tc.Function.Name]; ok {
      h(tc, &a.params)
    }
  }
}

func (a *Agent) Tools() []tools.Tool {
    return a.Registry.Tools()
}



func (a *Agent) Close() {
  a.Name = ""
  a.Model = ""
  a.params = openai.ChatCompletionNewParams{}
  a.Registry = nil
}

func (a *Agent) String() string {
	if a == nil || a.Name == "<none>" {
		return "No agent selected\n"
	}
	result := fmt.Sprintf("Agent: %s\nModel: %s\n", a.Name, a.Model)
	result += a.Registry.String()
	return result
}



// DumpMessages will pretty-print your prompt slice
func (a *Agent) DumpMessages() {
  b, err := json.MarshalIndent(a.params.Messages, "", "  ")
  if err != nil {
    fmt.Println("❌ failed to marshal messages:", err)
    return
  }
  fmt.Println("=== PARAMS.MESSAGES ===")
  fmt.Println(string(b))
  fmt.Println("=======================")
}

func (a *Agent) History() []ChatMessage {
  // copy to prevent callers mutating your internal slice
  out := make([]ChatMessage, len(a.history))
  copy(out, a.history)
  return out
}
