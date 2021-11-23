package myhttp

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"

	"my-go-work-shop/GoWorkShop/coolfacts/entrypoint/fact"
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

var req struct {
	Id          string `json:"id"`
	Image       string `json:"Image"`
	Description string `json:"Description"`
}

type FactsHandler struct {
	FactRepo FactRepository
}

type FactRepository interface {
	Add(f fact.Fact)
	GetAll() []fact.Fact
}

func NewFactsHandler(factRepo FactRepository) *FactsHandler {
	return &FactsHandler{
		FactRepo: factRepo,
	}
}

func (h *FactsHandler) Ping(w http.ResponseWriter, r *http.Request) {
	showPong(w, r)
}

func (h *FactsHandler) Facts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		showFacts(h, w)

	case http.MethodPost:
		postFacts(h, w, r)
	default:
	}
}

func postFacts(h *FactsHandler, w http.ResponseWriter, r *http.Request) {
	fmt.Println("got post")
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
}

func showFacts(h *FactsHandler, w http.ResponseWriter) {
	w.Header().Add("Content-Type", "text/html")
	tmpl, err := template.New("facts").Parse(newsTemplate)
	if err != nil {
		http.Error(w, `Error cresting and parsing html `+string(err.Error()), http.StatusInternalServerError)
		return
	}

	allFacts := h.FactRepo.GetAll()
	tmpl.Execute(w, allFacts)
}

func showPong(w http.ResponseWriter, r *http.Request) {
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