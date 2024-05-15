package scraper

import (
	"context"
	"encoding/xml"
	"io"
	"log"
	"net/http"
)

func sendRequest(ctx context.Context, url string) (Page, error) {

	var page Page

	req, generateRequestErr := http.NewRequestWithContext(ctx, "GET", url, nil)
	if generateRequestErr != nil {
		log.Printf("Encountered error: %s in function: sendRequest", generateRequestErr)
		return page, generateRequestErr
	}

	client := &http.Client{}

	res, clientRequestErr := client.Do(req)
	if clientRequestErr != nil {
		log.Printf("Encountered error: %s. in function: sendRequest", clientRequestErr)
		return page, clientRequestErr
	}
	defer res.Body.Close()

	data, responseBodyReadErr := io.ReadAll(res.Body)
	if responseBodyReadErr != nil {
		log.Printf("Encountered error: %s. in function: sendRequest", responseBodyReadErr)
		return page, responseBodyReadErr
	}

	xmlUnmarshalErr := xml.Unmarshal(data, &page)
	if xmlUnmarshalErr != nil {
		log.Printf("Encountered error: %s. in function: sendRequest", xmlUnmarshalErr)
		return page, xmlUnmarshalErr
	}

	return page, nil
}
