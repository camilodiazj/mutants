package mutant

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsMutant(t *testing.T) {
	input := []string{"ATGCGA", "CAGTGC", "TTATGT", "AGAAGG", "CCCCTA", "TCACGG"}
	verifier := NewMutanVerifier()
	result := verifier.IsMutant(input)
	assert.True(t, result)
}

func TestIsMutantShouldReturnFalse(t *testing.T) {
	input := []string{"ATGCGA", "CCGTTC", "TTATGT", "AGAAGG", "CCCCTA", "TCACGG"}
	verifier := NewMutanVerifier()
	result := verifier.IsMutant(input)
	assert.False(t, result)
}
