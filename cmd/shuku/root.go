package main

import (
	"fmt"

	"github.com/takumines/shuku/cmd/shuku/compress"
	"github.com/takumines/shuku/cmd/shuku/version"

	"github.com/urfave/cli/v2"
)

// rootCmd returns the root command.
func rootCmd() *cli.App {
	helpCommand := &cli.Command{
		Name:      "help",
		Aliases:   []string{"h"},
		Usage:     "Shows a list of commands or help for one command",
		ArgsUsage: "[command]",
		Action: func(c *cli.Context) error {
			args := c.Args()
			if args.Present() {
				return cli.ShowCommandHelp(c, args.First())
			}
			return cli.ShowAppHelp(c)
		},
	}

	return &cli.App{
		Name:      "shuku",
		Usage:     "A CLI tool for compressing images.",
		UsageText: "shuku [command] [options] [arguments]",
		Commands: []*cli.Command{
			compress.Cmd(),
			version.Cmd(),
			helpCommand,
		},
		Suggest: true,
		CommandNotFound: func(context *cli.Context, command string) {
			fmt.Printf("Command not found: %q\n", command)
		},
		HelpName: "shuku",
	}
}
