package cmd

import "github.com/nandiheath/k8s-node-monitor/internal/watcher"

type Command struct {
	Watcher Watcher `cmd:"" help:"start the node monitor and start update SRV records"`
}

type Watcher struct {
}

func (cli *Watcher) Run() error {
	w := watcher.New()
	w.Start()

	return nil
}
