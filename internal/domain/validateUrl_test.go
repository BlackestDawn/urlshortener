package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateURL_AcceptsHttps(t *testing.T) {
	url := "https://www.google.com"

	got, _ := ValidateURL(url)

	assert.True(t, got)
}

func TestValidateURL_AcceptsHttp(t *testing.T) {
	url := "http://www.google.com"

	got, _ := ValidateURL(url)

	assert.True(t, got)
}

func TestValidateURL_RejectsEmpty(t *testing.T) {
	url := ""

	got, _ := ValidateURL(url)

	assert.False(t, got)
}

func TestValidateURL_RejectsMalformed(t *testing.T) {
	url := "not-a-url"

	got, _ := ValidateURL(url)

	assert.False(t, got)
}

func TestValidateURL_RejectsNoScheme(t *testing.T) {
	url := "www.google.com"

	got, _ := ValidateURL(url)

	assert.False(t, got)
}
