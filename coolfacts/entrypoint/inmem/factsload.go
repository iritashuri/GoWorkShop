package inmem

import "my-go-work-shop/GoWorkShop/coolfacts/entrypoint/fact"

type factRepository struct {
	facts map[string]fact.Fact
}

func NewFactRepository() *factRepository {
	return &factRepository{}
}

func (r *factRepository) Add(f fact.Fact) {
	r.facts[f.ID] = f
}

func (r *factRepository) GetAll() map[string]fact.Fact {
	return r.facts
}
