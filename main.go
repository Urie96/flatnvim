package main

import (
	_ "embed"
	"log"
	"os"
	"os/exec"
	"syscall"

	"github.com/neovim/go-client/nvim"
)

//go:embed remote.lua
var lua string

func main() {
	addr := os.Getenv("NVIM")
	if addr == "" {
		binary, err := exec.LookPath("nvim")
		if err != nil {
			log.Panicln(err)
		}

		args := append([]string{"nvim"}, stripNoWait(os.Args[1:])...)
		env := os.Environ()

		syscall.Exec(binary, args, env)
		return
	}

	args, noWait := splitNoWait(os.Args[1:])
	if len(args) == 0 {
		log.Panicln("no arguments given")
	}

	nv, err := nvim.Dial(addr)
	if err != nil {
		log.Panicf("unable to connect to parent nvim instance: %v\n", err)
	}
	defer nv.Close()

	ch := make(chan struct{})
	nv.RegisterHandler("stop", func(event string, args ...any) {
		close(ch)
	})

	panicOnError(err)
	panicOnError(nv.ExecLua(lua, nil, nv.ChannelID(), args, noWait))

	<-ch
}

func panicOnError(err error) {
	if err != nil {
		log.Panicf("error: %v\n", err)
	}
}

// stripNoWait drops the --no-wait/-n flags so they never reach a directly
// invoked nvim, which would reject them as unknown options.
func stripNoWait(args []string) []string {
	out := args[:0]
	for _, a := range args {
		if a == "--no-wait" || a == "-n" {
			continue
		}
		out = append(out, a)
	}
	return out
}

// splitNoWait separates the no-wait flags from the remaining arguments.
func splitNoWait(args []string) ([]string, bool) {
	noWait := false
	var out []string
	for _, a := range args {
		if a == "--no-wait" || a == "-n" {
			noWait = true
			continue
		}
		out = append(out, a)
	}
	return out, noWait
}

