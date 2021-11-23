package providor

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"my-go-work-shop/GoWorkShop/coolfacts/entrypoint/fact"
)

type Provider struct{}

func NewProvider() *Provider {
	return &Provider{}
}

func (p Provider) Facts() (map[string]fact.Fact, error) {
	facts := make(map[string]fact.Fact)
	var items []struct {
		ID           string `json:"id"`
		FactText     string `json:"fact"`
		PrimaryImage string `json:"primaryImage"`
	}

	res, err := http.Get("http://mentalfloss.com/api/facts")
	if err != nil {
		fmt.Errorf(" `Error with get request! ,%v", err)
		return nil, err
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf(` Error read boby respond! ,%v`, err)
	}

	err = json.Unmarshal(b, &items)
	if err != nil {
		return nil, fmt.Errorf(` Error unmarshel ,%v`, err)
	}

	for _, v := range items {
		f := fact.Fact{
			ID:          v.ID,
			Image:       v.PrimaryImage,
			Description: v.FactText,
		}
		facts[f.ID] = f
	}

	defer res.Body.Close()
	return facts, nil
}
