package inmem

import "my-go-work-shop/GoWorkShop/coolfacts/entrypoint/fact"

type factRepository struct {
	facts []fact.Fact
}

func NewFactRepository() *factRepository {
	return &factRepository{}
}

func (r *factRepository) Add(f fact.Fact) {
	r.facts = append(r.facts, f)
}

func (r *factRepository) GetAll() []fact.Fact {
	return r.facts
}