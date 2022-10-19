package server

import (
	"github.com/mroth/weightedrand"

	"github.com/vmihailenco/msgpack/v5"
	"github.com/MetLifeLegalPlans/sidearm/config"
	"time"
)

type ServerContext struct {
	Conf    *config.Config
	Chooser *weightedrand.Chooser
	Start   time.Time
	End     time.Time
}

func NewContext(conf *config.Config) *ServerContext {
	sc := &ServerContext{}

	sc.Conf = conf
	sc.Start = time.Now()
	sc.End = sc.Start.Add(time.Duration(conf.Duration) * time.Second)

	// Pre-encoding all of our choices here so that our message
	// throughput is only limited by the speed of the socket
	choices := make([]weightedrand.Choice, len(conf.Scenarios))
	for idx, val := range conf.Scenarios {
		// Msgpack is SUBSTANTIALLY faster than JSON or YAML
		// and can be easily transmitted as raw bytes without
		// re-encoding
		packed, _ := msgpack.Marshal(val)
		choices[idx] = weightedrand.NewChoice(packed, val.Weight)
	}

	chooser, cerr := weightedrand.NewChooser(choices...)
	if cerr != nil {
		panic(cerr)
	}
	sc.Chooser = chooser

	return sc
}

func (s *ServerContext) Choose() []byte {
	return s.Chooser.Pick().([]byte)
}

func (s *ServerContext) RequestsForSecond() int64 {
	now := time.Now()
	// If we're not ramping up then we send out the full amount of requests every second
	if now.After(s.End) {
		return s.Conf.Requests
	}

	// If we are ramping up, then get the ratio of our timestamp to the target
	elapsed := now.Sub(s.Start)
	fullSpan := s.End.Sub(s.Start)

	ratio := float64(elapsed) / float64(fullSpan)
	desired := ratio * float64(s.Conf.Requests)

	// float -> int conversion truncates past the decimal point
	// Basically just a faster Math.Floor
	return int64(desired)
}
