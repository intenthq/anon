package main

import (
	"encoding/json"
	"os"
)

// CsvConfig stores the config to read and write the csv file
type CsvConfig struct {
	Delimiter string
	IDColumn  uint32
}

// JsonConfig stores the config to read and write the json file
type JsonConfig struct {
	IDField string
}

// SamplingConfig stores the config to know how to sample the file
type SamplingConfig struct {
	Mod uint32
}

// Config stores all the configuration
type Config struct {
	Csv      *CsvConfig
	Json     *JsonConfig
	Sampling SamplingConfig
	Actions  []ActionConfig
}

var defaultSamplingConfig = SamplingConfig{
	Mod: 1,
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
		Sampling: defaultSamplingConfig,
		Actions:  defaultActionsConfig,
	}
	err = decoder.Decode(&conf)
	if err != nil {
		return nil, err
	}
	return &conf, err
}
