package main

import "github.com/fitzix/spider/cmd"

var (
	version = "dev"
)

func main() {
	cmd.Execute(version)
}
