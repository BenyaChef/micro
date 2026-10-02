package passwordhasher

import "golang.org/x/crypto/bcrypt"

type Builder struct {
	cost int
}

func NewBuilder() *Builder {
	return &Builder{cost: bcrypt.DefaultCost}
}

func (b *Builder) Cost(cost int) *Builder {
	b.cost = cost

	return b
}

func (b *Builder) Build() (*Hasher, error) {
	if b.cost < bcrypt.MinCost || b.cost > bcrypt.MaxCost {
		return nil, ErrCostOutOfRange(b.cost, bcrypt.MinCost, bcrypt.MaxCost)
	}

	return &Hasher{cost: b.cost}, nil
}
