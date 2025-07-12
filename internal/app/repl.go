package app

import (
    "context"
  	"github.com/johnjallday/dolphin-tool-calling-agent/internal/user"
    openai "github.com/openai/openai-go"
)

type REPLApp struct {
    auth   Authenticator
    userL  UserLoader
    agentM AgentManager

    usr     *user.User
    agent   string
    client  *openai.Client
}

func NewREPLApp(
    auth Authenticator,
    userL UserLoader,
    agentM AgentManager,
) *REPLApp {
    return &REPLApp{
        auth:   auth,
        userL:  userL,
        agentM: agentM,
    }
}

func (a *REPLApp) Init(ctx context.Context) error {
    // 1) Login
    userID, err := a.auth.Login()
    if err != nil {
        return err
    }

    // 2) Load user (auto‐creates config dirs/files if needed)
    u, err := a.userL.Load(userID)
    if err != nil {
        return err
    }
    a.usr = u

    // 3) Pick or confirm default agent
    agentName, err := a.agentM.Select(ctx, u)
    if err != nil {
        return err
    }
    a.agent = agentName

    // 4) Load the agent configuration
    if err := a.agentM.Load(ctx, u, agentName); err != nil {
        return err
    }

    // 5) Wire up OpenAI (or whatever backend)
    client := openai.NewClient()
    a.client = &client

    return nil
}

func (a *REPLApp) Run(ctx context.Context) error {
    // now you can drive your REPL loop, calling a.client, etc.
    return nil
}
