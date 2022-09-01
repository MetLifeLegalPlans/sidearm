package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"sidearm/channels"
	"sidearm/config"

	. "github.com/logrusorgru/aurora/v3"
	"github.com/pebbe/zmq4"
	"github.com/vmihailenco/msgpack/v5"
)

var Client = &http.Client{}

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

func process(msg []byte) {
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

	resp, err = Client.Do(req)
	if err != nil {
		logResult(http.StatusRequestTimeout, data.Method, data.URL, Red)
		return
	}
	defer resp.Body.Close()

	statusCodeColor := Green
	if resp.StatusCode >= 400 || err != nil {
		statusCodeColor = Red
	}
	logResult(resp.StatusCode, data.Method, data.URL, statusCodeColor)
}

func Entrypoint(conf *config.Config) {
	receiver, err := zmq4.NewSocket(zmq4.PULL)
	if err != nil {
		panic(err)
	}
	defer receiver.Close()

	err = receiver.Connect(conf.QueueConfig.Connect)
	if err != nil {
		panic(err)
	}

	for {
		select {
		case <-channels.Running:
			return
		default:
			msg, _ := receiver.RecvBytes(0)
			go process(msg)
		}
	}
}
