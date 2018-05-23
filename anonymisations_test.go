package main

import (
	"math/rand"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"github.com/stretchr/testify/assert"
)

func TestAnonymisations(t *testing.T) {
	salt := "jump"
	conf := &[]ActionConfig{
		ActionConfig{
			Name: "nothing",
		},
		ActionConfig{
			Name: "hash",
			Salt: &salt,
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
	expectedRes, expectedErr = hash("jump")("a")
	actualRes, actualErr = anons[1]("a")
	assert.Equal(t, expectedRes, actualRes)
	assert.Equal(t, expectedErr, actualErr)
}

func TestActionConfig(t *testing.T) {
	t.Run("saltOrRandom", func(t *testing.T) {
		t.Run("if salt is not specified", func(t *testing.T) {
			rand.Seed(1)
			acNoSalt := ActionConfig{Name: "hash"}
			assert.Equal(t, "5577006791947779410", acNoSalt.saltOrRandom(), "should return a random salt")
		})
		t.Run("if salt is specified", func(t *testing.T) {
			emptySalt := ""
			acEmptySalt := ActionConfig{Name: "hash", Salt: &emptySalt}
			assert.Empty(t, acEmptySalt.saltOrRandom(), "should return the empty salt if empty")

			salt := "jump"
			acSalt := ActionConfig{Name: "hash", Salt: &salt}
			assert.Equal(t, "jump", acSalt.saltOrRandom(), "should return the salt")
		})
	})
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
	t.Run("should hash the values using sha1 without a salt", func(t *testing.T) {
		unsaltedHash := hash("")
		res, err := unsaltedHash("")
		assert.NoError(t, err)
		assert.Equal(t, "da39a3ee5e6b4b0d3255bfef95601890afd80709", res)
		res, err = unsaltedHash("hasselhoff")
		assert.Equal(t, "ffe3294fad149c2dd3579cb864a1aebb2201f38d", res)
	})
	t.Run("should use the salt if provided", func(t *testing.T) {
		properties := gopter.NewProperties(nil)

		properties.Property("hash(salt)(s) == hash(s+salt)", prop.ForAll(
			func(salt string, s string) bool {
				res1, err1 := hash(salt)(s)
				res2, err2 := hash("")(s + salt)
				return assert.NoError(t, err1) && assert.NoError(t, err2) && assert.Equal(t, res1, res2)
			},
			gen.AlphaString(),
			gen.AlphaString(),
		))
	})
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
