package client

import (
	"fmt"
	"net/http"
	"time"

	"github.com/MetLifeLegalPlans/sidearm/channels"
	"github.com/MetLifeLegalPlans/sidearm/config"

	. "github.com/logrusorgru/aurora/v3"
	"github.com/pebbe/zmq4"
)

var Client = &http.Client{
	Transport: &http.Transport{
		DisableKeepAlives: true,
	},
}
var sink *zmq4.Socket

func init() {
	Client.Timeout = 5 * time.Second
}

func logResult(code int, method, url string, color func(arg any) Value) {
	fmt.Println(
		fmt.Sprintf(
			"[%v] %v %v",
			Bold(color(code)),
			Bold(method),
			Cyan(url),
		),
	)
}

func Entrypoint(conf *config.Config, quiet bool) {
	receiver, err := zmq4.NewSocket(zmq4.PULL)
	defer receiver.Close()
	if err != nil {
		panic(err)
	}

	err = receiver.Connect(conf.QueueConfig.Connect)
	if err != nil {
		panic(err)
	}

	if conf.SinkConfig.Enabled() {
		sink, err = zmq4.NewSocket(zmq4.PUSH)
		if err != nil {
			panic(err)
		}
		defer sink.Close()
		sink.Bind(conf.SinkConfig.Bind)

		go reportWorker()
	}

	for {
		select {
		case <-channels.Running:
			return
		default:
			msg, _ := receiver.RecvBytes(0)
			go process(msg, quiet)
		}
	}
}
