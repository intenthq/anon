package main

import (
	"bytes"
	"encoding/csv"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitReader(t *testing.T) {
	t.Run("with an empty filename", func(t *testing.T) {
		tmpfile := tmpFile("content")
		defer os.Remove(tmpfile.Name()) // clean up

		oldStdin := os.Stdin
		defer func() { os.Stdin = oldStdin }() // Restore original Stdin
		os.Stdin = tmpfile

		r := initReader("", defaultCsvConfig)
		record, err := r.Read()

		assert.NoError(t, err, "should return no error")
		assert.Equal(t, []string{"content"}, record, "should return a csv reader that reads from stdin")
	})
	t.Run("with a valid filename", func(t *testing.T) {
		tmpfile := tmpFile("content")
		defer os.Remove(tmpfile.Name()) // clean up

		r := initReader(tmpfile.Name(), defaultCsvConfig)
		record, err := r.Read()

		assert.NoError(t, err, "should return no error")
		assert.Equal(t, []string{"content"}, record, "should return a csv reader that reads from the file")
	})
}

func tmpFile(content string) *os.File {
	tmpfile, err := ioutil.TempFile("", "anon-test")
	if err != nil {
		log.Fatal(err)
	}
	ioutil.WriteFile(tmpfile.Name(), []byte("content"), os.ModePerm)
	return tmpfile
}

func TestInitWriter(t *testing.T) {
	t.Run("with an empty filename", func(t *testing.T) {
		tmpfile := tmpFile("")
		defer os.Remove(tmpfile.Name()) // clean up

		oldStdout := os.Stdout
		defer func() { os.Stdout = oldStdout }() // Restore original Stdout
		os.Stdout = tmpfile

		w := initWriter("", defaultCsvConfig)
		err := w.Write([]string{"csv", "content"})
		w.Flush()

		content, _ := ioutil.ReadFile(tmpfile.Name())
		assert.NoError(t, err, "should return no error")
		assert.Equal(t, "csv,content\n", string(content), "should return a csv writer that writes to stdout")
	})
	t.Run("with a valid filename", func(t *testing.T) {
		tmpfile := tmpFile("")
		defer os.Remove(tmpfile.Name()) // clean up

		w := initWriter(tmpfile.Name(), defaultCsvConfig)
		err := w.Write([]string{"csv", "content"})
		w.Flush()

		content, _ := ioutil.ReadFile(tmpfile.Name())
		assert.NoError(t, err, "should return no error")
		assert.Equal(t, "csv,content\n", string(content), "should return a csv writer that writes to stdout")
	})
}
func TestFileOr(t *testing.T) {
	assert.Equal(t, fileOr("", os.Stdin, stdOutOk), os.Stdin, "with an empty filename returns the default value")
	assert.Equal(t, fileOr("something", os.Stdin, stdOutOk), os.Stdout, "with non empty filename returns the value returned by the action")
}

func stdOutOk(s string) (*os.File, error) {
	return os.Stdout, nil
}

func TestAnonymise(t *testing.T) {
	record := []string{"a", "b", "c"}
	actions := []Anonymisation{identity, hash(""), identity}
	output := []string{"a", "e9d71f5ee7c92d6dc9e92ffdad17b8bd49418f98", "c"}
	res, err := anonymise(record, actions)
	assert.NoError(t, err)
	assert.Equal(t, output, res, "should apply anonymisation functions to each column in the record")
}

func TestSample(t *testing.T) {
	conf := SamplingConfig{
		Mod: 2,
	}
	assert.True(t, sample("a", conf))
	assert.False(t, sample("b", conf))
}

func TestProcess(t *testing.T) {
	config := func(mod uint32, idColumn uint32) *Config {
		return &Config{Sampling: SamplingConfig{Mod: mod, IDColumn: idColumn}}
	}
	anons := &[]Anonymisation{identity, outcode}
	createReaderAndWriter := func(in string) (*csv.Reader, *csv.Writer, *bytes.Buffer) {
		var out bytes.Buffer
		r := csv.NewReader(strings.NewReader(in))

		w := csv.NewWriter(&out)
		return r, w, &out
	}
	t.Run("when the id column is out of range", func(t *testing.T) {
		r, w, out := createReaderAndWriter("a,b c\nd,e f\n")

		err := process(r, w, config(1, 100), anons)
		assert.Error(t, err, "should return an error")
		assert.Equal(t, "", out.String(), "shouldn't write any output")
	})
	t.Run("when there is an error writing the output", func(t *testing.T) {
		var out bytes.Buffer
		f, _ := os.Open("non existing file")
		r := csv.NewReader(f)

		w := csv.NewWriter(&out)
		err := process(r, w, config(1, 0), anons)
		assert.Error(t, err, "should return an error")
	})
	t.Run("when there is an error processing one of the rows", func(t *testing.T) {
		r, w, out := createReaderAndWriter("20020202\nfail\n10010101")

		y, _ := year("20060102")
		err := process(r, w, config(1, 0), &[]Anonymisation{y})
		assert.NoError(t, err, "should not return an error")
		assert.Equal(t, "2002\n1001\n", out.String(), "should skip that row")
	})
	t.Run("when sampling is defined", func(t *testing.T) {
		r, w, out := createReaderAndWriter("a,b c\nd,e f\ng,h i\nj,k l\n")

		err := process(r, w, config(2, 0), anons)
		assert.NoError(t, err, "should return no error")
		assert.Equal(t, "a,b\ng,h\n", out.String(), "should process some rows")
	})
	t.Run("when all the rows are valid", func(t *testing.T) {
		r, w, out := createReaderAndWriter("a,b c\nd,e f\n")

		err := process(r, w, config(1, 0), anons)
		assert.NoError(t, err, "should return no error")
		assert.Equal(t, "a,b\nd,e\n", out.String(), "should process all rows")
	})
}
