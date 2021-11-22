package mentalfloss

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/FTBpro/go-workshop/coolfacts/entrypoint/fact"
)

type Mentalfloss struct{}

func (mf Mentalfloss) Facts() ([]fact.Fact, error) {
	var facts []fact.Fact
	var items []struct {
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
		return nil, fmt.Errorf(` Error read body respond! ,%v`, err)
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

	defer res.Body.Close()
	return facts, nil
}
