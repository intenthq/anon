package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"strings"
)

type anonymisation func(string) string

// Returns the input without changing it
func identity(s string) string {
	return s
}

// Returns the sha1 of the input
func hash(s string) string {
	h := sha1.New()
	io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// Returns the prefix of the input until it finds a space
func outcode(s string) string {
	return strings.Split(s, " ")[0]
}
