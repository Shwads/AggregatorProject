package scraper

import (
	"AggregatorProject/internal/database"
	"context"
	"fmt"
	"log"
	"sync"
)

func scrape(ctx context.Context, db *database.Queries) (errors map[string]error) {
	errors = make(map[string]error)

	feeds, getNextFeedsErr := db.GetNextFeedsToFetch(ctx, 10)
	if getNextFeedsErr != nil {
		log.Printf("Encountered error: %s. In function: Scrape.", getNextFeedsErr)
		errors["database error"] = getNextFeedsErr
		return
	}

	errorChannel := make(chan error, len(feeds))

	var wg sync.WaitGroup
	var mutex sync.Mutex

	for _, feed := range feeds {
		wg.Add(1)
		go func(f database.Feed) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				return
			default:
				page, scrapeFeedErr := sendRequest(ctx, f.Url)
				if scrapeFeedErr != nil {
					log.Printf("Encountered error: %s. In function: Scrape. Fetching feedname: %s", scrapeFeedErr, f.Name)
					errorChannel <- scrapeFeedErr
					mutex.Lock()
					if _, ok := errors[f.Name]; !ok {
						errors[f.Name] = scrapeFeedErr
					}
					mutex.Unlock()
					return
				}

				log.Printf("===== Scraped Posts: %s =====", feed.Name)
				for _, post := range page.Posts {
					log.Print(post.Title)
					fmt.Println()
				}
			}
		}(feed)
	}

	go func() {
		wg.Wait()
		close(errorChannel)
	}()

	if _, ok := <-errorChannel; ok {
		return
	}

	return nil
}
