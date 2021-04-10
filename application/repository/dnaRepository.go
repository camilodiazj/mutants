package repository

type DnaEntity struct {
	Dna      string
	Id       string
	IsMutant bool
}

type Count struct {
	CountResult uint64
	TotalCount uint64
}

type DnaRepository interface {
	Save(dna *DnaEntity) error
	CountMutants() (*Count, error)
}
