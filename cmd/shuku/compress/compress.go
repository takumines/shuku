package compress

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

// Cmd returns the compress command.
func Cmd() *cli.Command {
	return &cli.Command{
		Name:    "compress",
		Usage:   "Compress an image.",
		Aliases: []string{"c"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "input",
				Aliases:  []string{"i"},
				Usage:    "Input image file path",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "Specify the output file name (optional, defaults to input file name + '_compressed').",
			},
			&cli.IntFlag{
				Name:    "quality",
				Aliases: []string{"q"},
				Usage:   "JPEG/WebP quality (0-100)",
				Value:   80,
			},
		},
		Action: compressAction,
	}
}

// compressAction is the action for the compress command.
func compressAction(c *cli.Context) error {
	// TODO: 圧縮処理の実装をする
	fmt.Println("Compressing image...")
	return nil
}
