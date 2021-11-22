package fact

type Fact struct {
	Image       string
	Description string
}

type Repository struct {
	facts []Fact
}

func (r *Repository) GetAll() []Fact {
	return r.facts
}

func (r *Repository) Add(f Fact) {
	r.facts = append(r.facts, f)
}
