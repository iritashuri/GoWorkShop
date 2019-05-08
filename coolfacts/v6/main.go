package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/FTBpro/go-workshop/coolfacts/v6/fact"
	"github.com/FTBpro/go-workshop/coolfacts/v6/mentalfloss"
)

type Handlerer struct {
	store *fact.Store
}

var newsTemplate = `<html>
                    <h1>Facts</h1>
                    <div>
                        {{range .}}
                            <div>
                                <h3>{{.Description}}</h3>
                                <img src="{{.Image}}" width="25%" height="25%"> </img>
                            </div>
                        {{end}}
                    <div>
                    </html>`

func main() {
	factsStore := fact.Store{}
	mf := mentalfloss.Mentalfloss{}
	handlerer := &Handlerer{&factsStore}

	err := fillStoreWithNewData(mf, &factsStore)

	http.HandleFunc("/ping", handlerer.Ping)
	http.HandleFunc("/facts", handlerer.Facts)

	err = http.ListenAndServe(":9002", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func fillStoreWithNewData(mf mentalfloss.Mentalfloss, factsStore *fact.Store) error {
	facts, err := mf.Facts()
	if err != nil {
		log.Fatal("can't reach mentalfloss: ", err)
	}
	for _, f := range facts {
		factsStore.Add(f)
	}
	return err
}

func (h *Handlerer) Ping(w http.ResponseWriter, r *http.Request) {
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

func (h *Handlerer) Facts(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.getFacts(w)
		return
	}
	if r.Method == http.MethodPost {
		h.postFacts(r, w)
		return
	}
	http.Error(w, "no http handler found", http.StatusNotFound)
}

func (h *Handlerer) postFacts(r *http.Request, w http.ResponseWriter) {
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
	h.store.Add(f)
	w.Write([]byte("SUCCESS"))
}

func (h *Handlerer) getFacts(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "text/html")
	tmpl, err := template.New("facts").Parse(newsTemplate)
	if err != nil {
		errMessage := fmt.Sprintf("error ghttp template writing: %v", err)
		http.Error(w, errMessage, http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, h.store.GetAll())
}