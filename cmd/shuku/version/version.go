package version

import (
	"github.com/urfave/cli/v2"
	"log"
)

// Cmd returns the version command.
func Cmd() *cli.Command {
	return &cli.Command{
		Name:    "version",
		Usage:   "Print the version",
		Aliases: []string{"v"},
		Action: func(c *cli.Context) error {
			log.Println("shuku version 1.0.0")
			return nil
		},
	}
}
