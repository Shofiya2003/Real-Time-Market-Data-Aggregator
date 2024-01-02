package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gocolly/colly/v2"
)

type Marketcap_Scraper struct {
	Url string
}

func (scraper Marketcap_Scraper) Scrape() error {
	c := colly.NewCollector(colly.UserAgent("Mozilla/5.0 (X11; Linux i686; rv:10.0) Gecko/20100101 Firefox/10.0"))
	names := []Currency{}
	c.OnHTML(".cmc-table-row", func(h *colly.HTMLElement) {

		fields := make([]string, 10)

		h.ForEach("td", func(i int, h *colly.HTMLElement) {
			if i != 0 && i != 9 {
				fields[i-1] = h.Text
			}
		})
		new_currency := Currency{fields[0], fields[1], fields[2], fields[3], fields[4], fields[5], fields[6], fields[7], fields[8], time.Now().String()}
		err := StoreCurrency(RedisClient, new_currency)
		if err != nil {
			fmt.Printf("error in storing currency %s : %s", new_currency.Name, err)
		}
		fmt.Println("moving to next currecny")
		// var curr Currency
		// curr, err = GetCurrency(RedisClient, new_currency.Name)
		// if err != nil {
		// 	fmt.Printf("error in retreiving currency %s : %s", new_currency.Name, err)
		// }

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
