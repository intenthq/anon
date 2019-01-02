package main

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var defaultJsonConfig = JsonConfig{
	IDField: "id",
}

func TestAnonymiseJson(t *testing.T) {
	config := func(mod uint32) *Config {
		return &Config{Json: &defaultJsonConfig, Sampling: SamplingConfig{Mod: mod}}
	}
	createDecoderAndEncoder := func(in string) (*json.Decoder, *json.Encoder, *bytes.Buffer) {
		var out bytes.Buffer
		r := json.NewDecoder(strings.NewReader(in))
		w := json.NewEncoder(&out)
		return r, w, &out
	}
	t.Run("when the json is not valid", func(t *testing.T) {
		r, w, _ := createDecoderAndEncoder(`not a json`)

		err := anonymiseJson(r, w, config(1), map[string]Anonymisation{"id": hash("salt")})
		assert.Error(t, err, "should return an error")
	})
	t.Run("when there is an error applying one of the anonymisations", func(t *testing.T) {
		input := `{"id": "id", "date": "not a date"}`
		r, w, out := createDecoderAndEncoder(input)

		y, _ := year("20060102")
		err := anonymiseJson(r, w, config(1), map[string]Anonymisation{"id": identity, "date": y})
		assert.NoError(t, err, "should return no error")
		assert.JSONEq(t, input, out.String(), "should leave the field with the error untouched")
	})
	t.Run("when a field doesn't have an anonymisation defined", func(t *testing.T) {
		r, w, out := createDecoderAndEncoder(`{"id": "id", "field": "don't touch it"}`)

		err := anonymiseJson(r, w, config(1), map[string]Anonymisation{"id": hash("salt")})
		assert.NoError(t, err, "should return no error")
		assert.JSONEq(t, `{"id": "58619739af7a7374f30a027fe40313491e678ed9", "field": "don't touch it"}`, out.String(), "should leave the field with the error untouched")
	})
	t.Run("when sampling is defined", func(t *testing.T) {
		r, w, out := createDecoderAndEncoder(`{"id": "1"}
		{"id": "2"}
		{"id": "3"}
		{"id": "4"}`)

		err := anonymiseJson(r, w, config(2), map[string]Anonymisation{"id": identity})
		assert.NoError(t, err, "should return no error")
		assert.Equal(t, "{\"id\":\"1\"}\n{\"id\":\"3\"}\n", out.String(), "should process some rows")
	})
	t.Run("when all the rows are valid", func(t *testing.T) {
		r, w, out := createDecoderAndEncoder(`{"id": "1"}
		{"id": "2"}`)

		err := anonymiseJson(r, w, config(1), map[string]Anonymisation{"id": identity})
		assert.NoError(t, err, "should return no error")
		assert.Equal(t, "{\"id\":\"1\"}\n{\"id\":\"2\"}\n", out.String(), "should process all rows")
	})
}
