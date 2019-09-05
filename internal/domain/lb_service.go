package domain

type LbService struct {
	Name        string
	LbSrcRanges []Cidr
	UpdateOps   LbServicesUpdateOps
}

type LbServicesReadOps interface {
	FilterFor(l *LbSrcRanger) ([]LbService, error)
}

type LbServicesUpdateOps interface {
	UpdateCidrs(cidrs *[]Cidr) error
}
