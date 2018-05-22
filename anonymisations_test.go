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
			return assert.NoError(t, err) && assert.Equal(t, v, res)
		},
		gen.AnyString(),
	))

	properties.TestingRun(t)
}

func TestHash(t *testing.T) {
	res, err := hash("")
	assert.NoError(t, err)
	assert.Equal(t, "da39a3ee5e6b4b0d3255bfef95601890afd80709", res)
	res, err = hash("hasselhoff")
	assert.Equal(t, "ffe3294fad149c2dd3579cb864a1aebb2201f38d", res)
}

func TestOutcode(t *testing.T) {
	properties := gopter.NewProperties(nil)

	properties.Property("Same output as input", prop.ForAll(
		func(v1 string, v2 string) bool {
			res, err := outcode(v1 + " " + v2)
			return assert.NoError(t, err) && assert.Equal(t, v1, res)
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
		assert.NoError(t, err, "should return no error")
		assert.Equal(t, "2012", res, "should return the year")
	})
	t.Run("if the date cannot be parsed", func(t *testing.T) {
		res, err := f("input")
		assert.Error(t, err, "should return an error")
		assert.Equal(t, "input", res, "should return the input unchanged")
	})
}
func TestRanges(t *testing.T) {
	min := 0.0
	max := 100.0
	output := "0-100"
	f := ranges([]RangeConfig{RangeConfig{Gt: &min, Lte: &max, Output: &output}})
	t.Run("if the value is not a float", func(t *testing.T) {
		res, err := f("input")
		assert.Error(t, err, "should return an error")
		assert.Equal(t, "input", res, "should return the input unchanged")
	})
	t.Run("if the value is a float", func(t *testing.T) {
		t.Run("not in any range", func(t *testing.T) {
			res, err := f("2000")
			assert.Error(t, err, "should return an error")
			assert.Equal(t, "2000", res, "should return the input unchanged")
		})
		t.Run("inside a range", func(t *testing.T) {
			res, err := f("10")
			assert.NoError(t, err, "should return no error")
			assert.Equal(t, output, res, "should return the output")
		})
	})
}

func TestRangeConfigContains(t *testing.T) {
	min := 0.0
	max := 100.0
	t.Run("range containing only lt", func(t *testing.T) {
		conf := RangeConfig{Lt: &max}
		assert.True(t, conf.contains(max-1))
		assert.False(t, conf.contains(max))
		assert.False(t, conf.contains(max+1))
	})
	t.Run("range containing only lte", func(t *testing.T) {
		conf := RangeConfig{Lte: &max}
		assert.True(t, conf.contains(max-1))
		assert.True(t, conf.contains(max))
		assert.False(t, conf.contains(max+1))
	})
	t.Run("range containing only gt", func(t *testing.T) {
		conf := RangeConfig{Gt: &min}
		assert.False(t, conf.contains(min-1))
		assert.False(t, conf.contains(min))
		assert.True(t, conf.contains(min+1))
	})
	t.Run("range containing only gte", func(t *testing.T) {
		conf := RangeConfig{Gte: &min}
		assert.False(t, conf.contains(min-1))
		assert.True(t, conf.contains(min))
		assert.True(t, conf.contains(min+1))
	})
	t.Run("range containing gt and lt", func(t *testing.T) {
		conf := RangeConfig{Gt: &min, Lt: &max}
		assert.False(t, conf.contains(min-1))
		assert.False(t, conf.contains(min))
		assert.True(t, conf.contains(min+1))
		assert.False(t, conf.contains(max))
		assert.False(t, conf.contains(max+1))
	})
	t.Run("range containing gt and lte", func(t *testing.T) {
		conf := RangeConfig{Gt: &min, Lte: &max}
		assert.False(t, conf.contains(min-1))
		assert.False(t, conf.contains(min))
		assert.True(t, conf.contains(min+1))
		assert.True(t, conf.contains(max))
		assert.False(t, conf.contains(max+1))
	})
	t.Run("range containing gte and lt", func(t *testing.T) {
		conf := RangeConfig{Gte: &min, Lt: &max}
		assert.False(t, conf.contains(min-1))
		assert.True(t, conf.contains(min))
		assert.True(t, conf.contains(min+1))
		assert.False(t, conf.contains(max))
		assert.False(t, conf.contains(max+1))
	})
	t.Run("range containing gte and lte", func(t *testing.T) {
		conf := RangeConfig{Gte: &min, Lte: &max}
		assert.False(t, conf.contains(min-1))
		assert.True(t, conf.contains(min))
		assert.True(t, conf.contains(min+1))
		assert.True(t, conf.contains(max))
		assert.False(t, conf.contains(max+1))
	})
}
