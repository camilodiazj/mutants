package mutant

type Direction int

const (
	HOR Direction = iota
	VER
	OLR
	ORL
)

const (
	expectedLength = 4
	invalidBase    = "X"
)

type Line struct {
	Positions map[int]int
	Line      []string
}

type MutanService interface {
	IsMutant(dna []string) bool
}

type mutanService struct{}

func NewMutanService() MutanService {
	return &mutanService{}
}

func (mutanService) IsMutant(dna []string) bool {
	if !isValidDna(dna) {
		return false
	}
	matrix := ConvertStringSliceToMatrix(dna)

	sequenceFound := 0

	indexLimit := len(dna) - expectedLength
	for row, slice := range matrix {
		for column, base := range slice {
			if sequenceFound > 1 {
				return true
			}
			if base == invalidBase {
				continue
			}
			rightPositionsAvailable := column <= indexLimit
			downPositionsAvailable := row <= indexLimit
			leftPositionsAvailable := column >= expectedLength-1
			if rightPositionsAvailable && horizontalValidation(matrix, row, column) {
				sequenceFound++
			} else if downPositionsAvailable && verticalValidation(matrix, row, column) {
				sequenceFound++
			} else if rightPositionsAvailable && downPositionsAvailable && obliqueLeftToRightValidation(matrix, row, column) {
				sequenceFound++
			} else if downPositionsAvailable && leftPositionsAvailable && obliqueRightToLeftValidation(matrix, row, column) {
				sequenceFound++
			}
		}
	}
	return false
}

func horizontalValidation(dna [][]string, row int, column int) bool {
	fourSizeLine := dna[row][column : column+expectedLength]
	return validateLine(dna, fourSizeLine, row, column, HOR)
}

func verticalValidation(dna [][]string, row int, column int) bool {
	line := getFourSizeLine(dna, VER, row, column)
	return validateLine(dna, line.Line, row, column, VER)
}

func obliqueLeftToRightValidation(dna [][]string, row int, column int) bool {
	line := getFourSizeLine(dna, OLR, row, column)
	return isObliqueSequence(dna, line.Line, line.Positions)
}

func obliqueRightToLeftValidation(dna [][]string, row int, column int) bool {
	line := getFourSizeLine(dna, ORL, row, column)
	return isObliqueSequence(dna, line.Line, line.Positions)
}

func getFourSizeLine(dna [][]string, direction Direction, row int, column int) Line {
	var line []string
	var positions map[int]int

	switch direction {
	case VER:
		line = []string{dna[row][column], dna[row+1][column], dna[row+2][column], dna[row+3][column]}
	case OLR:
		line = []string{dna[row][column], dna[row+1][column+1], dna[row+2][column+2], dna[row+3][column+3]}
		positions = map[int]int{row: column, row + 1: column + 1, row + 2: column + 2, row + 3: column + 3}
	case ORL:
		line = []string{dna[row][column], dna[row+1][column-1], dna[row+2][column-2], dna[row+3][column-3]}
		positions = map[int]int{row: column, row + 1: column - 1, row + 2: column - 2, row + 3: column - 3}
	}

	return Line{
		Positions: positions,
		Line:      line,
	}
}

func validateLine(dna [][]string, line []string, row int, column int, direction Direction) bool {
	if AllSameStrings(line) {
		switch direction {
		case VER:
			return disableValues(dna, row, column, false)
		case HOR:
			return disableValues(dna, column, row, true)
		}
	}
	return false
}

func isObliqueSequence(dna [][]string, line []string, positions map[int]int) bool {
	if AllSameStrings(line) {
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

func disableValues(dna [][]string, init int, index int, isStaticRow bool) bool {
	for i := init; i < init+expectedLength; i++ {
		if isStaticRow {
			dna[index][i] = invalidBase
		} else {
			dna[i][index] = invalidBase
		}
	}
	return true
}
