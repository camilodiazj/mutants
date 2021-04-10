package mutant

import (
	"errors"
	"github.com/camilodiazj/mutants/domain/utils"
	"log"
)

type Direction int

const (
	HOR Direction = iota
	VER
)

const (
	expectedLength = 4
	invalidBase    = "X"
)

type Line struct {
	Positions  map[int]int
	Line       []string
	IsSequence bool
}

type MutanVerifier interface {
	IsMutant(dna []string) (bool, error)
}

type mutanService struct{}

func NewMutanVerifier() MutanVerifier {
	return &mutanService{}
}

func (mutanService) IsMutant(dna []string) (bool, error) {
	if !IsValidDna(dna) {
		log.Println("Fail due invalid DNA input")
		return false, errors.New("invalid DNA")
	}
	matrix := utils.ConvertStringSliceToMatrix(dna)
	sequenceFound := 0
	indexLimit := len(dna) - expectedLength
	for row, slice := range matrix {
		for column, base := range slice {
			if sequenceFound > 1 {
				return true, nil
			}
			if base == invalidBase {
				continue
			}
			rightPositionsAvailable := column <= indexLimit
			downPositionsAvailable := row <= indexLimit
			leftPositionsAvailable := column >= expectedLength-1
			if rightPositionsAvailable && dna[row][column] == dna[row][column+1] && horizontalValidation(matrix, row, column) {
				sequenceFound++
			} else if downPositionsAvailable && dna[row][column] == dna[row+1][column] && verticalValidation(matrix, row, column) {
				sequenceFound++
			} else if rightPositionsAvailable && downPositionsAvailable && dna[row][column] == dna[row+1][column+1] &&
				obliqueLeftToRightValidation(matrix, row, column) {
				sequenceFound++
			} else if downPositionsAvailable && leftPositionsAvailable && dna[row][column] == dna[row+1][column-1] &&
				obliqueRightToLeftValidation(matrix, row, column) {
				sequenceFound++
			}
		}
	}
	return false, nil
}

func horizontalValidation(dna [][]string, row int, column int) bool {
	line := dna[row][column : column+expectedLength]
	if utils.AllSameStrings(line) {
		return disableValues(dna, row, column, HOR)
	}
	return false
}

func verticalValidation(dna [][]string, row int, column int) bool {
	line := []string{dna[row][column], dna[row+1][column], dna[row+2][column], dna[row+3][column]}
	if utils.AllSameStrings(line) {
		return disableValues(dna, row, column, VER)
	}
	return false
}

func obliqueLeftToRightValidation(dna [][]string, row int, column int) bool {
	line := []string{dna[row][column], dna[row+1][column+1], dna[row+2][column+2], dna[row+3][column+3]}
	if utils.AllSameStrings(line) {
		positions := map[int]int{row: column, row + 1: column + 1, row + 2: column + 2, row + 3: column + 3}
		return disableObliquesValues(dna, positions)
	}
	return false
}

func obliqueRightToLeftValidation(dna [][]string, row int, column int) bool {
	line := []string{dna[row][column], dna[row+1][column-1], dna[row+2][column-2], dna[row+3][column-3]}
	if utils.AllSameStrings(line) {
		positions := map[int]int{row: column, row + 1: column - 1, row + 2: column - 2, row + 3: column - 3}
		return disableObliquesValues(dna, positions)
	}
	return false
}

func disableObliquesValues(dna [][]string, positions map[int]int) bool {
	for i, k := range positions {
		dna[i][k] = invalidBase
	}
	return true
}

func disableValues(dna [][]string, row int, column int, direction Direction) bool {
	dna[row][column] = invalidBase
	switch direction {
	case HOR:
		dna[row][column+1] = invalidBase
		dna[row][column+2] = invalidBase
		dna[row][column+3] = invalidBase
	case VER:
		dna[row+1][column] = invalidBase
		dna[row+2][column] = invalidBase
		dna[row+3][column] = invalidBase
	}
	return true
}
