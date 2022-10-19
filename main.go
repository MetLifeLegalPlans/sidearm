package main

import (
	"math/rand"
	"github.com/MetLifeLegalPlans/sidearm/cmd"
	"time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	cmd.Execute()
}
