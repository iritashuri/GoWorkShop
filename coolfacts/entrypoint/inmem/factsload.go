package inmem

import "my-go-work-shop/GoWorkShop/coolfacts/entrypoint/fact"

type factRepository struct {
	facts map[string]fact.Fact
}

func NewFactRepository() *factRepository {
	return &factRepository{map[string]fact.Fact{}}
}

func (r *factRepository) Add(f fact.Fact) {
	r.facts[f.ID] = f
}

func (r *factRepository) GetAll() []fact.Fact {
	var facts []fact.Fact
	for  _, v := range r.facts{
		facts = append(facts, v)
	}
	return facts
}
