package scraper

type Scraper interface {
	Scrape() error
}

func main() {
	// scrapers := []Scraper{Marketcap_Scraper{"https://coinmarketcap.com/all/views/all/"}}

}
