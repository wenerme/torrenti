package btbtt

import (
	"context"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/wenerme/torrenti/pkg/scrape"
	"github.com/wenerme/torrenti/pkg/scrape/handlers"
	"github.com/wenerme/torrenti/pkg/scrape/handlers/archives"

	"github.com/gocolly/colly/v2"

	"github.com/wenerme/torrenti/pkg/subi"

	"github.com/rs/zerolog"

	"golang.org/x/exp/slices"

	"github.com/pkg/errors"

	"github.com/rs/zerolog/log"
	"github.com/wenerme/torrenti/pkg/torrenti"
	"github.com/wenerme/torrenti/pkg/torrenti/util"
)

type VisitOptions struct {
	URL     string
	Source  string
	Reason  string
	Request *colly.Request
}

const KeySkipMarkVisit = "SkipMarkVisit"

var (
	threadIndex = regexp.MustCompile(`^/thread-index-fid-(\d+)-tid-(\d+)(-page-(\d+))?.htm$`)
	indexPage   = regexp.MustCompile(`^/index-index-page-(\d+).htm$`)
	forumIndex  = regexp.MustCompile(`^/forum-index-fid-(\d+)(-page-(\d+))?.htm$`)
)

var fids = []string{
	"950",  // 剧集
	"951",  // 电影
	"1183", // 高清电影
	"953",  // 音乐
	"981",  // 动漫
	"1107", // 综艺
	"1151", // 图书
	"957",  // 美图
	"1191", // 音轨字幕
}

func init() {
	slices.Sort(fids)

	Name := "btbtt"
	scrape.RegisterScraper(&scrape.Scraper{
		Name: Name,
		Support: func(ctx *scrape.Context) bool {
			return ctx.Seed.Hostname() == "www.btbtt12.com"
		},
		InitCollector: func(ctx *scrape.Context, c *colly.Collector) error {
			c.SetRequestTimeout(30 * time.Second)
			// 20MB
			c.MaxBodySize = 20 * 1024 * 1024
			c.AllowedDomains = []string{"www.btbtt12.com"}
			return nil
		},
		SetupCollector: func(sc *scrape.Context, c *colly.Collector) error {
			ctx := sc.Context
			st := sc.Stat

			lastReport := time.Now()
			report := func(log zerolog.Logger) {
				lastReport = time.Now()
				log.Info().
					Object("stat", st).
					Msg("report")
			}

			vis := func(o scrape.QueueVisitOptions) {
				var err error

				err = sc.QueueVisit(o)
				sc.OnError(&scrape.OnErrorEvent{
					Error:   err,
					Request: o.Request,
					Message: "queue visit",
				})
			}

			c.OnScraped(func(resp *colly.Response) {
				u := resp.Request.URL.String()
				log := log.With().Str("url", u).Logger()

				now := time.Now()

				defer func() {
					if (st.ScrapedCount%1000 == 0 && now.Sub(lastReport) > (time.Second*10)) || now.Sub(lastReport) > (time.Second*30) {
						lastReport = now
						report(log)
					}
				}()
			})

			c.OnHTML("a[href]", func(e *colly.HTMLElement) {
				src := e.Request.URL.Path
				target, _ := url.Parse(e.Request.AbsoluteURL(e.Attr("href")))
				var to string
				if target != nil {
					to = target.Path
				}

				visit := func(url string, res string) {
					vis(scrape.QueueVisitOptions{
						URL:     url,
						Source:  src,
						Reason:  res,
						Referer: e.Request,
					})
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
				case forumIndex.MatchString(to):
					fid := forumIndex.FindStringSubmatch(to)[0]
					if util.BinarySearchContain(fids, fid) {
						visit(to, "forum index")
					}
				case strings.HasPrefix(src, "/") && indexPage.MatchString(to):
					visit(to, "home pagination")
				case strings.HasPrefix(src, "/index-index-page") && indexPage.MatchString(to):
					visit(to, "index pagination")
				case strings.HasPrefix(src, "/") && threadIndex.MatchString(to):
					visit(to, "home to thread")
				case strings.HasPrefix(src, "/index-index-page") && threadIndex.MatchString(to):
					visit(to, "list to thread")
				case strings.HasPrefix(to, "/attach-download-fid"):
					visit(to, "download file")
				case strings.HasPrefix(to, "/attach-dialog-fid"):
					visit(strings.ReplaceAll(to, "/attach-dialog-fid", "/attach-download-fid"), "download page")
				case strings.HasPrefix(src, "/thread-index-fid") && threadIndex.MatchString(to):
					visit(to, "xref")
				default:
					log.Trace().Str("href", to).Str("from", src).Msg("href")
				}
			})

			c.OnResponse(func(resp *colly.Response) {
				if hdr := resp.Headers.Get("Content-Disposition"); hdr != "" {
					// resp.Ctx.Put(KeySkipMarkVisit, true)
					st.FileCount++

					u := resp.Request.URL.String()
					_, params, _ := mime.ParseMediaType(hdr)
					filename := params["filename"]

					log := log.With().Int("depth", resp.Request.Depth).Str("url", u).Str("file", filename).Logger()

					if err := sc.Store.MarkFile(u); err != nil {
						log.Err(err).Msg("mark file")
					}

					log.Debug().
						Int("file_count", st.FileCount).
						Msg("detect")
					fi := &util.File{
						Path:   filename,
						Length: int64(len(resp.Body)),
						Data:   resp.Body,
						URL:    u,
					}
					err := handle(ctx, st, fi)
					if err != nil {
						sc.OnError(&scrape.OnErrorEvent{
							Response: resp,
							Error:    err,
							Message:  "handle",
							Log: func(log zerolog.Context) zerolog.Context {
								return log.
									Str("file", filename).
									Str("ext", filepath.Ext(filename)).
									Str("mime", http.DetectContentType(fi.Data))
							},
						})
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

func handle(ctx context.Context, st *scrape.Stat, f *util.File) (err error) {
	if f.IsDir() {
		return nil
	}
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		err, _ = r.(error)
		if err == nil {
			err = errors.Errorf("handle panic: %v", r)
		}
	}()
	cb := func(ctx context.Context, file *util.File) error {
		return handle(ctx, st, file)
	}
	ext := handlers.FileExt(f)
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
