package main

import "sync-ethereum/cmd"

//go:generate wire ./...

func main() {
	cmd.Execute()
}
