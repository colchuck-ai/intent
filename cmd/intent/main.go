// Command intent is the single-binary CLI for operating an intent.yaml tree.
package main

import (
	"os"

	"github.com/colchuck-ai/intent/internal/cli"
)

func main() {
	os.Exit(cli.Execute())
}
