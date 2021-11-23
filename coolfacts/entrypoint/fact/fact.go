package fact

import "fmt"

type Fact struct {
	Image       string
	Description string
}

type Provider interface {
	Facts() ([]Fact, error)
}

type Repository interface {
	Add(f Fact)
	GetAll() []Fact
}

type service struct {
	provider Provider
	repository    Repository
}

func (s *service) UpdateFacts() func() error {
	return func() error {
		facts, err := s.provider.Facts()
		if err != nil {
			return fmt.Errorf(`Error reading content %v `, err)
		}

		for _, fact := range facts {
			s.repository.Add(fact)
		}
		return nil
	}
}

func NewService(r Repository, p Provider) *service {
	return &service{
		provider:   p,
		repository: r,
	}
}

