package main

import (
	"context"
	"net/url"
	"os"
	"path/filepath"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/multierr"

	"github.com/gocolly/colly/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"github.com/wenerme/torrenti/pkg/scrape"
	"github.com/wenerme/torrenti/pkg/subi"
	"github.com/wenerme/torrenti/pkg/torrenti"
	"github.com/wenerme/torrenti/pkg/torrenti/util"
	"gorm.io/gorm"
)

func runScrape(cc *cli.Context) error {
	svc := newServeContext(cc)
	ctx, cancel := context.WithCancel(svc.Context)
	defer cancel()

	cache := filepath.Join(_conf.CacheDir, "web")
	if err := os.MkdirAll(cache, 0o755); err != nil {
		log.Fatal().Err(err).Msg("make cache dir")
	}
	log.Debug().Str("cache_dir", cache).Msg("cache dir")
	if cc.NArg() != 1 {
		log.Fatal().Msgf("must scrap only one url got %v", cc.NArg())
	}
	first := cc.Args().First()
	u, err := url.Parse(first)
	if err != nil {
		log.Fatal().Err(err).Msg("invalid url")
	}
	first = u.String()

	ctx = torrenti.IndexerContextKey.WithValue(ctx, getTorrentIndexer())
	ctx = subi.IndexerContextKey.WithValue(ctx, getSubIndexer())
	ctx = util.DirConfContextKey.WithValue(ctx, &_conf.DirConf)
	svc.Context = ctx

	c := colly.NewCollector(
		colly.CacheDir(cache),
		//colly.Debugger(&debug.WebDebugger{
		//	Address: _conf.Scrape.Debug.Addr,
		//}),
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.74 Safari/537.36"),
	)

	var sdb *gorm.DB
	_, sdb, err = newDB(&_conf.Scrape.Store.DB)
	if err != nil {
		return errors.Wrap(err, "new scraper store")
	}
	sc := &scrape.Context{
		Seed:      u,
		Fatal:     cc.Bool("fatal"),
		DB:        sdb,
		Context:   ctx,
		Collector: c,
	}

	registerDebug(svc)
	registerScrapeMetrics(sc)
	err = multierr.Combine(
		errors.Wrap(sc.Init(), "init scraper"),
		serveDebug(svc),
	)
	if err != nil {
		return err
	}

	if cc.Bool("seed") {
		return sc.Collector.Visit(first)
	}

	if !cc.Bool("pull") {
		if err = sc.Queue.AddURL(first); err != nil {
			return errors.Wrap(err, "add seed url to queue")
		}
	}

	svc.G.Add(func() error {
		return sc.Queue.Run(c)
	}, func(err error) {
		sc.Queue.Stop()
	})

	return svc.G.Run()
}

func registerScrapeMetrics(sc *scrape.Context) {
	promauto.NewGaugeFunc(prometheus.GaugeOpts{
		Name: "scrape_queue_size",
		Help: "Scrape queue size",
	}, func() float64 {
		size, err := sc.Queue.Size()
		if err != nil {
			log.Err(err).Msg("get scrape queue size")
		}
		return float64(size)
	})
	promauto.NewGaugeFunc(prometheus.GaugeOpts{
		Name: "scrape_request_count",
	}, func() float64 {
		return float64(sc.Stat.RequestCount)
	})
	promauto.NewGaugeFunc(prometheus.GaugeOpts{
		Name: "scrape_error_count",
	}, func() float64 {
		return float64(sc.Stat.ErrorCount)
	})
	promauto.NewGaugeFunc(prometheus.GaugeOpts{
		Name: "scrape_file_count",
	}, func() float64 {
		return float64(sc.Stat.FileCount)
	})
	promauto.NewGaugeFunc(prometheus.GaugeOpts{
		Name: "scrape_ext_total",
	}, func() float64 {
		return float64(sc.Stat.ExtensionCount)
	})
	promauto.NewGaugeFunc(prometheus.GaugeOpts{
		Name: "scrape_cache_hit_count",
	}, func() float64 {
		return float64(colly.CacheHit)
	})
	promauto.NewGaugeFunc(prometheus.GaugeOpts{
		Name: "scrape_cache_miss_count",
	}, func() float64 {
		return float64(colly.CacheMiss)
	})
	promauto.NewGaugeFunc(prometheus.GaugeOpts{
		Name: "scrape_cache_skip_count",
	}, func() float64 {
		return float64(colly.CacheSkip)
	})
	promauto.NewGaugeFunc(prometheus.GaugeOpts{
		Name: "scrape_cache_invalid_count",
	}, func() float64 {
		return float64(colly.CacheInvalid)
	})
}
