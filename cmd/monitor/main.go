package main

import (
	"github.com/alecthomas/kong"
	"github.com/nandiheath/k8s-node-monitor/internal/cmd"
)

func main() {

	c := cmd.Command{}
	ctx := kong.Parse(&c)
	ctx.Run()

}
