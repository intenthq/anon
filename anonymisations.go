package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"strings"
)

type anonymisation func(string) string

func identity(s string) string {
	return s
}

func hash(s string) string {
	h := sha1.New()
	io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func outcode(s string) string {
	return strings.Split(s, " ")[0]
}
