package main

import (
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"github.com/stretchr/testify/assert"
)

func TestIdentity(t *testing.T) {
	properties := gopter.NewProperties(nil)

	properties.Property("Same output as input", prop.ForAll(
		func(v string) bool {
			return assert.Equal(t, v, identity(v))
		},
		gen.AnyString(),
	))

	properties.TestingRun(t)
}

func TestHash(t *testing.T) {
	assert.Equal(t, "da39a3ee5e6b4b0d3255bfef95601890afd80709", hash(""))
	assert.Equal(t, "ffe3294fad149c2dd3579cb864a1aebb2201f38d", hash("hasselhoff"))
}

func TestOutcode(t *testing.T) {
	properties := gopter.NewProperties(nil)

	properties.Property("Same output as input", prop.ForAll(
		func(v1 string, v2 string) bool {
			return assert.Equal(t, v1, outcode(v1+" "+v2))
		},
		gen.AlphaString(),
		gen.AlphaString(),
	))

	properties.TestingRun(t)
}
