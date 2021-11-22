package main

import (
	"encoding/json"
	"fmt"
	"github.com/FTBpro/go-workshop/coolfacts/entrypoint/fact"
	"io/ioutil"
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

		switch r.Method {
		case "GET":
			allFacts := myFacts.GetAll()
			str, err := json.Marshal(allFacts)
			if err != nil {
				fmt.Errorf(" `Error! ,%v", err)
			}

			_, err = fmt.Fprint(w, string(str))
			if err != nil {
				http.Error(w, "Error", http.StatusInternalServerError)
			}
		case "POST":
			var req struct {
				Image       string `json:"image"`
				Description string `json:"description"`
			}

			if err := r.ParseForm(); err != nil {
				fmt.Fprintf(w, "ParseForm() err: %v", err)
				return
			}
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				fmt.Fprintf(w, "Error while reading erq body - %v", err)
				return
			}

			err = json.Unmarshal(b, &req)
			if err != nil {
				http.Error(w, `Error Unmarshel req `+string(err.Error()), http.StatusInternalServerError)
				return
			}
			k := fact.Fact{
				Image:       req.Image,
				Description: req.Description,
			}
			myFacts.Add(k)
			w.Write([]byte("SUCCESS"))

		default:
		}
	})
	log.Fatal(http.ListenAndServe(":9002", nil))
}
