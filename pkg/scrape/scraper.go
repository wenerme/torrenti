package scrape

import (
	"context"
	"net/url"
	"runtime"
	"time"

	"github.com/rs/zerolog"

	"github.com/gocolly/colly/v2/queue"
	"github.com/rs/zerolog/log"
	"go.uber.org/multierr"
	"gorm.io/gorm"

	"github.com/gocolly/colly/v2"

	"github.com/wenerme/torrenti/pkg/torrenti/util"
)

type NewContextOptions struct {
	DB              *gorm.DB
	Concurrent      int
	Fatal           bool
	DirectSeedDepth int
	Seed            string
}

type Context struct {
	Seed            *url.URL
	DirectSeedDepth int
	Fatal           bool
	Store           *VisitStore
	Queue           *queue.Queue
	QueueStorage    queue.Storage
	DB              *gorm.DB
	Context         context.Context
	Stat            *Stat
	Collector       *colly.Collector
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
	if sc.Context == nil {
		sc.Context = context.Background()
	}
	if !ContextKey.Exists(sc.Context) {
		sc.Context = ContextKey.WithValue(sc.Context, sc)
	}

	if sc.Store == nil {
		sc.Store = &VisitStore{DB: sc.DB}
	}
	if sc.QueueStorage == nil {
		sc.QueueStorage = &QueueStorage{DB: sc.DB}
	}
	if sc.Queue == nil {
		sc.Queue, err = queue.New(runtime.GOMAXPROCS(1), sc.QueueStorage)
	}
	if sc.Stat == nil {
		sc.Stat = &Stat{}
	}
	if sc.Collector == nil {
		sc.Collector = colly.NewCollector()
	}

	err = multierr.Combine(err, sc.Store.Init())
	if err != nil {
		return
	}

	c := sc.Collector
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
		if r.Depth <= sc.DirectSeedDepth {
			r.Headers.Set("Cache-Control", "no-cache")
		}

		u := r.URL.String()
		if err := sc.Store.MarkVisiting(u); err != nil {
			log.Err(err).Str("url", u).Msg("mark visiting")
		}
		log.Debug().Str("url", u).Msg("visiting")
	})

	c.OnError(func(r *colly.Response, err error) {
		sc.OnError(&OnErrorEvent{
			Response: r,
			Error:    err,
		})
	})

	c.OnScraped(sc.onScraped)
	return
}

type QueueVisitOptions struct {
	URL     string
	Source  string
	Reason  string
	Request *colly.Request
	Referer *colly.Request
}

func (sc *Context) QueueVisit(o QueueVisitOptions) (err error) {
	referer := ""
	if o.Referer != nil {
		referer = o.Referer.URL.String()
	}
	if o.Referer == nil {
		o.Referer = &colly.Request{
			Method: "GET",
			URL:    sc.Seed,
		}
	}

	r := o.Request
	u := o.URL
	u = o.Referer.AbsoluteURL(u)
	o.Source = o.Referer.AbsoluteURL(o.Source)

	if r == nil {
		r, err = o.Referer.New("GET", u, nil)
		r.Depth = o.Referer.Depth + 1
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

	if referer != "" {
		if r.Ctx == nil {
			r.Ctx = colly.NewContext()
		}
		r.Ctx.Put(ctxKeyReferer, referer)
	}
	log.Debug().Msg("queue")
	err = sc.Queue.AddRequest(r)
	return
}

type OnErrorEvent struct {
	URL      string
	Request  *colly.Request
	Response *colly.Response
	Error    error
	Log      func(zerolog.Context) zerolog.Context
	Message  string
}

func (sc *Context) OnError(e *OnErrorEvent) {
	if e.Error == nil {
		return
	}

	sc.Stat.ErrorCount++
	if e.Request == nil && e.Response != nil {
		e.Request = e.Response.Request
	}

	ec := log.With()
	if e.Request != nil {
		if e.URL == "" {
			e.URL = e.Request.URL.String()
		}
		ec = ec.Int("depth", e.Request.Depth)
	}

	if e.URL != "" {
		ec = ec.Str("url", e.URL)
	}
	el := ec.Logger()
	if e.Log != nil {
		el.UpdateContext(e.Log)
	}
	if sc.Fatal {
		el.Fatal().Err(e.Error).Msg(e.Message)
	} else {
		el.Error().Err(e.Error).Msg(e.Message)
	}
	if err := sc.Store.MarkError(e.URL, e.Error); err != nil {
		el.Err(err).Msg("mark error")
	}
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
