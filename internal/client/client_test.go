package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSucceeds(t *testing.T) {
	t.Parallel()

	body, err := Get("http://www:8080")
	assert.NoError(t, err)
	assert.Equal(t, "Hello, world!", body)
}

func TestGetFails(t *testing.T) {
	t.Parallel()

	body, err := Get("http://foobar:8080")
	assert.Error(t, err)
	assert.Equal(t, "", body)
}
