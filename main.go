package main

import (
	"flag"
	"hash/fnv"
	"log"
	"math/rand"
	"os"
	"time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	//TODO move args parsing to a function
	configFile := flag.String("config", "config.json", "Configuration of the data to be anonymised. Default is 'config.json'")
	outputFile := flag.String("output", "", "Output file. Default is stdout.")
	flag.Parse()
	log.Printf("Using configuration in file %s\n", *configFile)
	conf, err := loadConfig(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	if conf.Json != nil && conf.Csv != nil {
		log.Fatal("Need to specify only one of Json or Csv.")
	} else if &conf.Json != nil {
		if err := processJson(flag.Arg(0), *outputFile, conf); err != nil {
			log.Fatal(err)
		}
	} else if &conf.Csv != nil {
		if err := processCsv(flag.Arg(0), *outputFile, conf); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal("Need to specify at least one of Json or Csv.")
	}

}

func sample(s string, conf SamplingConfig) bool {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()%conf.Mod == 0
}

// If filename is empty, will return `def`, if it's not, will return the
// result of the function `action` after passing `filename` ot it.
func fileOr(filename string, def *os.File, action func(string) (*os.File, error)) *os.File {
	if filename == "" {
		return def
	}
	f, err := action(filename)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

func anonymise(record []string, anons []Anonymisation) ([]string, error) {
	var err error
	for i := range record {
		// TODO decide if we fail if not enough anonmisations are defined
		// or we just skip the column (i.e. we apply identity)
		if i < len(anons) {
			if record[i], err = anons[i](record[i]); err != nil {
				return nil, err
			}
		}
	}
	return record, nil
}
