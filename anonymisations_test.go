package main

import (
	"math/rand"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var salt = "jump"

const seed = int64(1)

//this is the first random salt with the seed above
const firstSalt = "5577006791947779410"

// can't test that the functions are equal because of https://github.com/stretchr/testify/issues/182
// and https://github.com/stretchr/testify/issues/159#issuecomment-99557398
// will have to test that the functions return the same
func assertAnonymisationFunction(t *testing.T, expected Anonymisation, actual Anonymisation, value string) {
	require.NotNil(t, expected)
	require.NotNil(t, actual)
	expectedRes, expectedErr := expected(value)
	actualRes, actualErr := actual(value)
	assert.Equal(t, expectedRes, actualRes)
	assert.Equal(t, expectedErr, actualErr)
}

func TestAnonymisations(t *testing.T) {
	t.Run("a valid configuration", func(t *testing.T) {
		conf := &[]ActionConfig{
			ActionConfig{
				Name: "nothing",
			},
			ActionConfig{
				Name: "hash",
				Salt: &salt,
			},
		}
		anons, err := anonymisations(conf)
		assert.NoError(t, err)
		assertAnonymisationFunction(t, identity, anons[0], "a")
		assertAnonymisationFunction(t, hash(salt), anons[1], "a")
	})
	t.Run("an invalid configuration", func(t *testing.T) {
		conf := &[]ActionConfig{ActionConfig{Name: "year", DateConfig: DateConfig{Format: "3333"}}}
		anons, err := anonymisations(conf)
		assert.Error(t, err, "should return an error")
		assert.Nil(t, anons)
	})
}

func TestActionConfigSaltOrRandom(t *testing.T) {
	t.Run("if salt is not specified", func(t *testing.T) {
		rand.Seed(seed)
		acNoSalt := ActionConfig{Name: "hash"}
		assert.Equal(t, firstSalt, acNoSalt.saltOrRandom(), "should return a random salt")
	})
	t.Run("if salt is specified", func(t *testing.T) {
		emptySalt := ""
		acEmptySalt := ActionConfig{Name: "hash", Salt: &emptySalt}
		assert.Empty(t, acEmptySalt.saltOrRandom(), "should return the empty salt if empty")

		acSalt := ActionConfig{Name: "hash", Salt: &salt}
		assert.Equal(t, "jump", acSalt.saltOrRandom(), "should return the salt")
	})
}

func TestActionConfigCreate(t *testing.T) {
	t.Run("invalid name", func(t *testing.T) {
		ac := ActionConfig{Name: "invalid name"}
		res, err := ac.create()
		assert.Error(t, err)
		assert.Nil(t, res)
	})
	t.Run("identity", func(t *testing.T) {
		ac := ActionConfig{Name: "nothing"}
		res, err := ac.create()
		assert.NoError(t, err)
		assertAnonymisationFunction(t, identity, res, "a")
	})
	t.Run("outcode", func(t *testing.T) {
		ac := ActionConfig{Name: "outcode"}
		res, err := ac.create()
		assert.NoError(t, err)
		assertAnonymisationFunction(t, outcode, res, "a")
	})
	t.Run("hash", func(t *testing.T) {
		t.Run("if salt is not specified uses a random salt", func(t *testing.T) {
			rand.Seed(1)
			ac := ActionConfig{Name: "hash"}
			res, err := ac.create()
			assert.NoError(t, err)
			assertAnonymisationFunction(t, hash(firstSalt), res, "a")
		})
		t.Run("if salt is specified uses it", func(t *testing.T) {
			ac := ActionConfig{Name: "hash", Salt: &salt}
			res, err := ac.create()
			assert.NoError(t, err)
			assertAnonymisationFunction(t, hash(salt), res, "a")
		})
	})
	t.Run("year", func(t *testing.T) {
		t.Run("with an invalid format", func(t *testing.T) {
			ac := ActionConfig{Name: "year", DateConfig: DateConfig{Format: "11112233"}}
			res, err := ac.create()
			assert.Error(t, err, "should fail")
			assert.Nil(t, res)
		})
		t.Run("with a valid format", func(t *testing.T) {
			ac := ActionConfig{Name: "year", DateConfig: DateConfig{Format: "20060102"}}
			res, err := ac.create()
			assert.NoError(t, err, "should not fail")
			y, err := year("20060102")
			assert.NoError(t, err)
			assertAnonymisationFunction(t, y, res, "21121212")
		})
	})
	t.Run("ranges", func(t *testing.T) {
		num := 2.0
		output := "0-100"
		t.Run("range has at least one of lt, lte, gt, gte", func(t *testing.T) {
			ac := ActionConfig{
				Name:        "ranges",
				RangeConfig: []RangeConfig{RangeConfig{Output: &output}},
			}
			r, err := ac.create()
			assert.Error(t, err, "if not should return an error")
			assert.Nil(t, r)
		})
		t.Run("range contains both lt and lte", func(t *testing.T) {
			ac := ActionConfig{
				Name:        "ranges",
				RangeConfig: []RangeConfig{RangeConfig{Lt: &num, Lte: &num, Output: &output}},
			}
			r, err := ac.create()
			assert.Error(t, err, "if not should return an error")
			assert.Nil(t, r)
		})
		t.Run("range contains both gt and gte", func(t *testing.T) {
			ac := ActionConfig{
				Name:        "ranges",
				RangeConfig: []RangeConfig{RangeConfig{Gt: &num, Gte: &num, Output: &output}},
			}
			r, err := ac.create()
			assert.Error(t, err, "if not should return an error")
			assert.Nil(t, r)
		})
		t.Run("range without output defined", func(t *testing.T) {
			ac := ActionConfig{
				Name:        "ranges",
				RangeConfig: []RangeConfig{RangeConfig{Lt: &num, Gte: &num}},
			}
			r, err := ac.create()
			assert.Error(t, err, "if not should return an error")
			assert.Nil(t, r)
		})
		t.Run("valid range", func(t *testing.T) {
			rangeConfigs := []RangeConfig{RangeConfig{Lte: &num, Gte: &num, Output: &output}}
			ac := ActionConfig{
				Name:        "ranges",
				RangeConfig: rangeConfigs,
			}
			r, err := ac.create()
			expected, _ := ranges(rangeConfigs)
			assert.NoError(t, err)
			assertAnonymisationFunction(t, expected, r, "2")
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
	f, _ := year("20060102")
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
	f, _ := ranges([]RangeConfig{RangeConfig{Gt: &min, Lte: &max, Output: &output}})
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
