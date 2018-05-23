package main

import (
	"io/ioutil"
	"log"
	"os"
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
	actions := []Anonymisation{identity, hash, identity}
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
