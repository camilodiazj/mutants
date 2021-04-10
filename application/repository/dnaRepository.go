package repository

type DnaEntity struct {
	Dna      string
	Id       string
	IsMutant bool
}

type Counter struct {
	CountResult uint64
	TotalCount uint64
}

type DnaRepository interface {
	Save(dna *DnaEntity) error
	CountMutants() (*Counter, error)
}
