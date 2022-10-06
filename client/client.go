package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"sidearm/channels"
	"sidearm/config"
	"sidearm/dashboard"

	. "github.com/logrusorgru/aurora/v3"
	"github.com/pebbe/zmq4"
	"github.com/vmihailenco/msgpack/v5"
)

var Client = &http.Client{}
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

func report(resp *http.Response, duration int64, scenario *config.Scenario) {
	if sink == nil {
		return
	}

	r := &dashboard.Result{
		StatusCode: resp.StatusCode,
		Duration:   duration,
		Scenario:   *scenario,
	}

	msg, _ := msgpack.Marshal(r)
	sink.SendBytes(msg, zmq4.DONTWAIT)
}

func process(msg []byte, quiet bool) {
	data := &config.Scenario{}
	msgpack.Unmarshal(msg, data)

	var (
		req  *http.Request
		resp *http.Response
		err  error
	)

	switch data.Method {
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		requestBody, _ := json.Marshal(data.Body)
		req, _ = http.NewRequest(data.Method, data.URL, bytes.NewBuffer(requestBody))
		req.Header.Add("Content-Type", "application/json")
	default:
		req, _ = http.NewRequest(data.Method, data.URL, nil)
	}

	start := time.Now()
	resp, err = Client.Do(req)
	if err != nil && !quiet {
		logResult(http.StatusRequestTimeout, data.Method, data.URL, Red)
	}
	defer resp.Body.Close()
	elapsed := time.Since(start).Milliseconds()

	report(resp, elapsed, data)

	if quiet || err != nil {
		return
	}

	statusCodeColor := Green
	if resp.StatusCode >= 400 || err != nil {
		statusCodeColor = Red
	}
	logResult(resp.StatusCode, data.Method, data.URL, statusCodeColor)
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
