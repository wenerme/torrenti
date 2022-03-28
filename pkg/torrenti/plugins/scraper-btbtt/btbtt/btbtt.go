package btbtt

import (
	"context"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"

	"github.com/wenerme/torrenti/pkg/subi"

	"github.com/rs/zerolog"

	"gorm.io/gorm"

	"golang.org/x/exp/slices"

	"github.com/wenerme/torrenti/pkg/torrenti/scraper/handlers"
	"github.com/wenerme/torrenti/pkg/torrenti/scraper/handlers/archives"

	"github.com/pkg/errors"

	"github.com/rs/zerolog/log"
	"github.com/wenerme/torrenti/pkg/torrenti"
	"github.com/wenerme/torrenti/pkg/torrenti/scraper"
	"github.com/wenerme/torrenti/pkg/torrenti/util"
)

type VisitOptions struct {
	URL    string
	Source string
	Reason string
}

type ScrapeContext struct {
	Store *gorm.DB
}

const KeySkipMarkVisit = "SkipMarkVisit"

func init() {
	Name := "btbtt"
	scraper.RegisterScraper(&scraper.Scraper{
		Name: Name,
		Support: func(ctx context.Context) bool {
			return scraper.OptionContextKey.Must(ctx).Seed.Hostname() == "www.btbtt12.com"
		},
		InitCollector: func(ctx context.Context, c *colly.Collector) error {
			c.SetRequestTimeout(30 * time.Second)
			// 20MB
			c.MaxBodySize = 20 * 1024 * 1024
			c.AllowedDomains = []string{"www.btbtt12.com"}
			return nil
		},
		SetupCollector: func(ctx context.Context, c *colly.Collector) error {
			so := scraper.OptionContextKey.Must(ctx)
			st := &scraper.Stat{}

			lastReport := time.Now()
			report := func(log zerolog.Logger) {
				lastReport = time.Now()
				log.Info().
					Object("stat", st).
					Msg("report")
			}
			fatal := func(log zerolog.Logger, msg string) {
				st.ErrorCount++
				if so.Fatal {
					report(log)
					log.Fatal().Msg(msg)
				} else {
					log.Error().Msg(msg)
				}
			}

			store := so.Store
			vis := func(o VisitOptions) {
				log := log.With().Str("href", o.URL).Str("src", o.Source).Str("reason", o.Reason).Logger()
				if o.URL == "" {
					log.Warn().Msg("empty url")
					return
				}

				visited, err := store.IsScraped(o.URL)
				if err != nil {
					log.Err(err).Msg("query visit")
				}
				if visited {
					st.SkipVisitCount++
					log.Debug().Msg("duplicate visit")
					return
				}

				log.Debug().Msg("visit")

				if err = c.Visit(o.URL); err != nil {
					if colly.ErrAlreadyVisited == err {
						st.AlreadyVisitCount++
						return
					}

					fatal(log.With().Err(err).Logger(), "visit")
				}
			}

			c.OnScraped(func(resp *colly.Response) {
				st.ScrapedCount++
				u := resp.Request.URL.String()
				now := time.Now()
				log := log.With().Str("url", u).Logger()

				defer func() {
					if (st.ScrapedCount%1000 == 0 && now.Sub(lastReport) > (time.Second*10)) || now.Sub(lastReport) > (time.Second*30) {
						lastReport = now
						report(log)
					}
				}()
				st.LastScrapedAt = now

				if v, _ := resp.Ctx.GetAny(KeySkipMarkVisit).(bool); v == true {
					st.SkipMarkVisitCount++
					log.Trace().Msg("skip mark visit")
					return
				}

				if err := store.MarkScraped(u); err != nil {
					log.Err(err).Msg("mark visit")
				} else {
					log.Trace().Msg("mark visit")
				}
			})

			c.OnHTML("a[href]", func(e *colly.HTMLElement) {
				src := e.Request.URL.Path
				target, _ := url.Parse(e.Request.AbsoluteURL(e.Attr("href")))
				var to string
				if target != nil {
					to = target.Path
				}

				visit := func(url string, res string) {
					url = e.Request.AbsoluteURL(url)
					vis(VisitOptions{URL: url, Source: src, Reason: "link"})
				}
				switch {
				case target == nil:
				case strings.HasPrefix(to, "post-"): // POST 请求页面
				case strings.HasPrefix(to, "javascript:"):
				default:
					goto VALID
				}
				return
			VALID:

				switch {
				case strings.HasPrefix(src, "/") && strings.HasPrefix(to, "/index-index-page"):
					visit(to, "home pagination")
				case strings.HasPrefix(src, "/") && strings.HasPrefix(to, "/thread-index-fid"):
					visit(to, "home to thread")
				case strings.HasPrefix(src, "/index-index-page") && strings.HasPrefix(to, "/thread-index-fid"):
					visit(to, "list to thread")
				case strings.HasPrefix(src, "/attach-dialog-fid") && strings.HasPrefix(to, "/attach-download-fid"):
					visit(to, "download file")
				case strings.HasPrefix(src, "/thread-index-fid") && strings.HasPrefix(to, "/attach-dialog-fid"):
					visit(to, "download page")
				case strings.HasPrefix(src, "/thread-index-fid") && strings.HasPrefix(to, "/attach-dialog-fid"):
					visit(to, "xref")
				default:
					log.Trace().Str("href", to).Str("from", src).Msg("href")
				}
			})

			c.OnRequest(func(r *colly.Request) {
				st.RequestCount++
				u := r.URL.String()
				if err := store.MarkVisiting(u); err != nil {
					log.Err(err).Str("url", u).Msg("mark visiting")
				}
				log.Debug().Str("url", u).Msg("visiting")
			})

			c.OnResponse(func(resp *colly.Response) {
				if hdr := resp.Headers.Get("Content-Disposition"); hdr != "" {
					resp.Ctx.Put(KeySkipMarkVisit, true)
					st.FileCount++

					_, params, _ := mime.ParseMediaType(hdr)

					filename := params["filename"]

					log := log.With().Str("url", resp.Request.URL.String()).Str("file", filename).Logger()
					log.Debug().
						Int("file_count", st.FileCount).
						Msg("detect")
					fi := &util.File{
						Path:   filename,
						Length: int64(len(resp.Body)),
						Data:   resp.Body,
						URL:    resp.Request.URL.String(),
					}
					err := handle(ctx, st, fi)
					if err != nil {
						fatal(
							log.With().
								Err(err).
								Str("ext", filepath.Ext(filename)).
								Str("mime", http.DetectContentType(fi.Data)).
								Logger(),
							"handle",
						)
					}
				}
			})

			return nil
		},
	})
}

