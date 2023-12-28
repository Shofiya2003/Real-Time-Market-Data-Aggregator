package scraper

import (
	"fmt"
	"log"

	"github.com/gocolly/colly/v2"
)

type Marketcap_Scraper struct {
	Url string
}

type Currency struct {
	name       string
	symbol     string
	market_cap string
}

func (scraper Marketcap_Scraper) Scrape() error {
	c := colly.NewCollector(colly.UserAgent("Mozilla/5.0 (X11; Linux i686; rv:10.0) Gecko/20100101 Firefox/10.0"))
	names := []Currency{}
	c.OnHTML(".cmc-table-row", func(h *colly.HTMLElement) {
		new_currency := Currency{}
		h.ForEach("td", func(i int, h *colly.HTMLElement) {
			if i == 0 {
				new_currency.name = h.Text
			}
			if i == 1 {
				new_currency.symbol = h.Text
			}

			if i == 2 {
				new_currency.market_cap = h.Text
			}

		})
		names = append(names, new_currency)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Fatal("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.Visit(scraper.Url)

	fmt.Printf("size of the data is %d", len(names))
	return nil
}
