package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileOr(t *testing.T) {
	assert.Equal(t, fileOr("", os.Stdin, stdOutOk), os.Stdin, "with an empty filename returns the default value")
	assert.Equal(t, fileOr("something", os.Stdin, stdOutOk), os.Stdout, "with non empty filename returns the value returned by the action")
}

func stdOutOk(s string) (*os.File, error) {
	return os.Stdout, nil
}

func TestAnonymise(t *testing.T) {
	record := []string{"1", "2", "3"}
	actions := []anonymisation{identity, hash, identity}
	output := []string{identity("1"), hash("2"), identity("3")}
	assert.Equal(t, anonymise(record, actions), output, "should apply anonymisation functions to each column in the record")
}
