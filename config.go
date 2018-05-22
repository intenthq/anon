package main

import (
	"encoding/json"
	"errors"
	"os"
)

// CsvConfig stores the config to read and write the csv file
type CsvConfig struct {
	Delimiter string
}

// SamplingConfig stores the config to know how to sample the file
type SamplingConfig struct {
	Mod      uint32
	IDColumn uint32
}

// Config stores all the configuration
type Config struct {
	Csv      CsvConfig
	Sampling SamplingConfig
	Actions  []ActionConfig
}

var defaultCsvConfig = CsvConfig{
	Delimiter: ",",
}

var defaultSamplingConfig = SamplingConfig{
	Mod:      1,
	IDColumn: 0,
}

var defaultActionsConfig = []ActionConfig{}

func loadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(file)
	conf := Config{
		Csv:      defaultCsvConfig,
		Sampling: defaultSamplingConfig,
		Actions:  defaultActionsConfig,
	}
	err = decoder.Decode(&conf)
	if err != nil {
		return nil, err
	}
	return &conf, err
}

func (conf *Config) validate() error {
	for _, action := range conf.Actions {
		for _, rc := range action.RangeConfig {
			if rc.Gt != nil && rc.Gte != nil || rc.Lt != nil && rc.Lte != nil {
				return errors.New("You can only specify one of (gt, gte) and (lt, lte)")
			} else if rc.Gt == nil && rc.Gte == nil && rc.Lt == nil && rc.Lte == nil {
				return errors.New("You need to specify at least one of gt, gte, lt, lte")
			} else if rc.Output == nil {
				return errors.New("You need to specify the output for a range")
			}
		}
	}
	return nil
}
