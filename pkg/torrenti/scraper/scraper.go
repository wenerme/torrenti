package scraper

import (
	"context"
	"net/url"

	"github.com/wenerme/torrenti/pkg/torrenti/util"

	"github.com/gocolly/colly"
)

type ScrapeOptions struct {
	Seed  *url.URL
	Fatal bool
	Store *Store
}

var OptionContextKey = util.ContextKey[*ScrapeOptions]{Name: "ScrapeOptions"}

type Scraper struct {
	Name           string
	Support        func(ctx context.Context) bool
	InitContext    func(ctx context.Context) (context.Context, error)
	InitCollector  func(ctx context.Context, c *colly.Collector) error
	SetupCollector func(ctx context.Context, c *colly.Collector) error
}

var scrapers []*Scraper

func RegisterScraper(v *Scraper) {
	scrapers = append(scrapers, v)
}

func InitContext(ctx context.Context) (out context.Context, err error) {
	out = ctx
	err = support(ctx, func(s *Scraper) (e error) {
		if s.InitContext == nil {
			return nil
		}
		out, e = s.InitContext(ctx)
		return
	})
	return
}

func InitCollector(ctx context.Context) func(c *colly.Collector) error {
	return func(c *colly.Collector) error {
		return support(ctx, func(s *Scraper) error {
			if s.InitCollector != nil {
				return s.InitCollector(ctx, c)
			}
			return nil
		})
	}
}

func SetupCollector(ctx context.Context) func(c *colly.Collector) error {
	return func(c *colly.Collector) error {
		return support(ctx, func(s *Scraper) error {
			if s.SetupCollector == nil {
				return nil
			}
			return s.SetupCollector(ctx, c)
		})
	}
}

func support(c context.Context, cb func(*Scraper) error) error {
	for _, s := range scrapers {
		if s.Support(c) {
			if err := cb(s); err != nil {
				return err
			}
		}
	}
	return nil
}
