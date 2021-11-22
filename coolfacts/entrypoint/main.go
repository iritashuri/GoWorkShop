package main

import (
	"encoding/json"
	"fmt"
	"github.com/FTBpro/go-workshop/coolfacts/entrypoint/fact"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, "PONG\n")
		if err != nil {
			http.Error(w, "Error", http.StatusInternalServerError)
		}
	})

	myFacts := fact.Repository{
		Facts: []fact.Fact{
			{
				Image:       "onePic",
				Description: "oneDes",
			},
			{
				Image:       "twoPic",
				Description: "twoDes",
			},
			{
				Image:       "thirdPic",
				Description: "thirdDes",
			},
			{
				Image:       "foursPic",
				Description: "foursDes",
			},
		},
	}

	http.HandleFunc("/facts", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		allFacts := myFacts.GetAll()
		str, err := json.Marshal(allFacts)
		if err != nil {
			fmt.Errorf(" `Error! ,%v", err)
		}

		_, err = fmt.Fprint(w, string(str))
		if err != nil {
			http.Error(w, "Error", http.StatusInternalServerError)
		}
	})
	log.Fatal(http.ListenAndServe(":9002", nil))
}
