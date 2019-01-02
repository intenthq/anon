package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

func processJson(inputFile string, outputFile string, conf *Config) error {
	r := json.NewDecoder(fileOr(inputFile, os.Stdin, os.Open))
	w := json.NewEncoder(fileOr(outputFile, os.Stdout, os.Create))

	anonsMap, err := anonymisationsMap(&conf.Actions)
	if err != nil {
		return err
	}

	return anonymiseJson(r, w, conf, anonsMap)
}

func anonymiseJson(r *json.Decoder, w *json.Encoder, conf *Config, anonsMap map[string]Anonymisation) error {
	var values map[string]string
	for {
		err := r.Decode(&values)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		} else if sample(values[conf.Json.IDField], conf.Sampling) {
			for key, value := range values {
				if anon, exists := anonsMap[key]; exists {
					anonValue, err := anon(value)
					if err != nil {
						log.Printf("Error applying anonymisation, field won't be anonymised. Error: %s\n", err)
					} else {
						values[key] = anonValue
					}
				}
			}
			w.Encode(values)
		}
	}
	return nil
}
