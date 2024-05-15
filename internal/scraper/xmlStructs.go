package scraper

type Page struct {
	Posts []Post `xml:"channel>item"`
}

type Post struct {
	Title           string `xml:"title"`
	Link            string `xml:"link"`
	PublicationDate string `xml:"pubDate"`
	Description     string `xml:"description"`
}
