package main

import (
	"os"
)

func main() {
	if err := rootCmd().Run(os.Args); err != nil {
		os.Exit(1)
	}
}
