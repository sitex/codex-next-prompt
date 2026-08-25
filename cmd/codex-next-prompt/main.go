package main

import (
	"os"

	"github.com/sitex/codex-next-prompt/internal/hook"
)

func main() {
	command := ""
	if len(os.Args) == 2 {
		command = os.Args[1]
	}

	os.Exit(hook.Run(command, hook.Streams{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}))
}
