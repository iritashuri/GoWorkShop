package myhttp

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
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

type FactsHandler struct {
	FactRepo *fact.Repository
}

var req struct {
	Image       string `json:"Image"`
	Description string `json:"Description"`
}

func (h *FactsHandler) Ping(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "no myhttp handler found", http.StatusNotFound)
		return
	}

	w.Header().Add("Content-Type", "text/plain")

	_, err := fmt.Fprint(w, "PONG")
	if err != nil {
		errMessage := fmt.Sprintf("error writing response: %v", err)
		http.Error(w, errMessage, http.StatusInternalServerError)
	}
}

func (h *FactsHandler) Facts(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")

	switch r.Method {
	case http.MethodGet:
		tmpl, err := template.New("facts").Parse(newsTemplate)
		if err != nil {
			http.Error(w, `Error cresting and parsing html `+string(err.Error()), http.StatusInternalServerError)
			return
		}

		allFacts := h.FactRepo.GetAll()

		err = tmpl.Execute(w, allFacts)
		if err != nil {
			http.Error(w, `Error display content `+string(err.Error()), http.StatusInternalServerError)
			return
		}

	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			http.Error(w, `Error display content `+string(err.Error()), http.StatusInternalServerError)
			return
		}

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, `Error reading req body `+string(err.Error()), http.StatusInternalServerError)
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

		h.FactRepo.Add(k)

		w.Write([]byte("SUCCESS"))
	default:
	}
}
