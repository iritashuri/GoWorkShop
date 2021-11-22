package mentalfloss

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"my-go-work-shop/GoWorkShop/coolfacts/entrypoint/fact"
)

type Mentalfloss struct{}

func (mf Mentalfloss) Facts() ([]fact.Fact, error) {
	var facts []fact.Fact
	var items []struct {
		FactText     string `json:"fact"`
		PrimaryImage string `json:"primaryImage"`
	}

	resp, err := http.Get("http://mentalfloss.com/api/facts")
	if err != nil {
		fmt.Errorf(" `Error with get request! ,%v", err)
		return nil, err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf(` Error read boby respond! ,%v`, err)
	}

	err = json.Unmarshal(b, &items)
	if err != nil {
		return nil, fmt.Errorf(` Error unmarshel ,%v`, err)
	}

	for _, v := range items {
		f := fact.Fact{
			Image:       v.PrimaryImage,
			Description: v.FactText,
		}
		facts = append(facts, f)
	}

	defer resp.Body.Close()
	return facts, nil
}
