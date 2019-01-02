package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
)

func processCsv(inputFile string, outputFile string, conf *Config) error {
	r := initCsvReader(inputFile, conf.Csv)
	w := initCsvWriter(outputFile, conf.Csv)

	anons, err := anonymisations(&conf.Actions)
	if err != nil {
		return err
	}

	if err := anonymiseCsv(r, w, conf, &anons); err != nil {
		return err
	}

	return nil
}

func initCsvReader(filename string, conf *CsvConfig) *csv.Reader {
	reader := csv.NewReader(fileOr(filename, os.Stdin, os.Open))
	reader.Comma = []rune(conf.Delimiter)[0]
	return reader
}

func initCsvWriter(filename string, conf *CsvConfig) *csv.Writer {
	writer := csv.NewWriter(fileOr(filename, os.Stdout, os.Create))
	writer.Comma = []rune(conf.Delimiter)[0]
	return writer
}

func anonymiseCsv(r *csv.Reader, w *csv.Writer, conf *Config, anons *[]Anonymisation) error {
	i := 0

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		} else if pe, ok := err.(*csv.ParseError); ok && pe.Err == csv.ErrFieldCount {
			// we just print the error and skip the record
			log.Print(err)
		} else if err != nil {
			return err
		} else if int64(conf.Csv.IDColumn) >= int64(len(record)) {
			return fmt.Errorf("id column (%d) out of range, record has %d columns", conf.Csv.IDColumn, len(record))
		} else if sample(record[conf.Csv.IDColumn], conf.Sampling) {
			anonymised, err := anonymise(record, *anons)
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
	return nil
}
