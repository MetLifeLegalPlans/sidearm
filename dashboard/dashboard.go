package dashboard

import (
	"sidearm/channels"
	"sidearm/config"
	"sidearm/dashboard/ui"
	"sidearm/db"
	"sidearm/db/models"

	"github.com/pebbe/zmq4"
	"github.com/vmihailenco/msgpack/v5"
)

func Entrypoint(conf *config.Config) {
	receiver, err := zmq4.NewSocket(zmq4.PULL)
	defer receiver.Close()
	if err != nil {
		panic(err)
	}

	err = receiver.Connect(conf.SinkConfig.Connect)
	if err != nil {
		panic(err)
	}

	db.Setup(conf)

	go ui.Run(conf)
	go eventReceiver(conf)

	for {
		select {
		case <-channels.Running:
			return
		default:
			msg, _ := receiver.RecvBytes(0)
			record := models.Response{}
			msgpack.Unmarshal(msg, &record)
			resultQueue <- record
		}
	}
}
