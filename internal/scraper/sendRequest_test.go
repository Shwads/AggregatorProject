package scraper

import (
	"context"
	"log"
	"testing"
)

func TestSendRequest(t *testing.T) {
	page, fetchPageErr := sendRequest(context.Background(), "https://blog.boot.dev/index.xml")
	if fetchPageErr != nil {
		log.Printf("Encountered Error: %s. When testing function: sendRequest", fetchPageErr)
		return
	}

	for _, post := range page.Posts {
		log.Println()
		log.Printf("Title: %s", post.Title)
		log.Printf("Publication date: %s", post.PublicationDate)
		log.Printf("Link: %s", post.Link)
		log.Printf("Description: %s", post.Description)
		log.Println()
	}
}
