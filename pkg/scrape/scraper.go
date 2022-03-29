package scrape

import (
	"context"
	"net/url"
	"time"

	"github.com/gocolly/colly/v2/queue"
	"github.com/rs/zerolog/log"
	"go.uber.org/multierr"
	"gorm.io/gorm"

	"github.com/gocolly/colly/v2"

	"github.com/wenerme/torrenti/pkg/torrenti/util"
)

type Context struct {
	Seed    *url.URL
	Fatal   bool
	Store   *VisitStore
	Queue   *queue.Queue
	DB      *gorm.DB
	Context context.Context
	Stat    *Stat
	C       *colly.Collector
}

var ContextKey = util.ContextKey[*Context]{Name: "scrape.Context"}

type Scraper struct {
	Name           string
	Support        func(ctx *Context) bool
	InitCollector  func(ctx *Context, c *colly.Collector) error
	SetupCollector func(ctx *Context, c *colly.Collector) error
}

var scrapers []*Scraper

func (sc *Context) Init() (err error) {
	if sc.Store == nil {
		sc.Store = &VisitStore{DB: sc.DB}
	}
	if sc.Queue == nil {
		sc.Queue, err = queue.New(2, &QueueStorage{DB: sc.DB})
	}
	if sc.Stat == nil {
		sc.Stat = &Stat{}
	}

	err = multierr.Combine(err, sc.Store.Init())
	if err != nil {
		return
	}

	c := sc.C
	for _, s := range scrapers {
		if s.Support != nil && !s.Support(sc) {
			continue
		}
		if err = s.InitCollector(sc, c); err != nil {
			return
		}
		if err = s.SetupCollector(sc, c); err != nil {
			return
		}
	}

	c.OnRequest(func(r *colly.Request) {
		sc.Stat.RequestCount++

		// always request seed page
		if r.Depth < 2 {
			r.Headers.Set("Cache-Control", "no-cache")
		}

		u := r.URL.String()
		if err := sc.Store.MarkVisiting(u); err != nil {
			log.Err(err).Str("url", u).Msg("mark visiting")
		}
		log.Debug().Str("url", u).Msg("visiting")
	})

	c.OnError(func(r *colly.Response, err error) {
		sc.Stat.ErrorCount++
		log.Err(err).Str("method", r.Request.Method).Str("url", r.Request.URL.String()).Msg("error")
	})

	c.OnScraped(sc.onScraped)
	return
}

type QueueVisitOptions struct {
	URL     string
	Source  string
	Reason  string
	Request *colly.Request
	Origin  *colly.Request
}

func (sc *Context) QueueVisit(o QueueVisitOptions) (err error) {
	if o.Origin == nil {
		o.Origin = &colly.Request{
			Method: "GET",
			URL:    sc.Seed,
		}
	}

	r := o.Request
	u := o.URL
	u = o.Origin.AbsoluteURL(u)
	o.Source = o.Origin.AbsoluteURL(o.Source)

	if r == nil {
		r = &colly.Request{
			Method: "GET",
			Depth:  o.Origin.Depth + 1,
		}
		r.URL, err = url.Parse(u)
		if err != nil {
			return
		}
	}

	log := log.With().Str("href", u).Str("src", o.Source).Str("reason", o.Reason).Logger()

	visited, err := sc.Store.IsScraped(u)
	if err != nil {
		log.Err(err).Msg("query visit")
	}
	if visited {
		sc.Stat.SkipVisitCount++
		log.Debug().Msg("duplicate visit")
		return
	}

	log.Debug().Msg("queue")
	err = sc.Queue.AddRequest(r)
	return
}

func (sc *Context) onScraped(r *colly.Response) {
	stat := sc.Stat
	stat.ScrapedCount++
	stat.LastScrapedAt = time.Now()

	u := r.Request.URL.String()
	log := log.With().Str("url", u).Logger()

	if v, _ := r.Ctx.GetAny(KeySkipMarkVisit).(bool); v == true {
		stat.SkipMarkVisitCount++
		log.Trace().Msg("skip mark visit")
		return
	}

	if err := sc.Store.MarkScraped(u); err != nil {
		log.Err(err).Msg("mark visited")
	} else {
		log.Trace().Msg("mark visited")
	}
}

const KeySkipMarkVisit = "SkipMarkVisit"

func RegisterScraper(v *Scraper) {
	scrapers = append(scrapers, v)
}
