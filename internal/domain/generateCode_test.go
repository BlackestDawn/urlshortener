package domain

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateCode_HasCorrectLength(t *testing.T) {
	url := "https://www.google.com"
	want := CodeLength

	result, _ := GenerateCode(url)

	assert.Equal(t, want, len(result))
}

func TestGenerateCode_IsUrlSafe(t *testing.T) {
	url := "https://www.google.com"

	got, _ := GenerateCode(url)

	assert.Regexp(t, regexp.MustCompile("^[a-zA-Z0-9]*$"), got)
}

func TestGenerateCode_Uniqueness(t *testing.T) {
	baseUrl := "https://www.google.com/"
	iters := 10000

	got := make([]string, iters)

	for i := range iters {
		code, _ := GenerateCode(fmt.Sprintf("%s%d", baseUrl, i))
		got[i] = code
		if i < 1 {
			continue
		}
		for j := range i - 1 {
			assert.NotEqual(t, code, got[j])
		}
	}
}
