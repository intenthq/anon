package main

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

// Anonymisation is a function that transforms a string into another one
type Anonymisation func(string) (string, error)

// DateConfig stores the format (layout) of an input date
type DateConfig struct {
	Format string
}

// RangeConfig stores configuration to define a range of values
type RangeConfig struct {
	Gt     *float64
	Gte    *float64
	Lt     *float64
	Lte    *float64
	Output *string
}

// ActionConfig stores the config of an anonymisation action
type ActionConfig struct {
	Name        string
	DateConfig  DateConfig
	RangeConfig []RangeConfig
}

// Returns an array of anonymisations according to the config
func anonymisations(configs *[]ActionConfig) []Anonymisation {
	res := make([]Anonymisation, len(*configs))
	for i, config := range *configs {
		res[i] = config.create()
	}
	return res
}

func (ac *ActionConfig) create() Anonymisation {
	switch ac.Name {
	case "nothing":
		return identity
	case "outcode":
		return outcode
	case "hash":
		return hash
	case "year":
		return year(ac.DateConfig.Format)
	case "ranges":
		return ranges(ac.RangeConfig)
	}
	return identity
}

// The no-op, returns the input unchanged.
func identity(s string) (string, error) {
	return s, nil
}

// Hashes (SHA1) the input.
func hash(s string) (string, error) {
	h := sha1.New()
	io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// Takes a UK format postcode (eg. W1W 8BE) and just keeps
// the outcode (eg. W1W).
// i.e. returns the prefix of the input until it finds a space
func outcode(s string) (string, error) {
	return strings.Split(s, " ")[0], nil
}

// Given a date format/layout, it returns a function that
// given a date in that format, just keeps the year.
// If either the format is invalid or the year doesn't
// match that format, it will return an error and
// the input unchanged
func year(format string) Anonymisation {
	return func(s string) (string, error) {
		t, err := time.Parse(format, s)
		if err != nil {
			return s, err
		}
		return strconv.Itoa(t.Year()), nil
	}
}

// Given a list of ranges, it will summarise numeric
// values into groups of values, each group defined
// by a range and an output
func ranges(ranges []RangeConfig) Anonymisation {
	return func(s string) (string, error) {
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return s, err
		}
		for _, rang := range ranges {
			if rang.contains(v) {
				return *rang.Output, nil
			}
		}
		return s, errors.New("No range defined for value")
	}
}

func (r *RangeConfig) contains(v float64) bool {
	return (r.Gt == nil && r.Gte == nil || r.Gt != nil && *r.Gt < v || r.Gte != nil && *r.Gte <= v) &&
		(r.Lt == nil && r.Lte == nil || r.Lt != nil && *r.Lt > v || r.Lte != nil && *r.Lte >= v)
}
