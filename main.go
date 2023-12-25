package main

type Scraper interface {
	Scrape() error
}

func main() {
	scrapers := []Scraper{Marketcap_Scraper{"https://coinmarketcap.com/all/views/all/"}}
	for i := 0; i < len(scrapers); i++ {
		scrapers[i].Scrape()
	}
}
