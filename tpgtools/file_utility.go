package main

import (
	"bytes"
	"go/format"
	"strings"
)

func stripExt(s string) string {
	n := strings.LastIndexByte(s, '.')
	if n >= 0 {
		return s[:n]
	}
	return s
}

func formatSource(source *bytes.Buffer) ([]byte, error) {
	output, err := format.Source(source.Bytes())
	if err != nil {
		return []byte(source.String()), err
	}

	return output, nil
}
