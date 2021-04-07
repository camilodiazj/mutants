package mutant

import (
	"regexp"
	"strings"
)

type Dna struct {
	Sequence []string `json:"dna"`
}

func isValidDna(dna []string) bool {
	validBases := regexp.MustCompile(`^[ATCG]+$`).MatchString
	expectedSize := len(dna) * len(dna)
	nitrogenBases := strings.Join(dna, "")
	return len(nitrogenBases) == expectedSize && validBases(nitrogenBases)
}
