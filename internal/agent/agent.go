package agent

import (
	"os"
  "context"
  "fmt"
  "path/filepath"
  "plugin"
	"sort"

  "github.com/fatih/color"

  "github.com/openai/openai-go"
  "github.com/johnjallday/dolphin-tool-calling-agent/internal/registry"
  "github.com/johnjallday/dolphin-tool-calling-agent/pkg/tools"
)

type Agent struct {
  Name     string
  Model    string
  Registry *registry.ToolRegistry
  client   openai.Client
  params   openai.ChatCompletionNewParams
}

func NewAgent(name, model string, pluginNames []string) (*Agent, error) {
  client := openai.NewClient()
  a := &Agent{
    Name:     name,
    Model:    model,
    client:   client,
    Registry: registry.NewToolRegistry(),
    params: openai.ChatCompletionNewParams{
      Messages:    []openai.ChatCompletionMessageParamUnion{},
      Model:       model,
      Temperature: openai.Float(0),
      Seed:        openai.Int(0),
    },
  }

  cwd, _ := os.Getwd()
  pluginDir := filepath.Join(cwd, "plugins")

  for _, pname := range pluginNames {
    path := filepath.Join(pluginDir, pname+".so")
    plug, err := plugin.Open(path)
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
}

func (a *Agent) Print() {
	if a.Name == "<none>" {
    fmt.Println("No agent selected")
    return
  }
  cLabel := color.New(color.FgCyan, color.Bold)
  cValue := color.New(color.FgWhite)
  cList  := color.New(color.FgMagenta, color.Bold)
  cItem  := color.New(color.FgGreen)

  cLabel.Print("Agent: "); cValue.Println(a.Name)
  cLabel.Print("Model: "); cValue.Println(a.Model)

  cList.Println("Registered Tools:")
  handlers := a.Registry.Handlers()
  names := make([]string, 0, len(handlers))
  for name := range handlers {
    names = append(names, name)
  }
  sort.Strings(names)
  for _, name := range names {
    cItem.Println("  - " + name)
  }
}
