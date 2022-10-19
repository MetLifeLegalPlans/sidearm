package dashboard

import (
	"github.com/MetLifeLegalPlans/sidearm/config"
)

type Result struct {
	StatusCode int
	Duration   int64
	Scenario   config.Scenario
}
