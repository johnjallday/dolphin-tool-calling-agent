package agent

import (
  "context"
  "fmt"
  "path/filepath"
  "plugin"

  "github.com/openai/openai-go"
  "github.com/johnjallday/dolphin-tool-calling-agent/internal/registry"
  "github.com/johnjallday/dolphin-tool-calling-agent/pkg/tools"
)

type Agent struct {
  Name     string
  Model    string
  Registry *registry.ToolRegistry
  client   *openai.Client
  params   openai.ChatCompletionNewParams
}

func NewAgent(client *openai.Client, name, model string, toolPaths []string) (*Agent, error) {
  a := &Agent{
    Name:   name,
    Model:  model,
    client: client,
    Registry: registry.NewToolRegistry(),
    params: openai.ChatCompletionNewParams{
      Messages:    []openai.ChatCompletionMessageParamUnion{},
      Model:       model,
      Temperature: openai.Float(0),
      Seed:        openai.Int(0),
    },
  }

  for _, tp := range toolPaths {
    absP, err := filepath.Abs(tp)
    if err != nil {
      return nil, fmt.Errorf("resolve plugin path %q: %w", tp, err)
    }
    plug, err := plugin.Open(absP)
    if err != nil {
      return nil, fmt.Errorf("open plugin %q: %w", absP, err)
    }
    // same lookup logic as before...
    // register tools into a.Registry
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
  }

  a.Registry.Initialize(&a.params)
  return a, nil
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

func (a *Agent) Close() {
    a.Name = ""
    a.Model = ""
    a.params = openai.ChatCompletionNewParams{}
    a.Registry = nil
    a.client = nil
}
