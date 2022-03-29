package main

import (
	"context"
	"net/url"
	"os"
	"path/filepath"

	"github.com/gocolly/colly/v2/debug"

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

	var sdb *gorm.DB
	{
		dirs := _conf.DirConf

		dc := &DatabaseConf{
			Type:     "sqlite",
			Database: filepath.Join(dirs.CacheDir, "scraper-store.db"),
			Log: SQLLogConf{
				IgnoreNotFound: true,
			},
		}
		_, sdb, err = newDB(dc)
		if err != nil {
			return errors.Wrap(err, "new scraper store")
		}
	}

	sc := &scrape.Context{
		Seed:  u,
		Fatal: cc.Bool("fatal"),
		DB:    sdb,
	}

	ctx := context.Background()
	ctx = torrenti.IndexerContextKey.WithValue(ctx, getTorrentIndexer())
	ctx = subi.IndexerContextKey.WithValue(ctx, getSubIndexer())
	ctx = util.DirConfContextKey.WithValue(ctx, &_conf.DirConf)
	ctx = scrape.ContextKey.WithValue(ctx, sc)

	c := colly.NewCollector(
		colly.CacheDir(cache),
		colly.Debugger(&debug.WebDebugger{
			Address: _conf.Scrape.Debug.Addr,
		}),
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.74 Safari/537.36"),
	)

	sc.Context = ctx
	sc.C = c

	if err = sc.Init(); err != nil {
		return errors.Wrap(err, "init scraper")
	}

	if cc.Bool("seed") {
		return sc.C.Visit(first)
	}

	if !cc.Bool("queue") {
		if err = sc.Queue.AddURL(first); err != nil {
			return errors.Wrap(err, "add seed url to queue")
		}
	}

	return sc.Queue.Run(c)
}
