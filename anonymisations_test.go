package main

import (
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"github.com/stretchr/testify/assert"
)

func TestAnonymisations(t *testing.T) {
	conf := &[]ActionConfig{
		ActionConfig{
			Name: "nothing",
		},
		ActionConfig{
			Name: "hash",
		},
	}
	// can't test that the functions are equal because of https://github.com/stretchr/testify/issues/182
	// and https://github.com/stretchr/testify/issues/159#issuecomment-99557398
	// will have to test that the functions return the same
	anons := anonymisations(conf)
	expectedRes, expectedErr := identity("a")
	actualRes, actualErr := anons[0]("a")
	assert.Equal(t, expectedRes, actualRes)
	assert.Equal(t, expectedErr, actualErr)
	expectedRes, expectedErr = hash("a")
	actualRes, actualErr = anons[1]("a")
	assert.Equal(t, expectedRes, actualRes)
	assert.Equal(t, expectedErr, actualErr)
}

func TestIdentity(t *testing.T) {
	properties := gopter.NewProperties(nil)

	properties.Property("Same output as input", prop.ForAll(
		func(v string) bool {
			res, err := identity(v)
			return assert.Nil(t, err) && assert.Equal(t, v, res)
		},
		gen.AnyString(),
	))

	properties.TestingRun(t)
}

func TestHash(t *testing.T) {
	res, err := hash("")
	assert.Nil(t, err)
	assert.Equal(t, "da39a3ee5e6b4b0d3255bfef95601890afd80709", res)
	res, err = hash("hasselhoff")
	assert.Equal(t, "ffe3294fad149c2dd3579cb864a1aebb2201f38d", res)
}

func TestOutcode(t *testing.T) {
	properties := gopter.NewProperties(nil)

	properties.Property("Same output as input", prop.ForAll(
		func(v1 string, v2 string) bool {
			res, err := outcode(v1 + " " + v2)
			return assert.Nil(t, err) && assert.Equal(t, v1, res)
		},
		gen.AlphaString(),
		gen.AlphaString(),
	))

	properties.TestingRun(t)
}

func TestYear(t *testing.T) {
	f := year("20060102")
	t.Run("if the date can be parsed", func(t *testing.T) {
		res, err := f("20120102")
		assert.Nil(t, err, "should return no error")
		assert.Equal(t, "2012", res, "should return the year")
	})
	t.Run("if the date cannot be parsed", func(t *testing.T) {
		res, err := f("input")
		assert.Error(t, err, "should return an error")
		assert.Equal(t, "input", res, "should return the input unchanged")
	})
}
