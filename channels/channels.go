package channels

import (
	"os"
	"os/signal"
	"syscall"
)

var Running = make(chan bool)
var Interrupt = make(chan os.Signal, 1)

func init() {
	signal.Notify(Interrupt, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
}
