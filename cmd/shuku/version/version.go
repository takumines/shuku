package version

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

// ビルド時に設定される変数
var (
	Version = "dev"
	Commit  = "unknown"
	Date    = "unknown"
)

// Cmd returns the version command.
func Cmd() *cli.Command {
	return &cli.Command{
		Name:    "version",
		Usage:   "Print the version",
		Aliases: []string{"v"},
		Action: func(c *cli.Context) error {
			fmt.Printf("shuku version %s\n", Version)
			fmt.Printf("Commit: %s\n", Commit)
			fmt.Printf("Built: %s\n", Date)
			return nil
		},
	}
}
