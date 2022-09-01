package config

import "net/http"

func (s *Scenario) SetDefaults() {
	if s.URL == "" {
		panic("All scenarios must have URLs")
	}

	// Our weighted choice engine requires weights to be a minimum of 1
	s.Weight += 1

	if s.Method == "" {
		s.Method = http.MethodGet
	}
}
