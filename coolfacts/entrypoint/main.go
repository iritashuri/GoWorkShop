package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"my-go-work-shop/GoWorkShop/coolfacts/entrypoint/fact"
	"my-go-work-shop/GoWorkShop/coolfacts/entrypoint/mentalfloss"
	"my-go-work-shop/GoWorkShop/coolfacts/entrypoint/myhttp"
)

func main() {
	factsRepo := fact.Repository{}
	handlerer := myhttp.FactsHandler{
		FactRepo: &factsRepo,
	}

	mentalFloss := mentalfloss.Mentalfloss{}

	updateFunc(mentalFloss, &factsRepo)

	updater := updateFunc(mentalFloss, &factsRepo)
	if err := updater(); err != nil {
		log.Printf(`errror in updateFunc: %v`, err)
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	updateFactsWithTicker(ctx, updater)

	http.HandleFunc("/ping", handlerer.Ping)
	http.HandleFunc("/facts", handlerer.Facts)
	log.Fatal(http.ListenAndServe(":9002", nil))
}

func updateFactsWithTicker(ctx context.Context, updateFunc func() error) {
	ticker := time.NewTicker(3 * time.Second)

	go func(ctx context.Context) {
		for {
			select {
			case <-ticker.C:
				fmt.Println("updating")
				if err := updateFunc(); err != nil {
					log.Printf(`errror in updateFunc: %v`, err)
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}(ctx)
}

func updateFunc(mentalFloss mentalfloss.Mentalfloss, repo *fact.Repository) func() error {
	return func() error {
		facts, err := mentalFloss.Facts()
		if err != nil {
			return fmt.Errorf(`Error reading content %v `, err)
		}

		for _, fact := range facts {
			repo.Add(fact)
		}
		return nil
	}
}
