package main

import (
	"encoding/json"
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
