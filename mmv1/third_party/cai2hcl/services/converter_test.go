package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppendMap(t *testing.T) {
	m1 := map[string]string{"1": "a"}
	m2 := map[string]string{"2": "b"}
	appendMap(m1, m2)

	assert.Equal(t, m1, map[string]string{"1": "a", "2": "b"})
	assert.Equal(t, m2, map[string]string{"2": "b"})
}

func TestAppendMapConflict(t *testing.T) {
	assert.Panics(t, func() {
		m1 := map[string]string{"1": "a"}
		m2 := map[string]string{"1": "b"}
		appendMap(m1, m2)
	})
}
