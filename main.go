package main

import (
	_ "embed"
	"log"
	"os"
	"os/exec"

	"github.com/neovim/go-client/nvim"
)

//go:embed remote.lua
var lua string

func main() {
	addr := os.Getenv("NVIM")
	if addr == "" {
		cmd := exec.Command("nvim", os.Args[1:]...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			log.Panicln(err)
		}

		return
	}

	files := os.Args[1:]
	if len(files) == 0 {
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
	panicOnError(nv.ExecLua(lua, nil, nv.ChannelID(), files))

	<-ch
}

func panicOnError(err error) {
	if err != nil {
		log.Panicf("error: %v\n", err)
	}
}
