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
	record := []string{"a", "b", "c"}
	actions := []Anonymisation{identity, hash, identity}
	output := []string{"a", "e9d71f5ee7c92d6dc9e92ffdad17b8bd49418f98", "c"}
	res, err := anonymise(record, actions)
	assert.Nil(t, err)
	assert.Equal(t, output, res, "should apply anonymisation functions to each column in the record")
}

func TestSample(t *testing.T) {
	conf := SamplingConfig{
		Mod: 2,
	}
	assert.True(t, sample("a", conf))
	assert.False(t, sample("b", conf))
}
