package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/johnjallday/dolphin-tool-calling-agent/agent"
	"github.com/urfave/cli/v3"
)

func printLogo() {
	logo := `
		🐬
`
	fmt.Print("\033[36m" + logo + "\033[0m\n")
	fmt.Println("Dolphin Tool Calling CLI")
}

func main() {
	printLogo()
	cmd := &cli.Command{
		Commands: []*cli.Command{
			{
				Name:    "agents",
				Aliases: []string{"a"},
				Usage:   "list available agents",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					agent.ListAgents()
					return nil
				},
				Commands: []*cli.Command{
					{
						Name:  "reaper",
						Usage: "view my agent in detail",
						Action: func(ctx context.Context, cmd *cli.Command) error {
							fmt.Println("test")
							return nil
						},
					},
				},
			},
			{
				Name:    "tools",
				Aliases: []string{"t"},
				Usage:   "get list of tools",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					fmt.Println("completed task: ", cmd.Args().First())
					return nil
				},
			},
			{
				Name:    "agent_builder",
				Aliases: []string{"ab"},
				Usage:   "Build your own Agent",
				Commands: []*cli.Command{
					{
						Name:  "create",
						Usage: "Use AI to build your own agent",
						Action: func(ctx context.Context, cmd *cli.Command) error {
							fmt.Println("agent built")
							return nil
						},
					},
					{
						Name:  "delete",
						Usage: "delete existing AI agent",
						Action: func(ctx context.Context, cmd *cli.Command) error {
							//fmt.Println("removed task template: ", cmd.Args().First())
							fmt.Println("Agent Removed")
							return nil
						},
					},
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
