package main

import (
	"github.com/urfave/cli/v2"
	"shuku/cmd/shuku/compress"
	"shuku/cmd/shuku/version"
)

// rootCmd returns the root command.
func rootCmd() *cli.App {

	return &cli.App{
		Name:      "shuku",
		Usage:     "A CLI tool for compressing images.",
		UsageText: "shuku [command] [options] [arguments]",
		Commands: []*cli.Command{
			compress.Cmd(),
			version.Cmd(),
		},
		Suggest: true,
	}
}
