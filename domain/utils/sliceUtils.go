package utils

import "strings"

func ConvertStringSliceToMatrix(dna []string) [][]string {
	matrix := make([][]string, len(dna))
	for i, slice := range dna {
		matrix[i] = strings.Split(slice, "")
	}
	return matrix
}

func AllSameStrings(slice []string) bool {

	for i := 1; i < len(slice); i++ {
		if slice[i] != slice[0] {
			return false
		}
	}
	return true
}
