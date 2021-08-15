package main

import (
	"math/rand"
	_ "net/http/pprof"
	"time"

	"github.com/gaiaz-iusipov/rpi-violet/commands"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	commands.Execute()
}
