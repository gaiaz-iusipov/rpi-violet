package main

import (
	_ "net/http/pprof"

	"github.com/gaiaz-iusipov/rpi-violet/commands"
)

func main() {
	commands.Execute()
}
