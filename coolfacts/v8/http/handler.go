package http

import (
	"encoding/json"
	"fmt"
	"github.com/FTBpro/go-workshop/coolfacts/v8/fact"
	"html/template"
	"io/ioutil"
	"net/http"
)

type FactStore interface {
	Add(f fact.Fact)
	GetAll() []fact.Fact
}

type FactsHandler struct {
	FactStore FactStore
}

var newsTemplate = `
<html>
	<head>
		<title>Coolfacts</title>
	</head>
	<style>
body {
	font-family: Helvetica, Arial, sans-serif;
	color: #26323d;
  max-width: 720px;
  margin: auto;
}

article {
	border: 1px solid #0095c4;
	border-radius: 4px;
	max-width: 256px;
	text-align: center;
}

a {
	color: #26323d;
}
a:hover {
	color: #f16957;
}
img {
	border-radius: 4px;
}
	</style>
<body>
	<h1>Amazing Fact Generator</h1>
	<article>
		<a href="http://mentalfloss.com/api{{.Url}}">
				<h3>{{.Description}}</h3>
				<img src="{{.Image}}" width="100%" />
		</a>
	</article>
</body>
</html>`


func (h *FactsHandler) Ping(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "no http handler found", http.StatusNotFound)
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
	if h.FactStore == nil {
		http.Error(w, "fact store isn't initializes", http.StatusInternalServerError)
	}

	switch r.Method {
	case http.MethodGet:
		h.showFacts(w)
		return
	case http.MethodPost:
		h.postFacts(r, w)
		return
	default:
		http.Error(w, "no http handler found", http.StatusNotFound)
	}
}

func (h *FactsHandler) postFacts(r *http.Request, w http.ResponseWriter) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errMessage := fmt.Sprintf("error read from body: %v", err)
		http.Error(w, errMessage, http.StatusInternalServerError)
		return
	}
	var req struct {
		Image       string `json:"image"`
		Url         string `json:"url"`
		Description string `json:"description"`
	}
	err = json.Unmarshal(b, &req)
	if err != nil {
		errMessage := fmt.Sprintf("error parsing fact: %v", err)
		http.Error(w, errMessage, http.StatusBadRequest)
	}
	f := fact.Fact{
		Image:       req.Image,
		Url:         req.Url,
		Description: req.Description,
	}
	h.FactStore.Add(f)
	w.Write([]byte("SUCCESS"))
}

func (h *FactsHandler) showFacts(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "text/html")
	tmpl, err := template.New("facts").Parse(newsTemplate)
	if err != nil {
		errMessage := fmt.Sprintf("error ghttp template writing: %v", err)
		http.Error(w, errMessage, http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, h.FactStore.GetAll())
}