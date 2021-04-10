package configuration

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetInjections(t *testing.T) {
	os.Setenv("AWS_ACCESS_KEY_ID","")
	os.Setenv("AWS_SECRET_ACCESS_KEY","")
	injects := GetInjections()
	assert.NotNil(t, injects)
}
