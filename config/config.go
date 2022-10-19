package config

import (
	"errors"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	QueueConfig AddressPair `yaml:"queue"`
	SinkConfig  AddressPair `yaml:"sink"`
	Duration    int
	Requests    int64  `yaml:"requests"`
	DbPath      string `yaml:"dbPath"`
	BatchSize   int    `yaml:"batchSize"`

	Scenarios []Scenario
}

type AddressPair struct {
	Bind    string
	Connect string
}

type Scenario struct {
	URL    string
	Method string
	Body   map[string]any

	Weight uint
}

func (c *Config) SetDefaults() {
	if c.Requests <= 0 {
		panic("`requests` is a required configuration field")
	}

	if c.SinkConfig.Enabled() && c.DbPath == "" {
		envPath := os.Getenv("DB_PATH")

		if envPath == "" {
			log.Fatalf("Result sink enabled but config.dbPath is unset, exiting")
		}
		c.DbPath = envPath
	}

	if c.BatchSize == 0 {
		c.BatchSize = 100
	}

	c.QueueConfig.SetDefaults()
	c.SinkConfig.SetDefaults()

	for idx := range c.Scenarios {
		c.Scenarios[idx].SetDefaults()
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)

	return !errors.Is(err, os.ErrNotExist)
}

func Load(path string) (*Config, error) {
	result := Config{}

	if !fileExists(path) {
		return &result, os.ErrNotExist
	}

	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	result.SetDefaults()
	return &result, nil
}