func handleSubtitle(ctx context.Context, f *util.File) (err error) {
	idx := subi.IndexerContextKey.Get(ctx)
	if idx == nil {
		log.Trace().Msg("skip subtitle")
		return
	}
	err = idx.Index(f)
	return
}

func handleTorrent(ctx context.Context, f *util.File) (err error) {
	idx := torrenti.IndexerContextKey.Must(ctx)
	t := &torrenti.Torrent{
		URL: f.URL,
	}
	t.FileInfo = f
	t.Data = f.Data
	_, err = idx.IndexTorrent(ctx, t)
	return
}

var (
	ignoredExts = []string{}
	triExts     = []string{".txt", ".tv", ".url", ".ds_store", ".db", ".sqlite", ".ini"}
	officeExts  = []string{".docx", ".doc"}
	imagesExts  = []string{".jpg", ".jpeg", ".png"}
)

func init() {
	ignoredExts = append(ignoredExts, triExts...)
	ignoredExts = append(ignoredExts, officeExts...)
	ignoredExts = append(ignoredExts, imagesExts...)
	slices.Sort(ignoredExts)
}

func handle(ctx context.Context, st *scraper.Stat, f *util.File) (err error) {
	if f.IsDir() {
		return nil
	}
	cb := func(ctx context.Context, file *util.File) error {
		return handle(ctx, st, file)
	}
	ext := handlers.Ext(f)
	fn := f.Name()

	log := log.With().Str("file", f.Path).Str("ext", ext).Logger()

	st.CountExt(ext)
	log.Trace().Msg("handle")

	switch {
	case strings.HasPrefix(fn, "."):
		log.Trace().Msg("skip hidden")
		return
	case util.BinarySearchContain(ignoredExts, ext):
		log.Trace().Msg("skip uninterested")
		return
	}
	switch ext {
	case ".zip":
		err = archives.Unzip(ctx, f, cb)
	case ".rar":
		err = archives.Unrar(ctx, f, cb)
	case ".7z":
		err = archives.Un7z(ctx, f, cb)
	case ".torrent":
		err = handleTorrent(ctx, f)
	default:
		switch {
		case handlers.IsSubtitleExt(ext):
			err = handleSubtitle(ctx, f)
		default:
			err = errors.Errorf("unable to handle file: %q", f.Path)
		}
	}
	if err != nil {
		log.Error().
			Str("mime", http.DetectContentType(f.Data)).
			Str("dump", dump(f)).
			Msg("unable to handle")
	}
	return
}

func dump(f *util.File) string {
	fn := filepath.Join(os.TempDir(), f.Name())
	_ = os.WriteFile(fn, f.Data, 0o644)
	return fn
}
