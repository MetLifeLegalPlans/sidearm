package client

import (
	"net/http"

	"github.com/MetLifeLegalPlans/sidearm/channels"
	"github.com/MetLifeLegalPlans/sidearm/config"
	"github.com/MetLifeLegalPlans/sidearm/db/models"

	"github.com/pebbe/zmq4"
	"github.com/vmihailenco/msgpack/v5"
)

var reportQueue = make(chan *models.Response, 1024*8)

func report(resp *http.Response, duration int64, scenario *config.Scenario) {
	if sink == nil {
		return
	}

	r := &models.Response{
		StatusCode: 504,
		Method:     scenario.Method,
		URL:        scenario.URL,
		Duration:   duration,
	}

	if resp != nil {
		r.StatusCode = resp.StatusCode
	}

	reportQueue <- r
}

func reportWorker() {
	if sink == nil {
		return
	}

	for {
		select {
		case <-channels.Running:
			return
		case msg := <-reportQueue:
			packed, _ := msgpack.Marshal(msg)
			sink.SendBytes(packed, zmq4.DONTWAIT)
		}
	}
}
