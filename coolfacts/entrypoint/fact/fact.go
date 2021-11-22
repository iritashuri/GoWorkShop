package fact

type Fact struct {
	Image string
	Description string
}

type Repository struct {
	Facts []Fact
}

func (r *Repository) GetAll() []Fact {
	return r.Facts
}

func (r *Repository) Add(f Fact) {
	r.Facts = append(r.Facts, f)
}