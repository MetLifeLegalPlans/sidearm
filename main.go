package main

import (
	"math/rand"
	"sidearm/cmd"
	"time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	cmd.Execute()
}
