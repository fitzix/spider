package main

import "os/exec"

func main() {
	exec.Command(`set_time -local `)
}
