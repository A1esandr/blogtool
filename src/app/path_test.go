package app

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathConfig_compose(t *testing.T) {
	c := &pathConfig{}
	result := c.compose("", "http://example.com")
	assert.True(t, result != "")
	assert.True(t, strings.HasPrefix(result, "example.com/"))
}
