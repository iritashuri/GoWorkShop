package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/FTBpro/go-workshop/coolfacts/entrypoint/fact"
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
	myFacts := fact.Repository{
		Facts: []fact.Fact{},
	}

	mentalfloss := Mentalfloss{}
	facts, err := mentalfloss.Facts()
	if err != nil {
		fmt.Sprintf(`Error reading content %v `, err)
	}

	for _, fact := range facts {
		myFacts.Add(fact)
	}

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, "PONG\n")
		if err != nil {
			http.Error(w, "Error", http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/facts", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")

		switch r.Method {
		case http.MethodGet:
			w.Header().Add("Content-Type", "text/html")

			tmpl, err := template.New("facts").Parse(newsTemplate)
			if err != nil {
				http.Error(w, `Error cresting and parsing html `+string(err.Error()), http.StatusInternalServerError)
				return
			}

			err = tmpl.Execute(w, facts)
			if err != nil {
				http.Error(w, `Error executing `+string(err.Error()), http.StatusInternalServerError)
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
