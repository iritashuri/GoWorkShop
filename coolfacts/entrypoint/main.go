package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"my-go-work-shop/GoWorkShop/coolfacts/entrypoint/fact"
	"my-go-work-shop/GoWorkShop/coolfacts/entrypoint/inmem"
	"my-go-work-shop/GoWorkShop/coolfacts/entrypoint/myhttp"
	"my-go-work-shop/GoWorkShop/coolfacts/entrypoint/providor"
)

func main() {
	factsRepo := inmem.NewFactRepository()

	handlerer := myhttp.NewFactsHandler(factsRepo)

	provider := providor.NewProvider()
	service := fact.NewService(factsRepo, provider)

	updater := service.UpdateFacts()
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
