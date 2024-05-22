package scraper

import (
	"AggregatorProject/internal/database"
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

func scrape(ctx context.Context, db *database.Queries) (errors map[string]error) {
	// store any feed errors in a map object so that we can easily see which feeds failed
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

				postCount := 0

				for _, post := range page.Posts {
					log.Printf("%s\n", post.Title)
					postCount++

					postUUID := uuid.New()

					layout := "Mon, 02 Jan 2006 15:04:05 -0700"
					pubDate, parseErr := time.Parse(layout, post.PublicationDate)

					dbPost := database.CreatePostParams{
						ID:        postUUID,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
						Title:     post.Title,
						Url:       post.Link,
						FeedID:    feed.ID,
					}

					if parseErr == nil {
						dbPost.PublishedAt = pubDate
					}

					if len(post.Description) > 0 {
						desc := sql.NullString{
							Valid:  true,
							String: post.Description,
						}
						dbPost.Description = desc
					}

					_, createPostErr := db.CreatePost(ctx, dbPost)
					if createPostErr != nil {
						if parsedErr, ok := createPostErr.(*pq.Error); ok {
							if parsedErr.Code.Name() != "unique_violation" {
								log.Printf("Encountered error: %s", parsedErr)
							}
						}
					}
				}
				now := time.Now()
				nullTime := sql.NullTime{
					Time:  now,
					Valid: true,
				}
				fetchedParams := database.MarkFeedFetchedParams{
					ID:            feed.ID,
					LastFetchedAt: nullTime,
					UpdatedAt:     now,
				}
				_, markFetchedErr := db.MarkFeedFetched(ctx, fetchedParams)
				if markFetchedErr != nil && markFetchedErr != context.Canceled {
					log.Printf("Encountered Error: %s. When marking %s as fetched in function: scraper.", markFetchedErr, f.Name)
					return
				}

				log.Printf("\n=====================================================================================\n" +
					fmt.Sprintf("Successfully updated feed %s to fetched at %v\n", f.Name, now) +
					fmt.Sprintf("Fetched %v posts from feed %s\n", postCount, f.Name) +
					"=====================================================================================\n")
			}
		}(feed)
	}

	go func() {
		wg.Wait()
		close(errorChannel)
	}()

	for err := range errorChannel {
		if err != nil {
			return
		}
	}

	return nil
}
