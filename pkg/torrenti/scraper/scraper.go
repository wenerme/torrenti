package scraper

import (
	"context"
	"net/url"

	"github.com/wenerme/torrenti/pkg/torrenti/util"

	"github.com/rs/zerolog/log"

	"github.com/gocolly/colly"
)

type ScrapeOptions struct {
	Fatal bool
}

var OptionContextKey = util.ContextKey[*ScrapeOptions]{"ScrapeOptions"}

type Scraper struct {
	Name    string
	Support func(u *url.URL) bool
	Init    func(ctx context.Context, c *colly.Collector) error
	Setup   func(ctx context.Context, c *colly.Collector) error
}

var scrapers []*Scraper

func RegisterScraper(v *Scraper) {
	scrapers = append(scrapers, v)
}

func InitCollector(ctx context.Context, u *url.URL) func(c *colly.Collector) {
	return func(c *colly.Collector) {
		for _, s := range scrapers {
			if s.Support(u) {
				if err := s.Init(ctx, c); err != nil {
					log.Fatal().Err(err).Str("name", s.Name).Msg("init scraper failed")
				}
			}
		}
	}
}

func SetupCollector(ctx context.Context, u *url.URL) func(c *colly.Collector) {
	return func(c *colly.Collector) {
		for _, s := range scrapers {
			if s.Support(u) {
				if err := s.Setup(ctx, c); err != nil {
					log.Fatal().Err(err).Str("name", s.Name).Msg("setup scraper failed")
				}
			}
		}
	}
}
