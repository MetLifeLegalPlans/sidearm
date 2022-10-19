package client

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/MetLifeLegalPlans/sidearm/config"

	. "github.com/logrusorgru/aurora/v3"
	"github.com/vmihailenco/msgpack/v5"
)

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
