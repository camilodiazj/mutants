package dna

type Entity struct {
	Dna      string
	Id       string
	IsMutant bool
}

type Counter struct {
	CountResult uint64
	TotalCount uint64
}

type Repository interface {
	Save(dna *Entity) error
	//TODO: Change CountMutants to Cont; Less specific
	CountMutants() (*Counter, error)
}
