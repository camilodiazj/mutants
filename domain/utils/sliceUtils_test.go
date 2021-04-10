package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConvertStringSliceToMatrix(t *testing.T) {
	expected := [][]string{{"A","B"},{"B","A"}}
	input := []string{"AB","BA"}
	result := ConvertStringSliceToMatrix(input)
	assert.Equal(t, expected, result)
}

func TestAllSameStrings(t *testing.T) {
	input := []string{"A","A","A","A"}
	result := AllSameStrings(input)
	assert.True(t, result)
}

func TestAllSameStringsShouldReturnFalse(t *testing.T) {
	input := []string{"A","C","A","A"}
	result := AllSameStrings(input)
	assert.False(t, result)
}
