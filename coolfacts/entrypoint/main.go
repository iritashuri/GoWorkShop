package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/FTBpro/go-workshop/coolfacts/entrypoint/fact"
	"github.com/FTBpro/go-workshop/coolfacts/entrypoint/mentalfloss"
	"github.com/FTBpro/go-workshop/coolfacts/entrypoint/myhttp"
)

func main() {
	factsRepo := fact.Repository{
		Facts: []fact.Fact{},
	}

	handlerer := myhttp.FactsHandler{
		FactRepo: &factsRepo,
	}

	mentalfloss := mentalfloss.Mentalfloss{}
	facts, err := mentalfloss.Facts()
	if err != nil {
		fmt.Sprintf(`Error reading content %v `, err)
	}

	for _, fact := range facts {
		factsRepo.Add(fact)
	}

	http.HandleFunc("/ping", handlerer.Ping)
	http.HandleFunc("/facts", handlerer.Facts)
	log.Fatal(http.ListenAndServe(":9002", nil))
}
