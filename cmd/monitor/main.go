package main

import "github.com/nandiheath/k8s-node-monitor/internal/watcher"

func main() {
	w := watcher.New()
	w.Start()
}
