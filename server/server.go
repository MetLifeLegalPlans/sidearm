package server

import (
	"sidearm/channels"
	"sidearm/config"
	"time"

	"github.com/pebbe/zmq4"
)

func Entrypoint(conf *config.Config) {
	distributor, err := zmq4.NewSocket(zmq4.PUSH)
	if err != nil {
		panic(err)
	}
	defer distributor.Close()

	err = distributor.Bind(conf.QueueConfig.Bind)
	if err != nil {
		panic(err)
	}

	// Send out requests with a minimum of a second between them
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	context := NewContext(conf)

	for {
		select {
		case <-channels.Running:
			return
		case <-ticker.C:
			numRequests := context.RequestsForSecond()

			// Perform a non-blocking send across the socket
			for i := 0; i <= int(numRequests); i++ {
				distributor.SendBytes(context.Choose(), zmq4.DONTWAIT)
			}
		}
	}
}
