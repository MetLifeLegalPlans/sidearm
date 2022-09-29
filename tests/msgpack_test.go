package tests

import (
	"encoding/json"
	"github.com/vmihailenco/msgpack/v5"
	"testing"

	"sidearm/config"
)

var scenario = &config.Scenario{
	URL:    "https://legalplans.com",
	Method: "GET",
	Body: map[string]any{
		"test": "key",
	},
	Weight: 1,
}

var mbuf, _ = msgpack.Marshal(scenario)
var jbuf, _ = json.Marshal(scenario)

func BenchmarkMsgpackEncode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		msgpack.Marshal(scenario)
	}
}

func BenchmarkMsgpackDecode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		msgpack.Unmarshal(mbuf, &config.Scenario{})
	}
}

func BenchmarkJsonEncode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		json.Marshal(scenario)
	}
}

func BenchmarkJsonDecode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		json.Unmarshal(jbuf, &config.Scenario{})
	}
}
