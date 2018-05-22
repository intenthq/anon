package main

import (
	"encoding/csv"
	"flag"
	"hash/fnv"
	"io"
	"log"
	"os"
)

func main() {
	//TODO move args parsing to a function
	configFile := flag.String("config", "config.json", "Configuration of the data to be anonymised. Default is 'config.json'")
	outputFile := flag.String("output", "", "Output file. Default is stdout.")
	flag.Parse()
	log.Printf("Using configuration in file %s\n", *configFile)
	conf, err := loadConfig(*configFile)
	if err != nil {
		log.Fatal(err)
	} else if err = conf.validate(); err != nil {
		log.Fatal(err)
	}
	r := initReader(flag.Arg(0), conf.Csv)
	w := initWriter(*outputFile, conf.Csv)
	anons := anonymisations(&conf.Actions)
	i := 0

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		} else if pe, ok := err.(*csv.ParseError); ok && pe.Err == csv.ErrFieldCount {
			// we just print the error and skip the record
			log.Print(err)
		} else if err != nil {
			log.Fatal(err)
		} else if sample(record[conf.Sampling.IDColumn], conf.Sampling) {
			anonymised, err := anonymise(record, anons)
			if err != nil {
				// we just print the error and skip the record
				log.Print(err)
			} else {
				w.Write(anonymised)
			}
			//TODO decide how often do we want to flush
			if i%100 == 0 {
				w.Flush()
			}
		}
		i++
	}
	w.Flush()
}

func sample(s string, conf SamplingConfig) bool {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()%conf.Mod == 0
}

func initReader(filename string, conf CsvConfig) *csv.Reader {
	reader := csv.NewReader(fileOr(filename, os.Stdout, os.Open))
	reader.Comma = []rune(conf.Delimiter)[0]
	return reader
}

func initWriter(filename string, conf CsvConfig) *csv.Writer {
	writer := csv.NewWriter(fileOr(filename, os.Stdin, os.Create))
	writer.Comma = []rune(conf.Delimiter)[0]
	return writer
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
