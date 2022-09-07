package dashboard

import (
	"fmt"
	"sidearm/channels"
	"sidearm/config"

	"github.com/pebbe/zmq4"
)

func Entrypoint(conf *config.Config) {
	receiver, err := zmq4.NewSocket(zmq4.PULL)
	if err != nil {
		panic(err)
	}
	defer receiver.Close()

	err = receiver.Bind(conf.SinkConfig.Bind)
	if err != nil {
		panic(err)
	}

	go ui()

	for {
		select {
		case <-channels.Running:
			return
		default:
			msg, _ := receiver.RecvBytes(0)
			fmt.Println(msg)
		}
	}
}
