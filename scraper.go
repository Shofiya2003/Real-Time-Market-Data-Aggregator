package main

type Scraper interface {
	Scrape() error
}

type Currency struct {
	Name                  string
	Symbol                string
	MarketCap             string
	Price                 string
	CirculatingSupply     string
	Volume                string
	OneHourChange         string
	TwentyFourHoursChange string
	SevenDaysChange       string
	UpdatedAt             string
}
