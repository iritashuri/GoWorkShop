package main

import (
	"encoding/json"
	"fmt"
	"github.com/FTBpro/go-workshop/coolfacts/entrypoint/fact"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

var newsTemplate = `<!DOCTYPE html>
<html>
  <head><style>/* copy coolfacts/styles.css for some color 🎨*/</style></head>
  <body>
  <h1>Facts List</h1>
  <div>
    {{ range . }}
       <article>
            <h3>{{.Description}}</h3>
            <img src="{{.Image}}" width="30%" />
       </article>
    {{ end }}
  <div>
  </body>
</html>`

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
		w.Header().Add("Content-Type", "text/html")

		switch r.Method {
		case http.MethodGet:
			w.Header().Add("Content-Type", "text/html")
			tmpl, err := template.New("facts").Parse(newsTemplate)
			if err != nil {
				http.Error(w, `Error cresting and parsing html `+ string(err.Error()), http.StatusInternalServerError)
				return
			}

			allFacts := myFacts.GetAll()
			err = tmpl.Execute(w, allFacts)
			if err != nil {
				http.Error(w, `Error executing `+ string(err.Error()), http.StatusInternalServerError)
				return
			}

		case http.MethodPost:
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
