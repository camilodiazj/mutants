package mutant

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsValidDna(t *testing.T) {
	input := []string{"ATGCGA","CAGTGC","TTATGT","AGAAGG","CCCCTA","TCACGG"}
	result := IsValidDna(input)
	assert.True(t, result)
}

func TestIsValidDnaShouldReturnFalse(t *testing.T) {
	input := []string{"ATGCGA","CAGTGC","TTATGT","AGAAGG","CCCCTA","TCACGX"}
	result := IsValidDna(input)
	assert.False(t, result)
}

func TestIsValidDnaShouldReturnFalseDueInvalidSize(t *testing.T) {
	input := []string{"ATGCGA","CAGTGC","TTATGT","AGAAGG","CCCCTA","TCACGCC"}
	result := IsValidDna(input)
	assert.False(t, result)
}
