package dashboard

import (
	"sidearm/config"
)

type Result struct {
	StatusCode int
	Duration   int64
	Scenario   config.Scenario
}
