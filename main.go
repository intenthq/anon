package main

import (
	"encoding/csv"
	"flag"
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
	r := initReader(flag.Arg(0))
	w := initWriter(*outputFile)
	//TODO create the anons array using the config provided
	anons := []anonymisation{hash, identity, identity, outcode}
	i := 0

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		} else if pe, ok := err.(*csv.ParseError); ok && pe.Err == csv.ErrFieldCount {
			// we just print the error and skip the column
			log.Print(err)
		} else if err != nil {
			log.Fatal(err)
		} else {
			w.Write(anonymise(record, anons))
			//TODO decide how often do we want to flush
			if i%100 == 0 {
				w.Flush()
			}
		}
		i++
	}
	w.Flush()
}

func initReader(file string) *csv.Reader {
	return csv.NewReader(fileOr(file, os.Stdout, os.Open))
}

func initWriter(file string) *csv.Writer {
	return csv.NewWriter(fileOr(file, os.Stdin, os.Create))
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

func anonymise(record []string, anons []anonymisation) []string {
	for i := range record {
		//TODO decide if we fail if not enough anonmisations are defined or we just skip the record
		if i < len(anons) {
			record[i] = anons[i](record[i])
		}
	}
	return record
}
