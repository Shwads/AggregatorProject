package scraper

import (
	"AggregatorProject/internal/database"
	"context"
	"log"
	"time"
)

func ScraperWorker(ctx context.Context, db *database.Queries) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Shutting down worker")
			return
		default:
			scraperErrors := scrape(ctx, db)

			if len(scraperErrors) > 0 {
				for key, value := range scraperErrors {
					log.Printf("Encountered error: %s", value)
					log.Printf("\tWhen scraping feed: %s", key)
					log.Println()
				}
			}
			time.Sleep(60 * time.Second)
		}
	}
}
