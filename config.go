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

func loadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(file)
	conf := Config{}
	err = decoder.Decode(&conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}
