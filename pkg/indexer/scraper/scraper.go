package scraper

import (
	"context"
	"net/url"

	"github.com/gocolly/colly"
)

type Scraper struct {
	Name    string
	Support func(u *url.URL) bool
	Init    func(ctx context.Context, c *colly.Collector)
	Setup   func(ctx context.Context, c *colly.Collector)
}

var scrapers []*Scraper

func RegisterScraper(v *Scraper) {
	scrapers = append(scrapers, v)
}

func InitCollector(ctx context.Context, u *url.URL) func(c *colly.Collector) {
	return func(c *colly.Collector) {
		for _, s := range scrapers {
			if s.Support(u) {
				s.Init(ctx, c)
			}
		}
	}
}

func SetupCollector(ctx context.Context, u *url.URL) func(c *colly.Collector) {
	return func(c *colly.Collector) {
		for _, s := range scrapers {
			if s.Support(u) {
				s.Setup(ctx, c)
			}
		}
	}
}
