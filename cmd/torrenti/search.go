package main

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/blugelabs/bluge/search/highlight"
	"github.com/dustin/go-humanize"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"github.com/urfave/cli/v2"
	torrentiv1 "github.com/wenerme/torrenti/pkg/apis/media/torrenti/v1"
	webv1 "github.com/wenerme/torrenti/pkg/apis/media/web/v1"
	"github.com/wenerme/torrenti/pkg/search"
	"github.com/wenerme/torrenti/pkg/serve"
	"github.com/wenerme/torrenti/pkg/subi"
	"github.com/wenerme/torrenti/pkg/torrenti"
	"github.com/wenerme/torrenti/pkg/torrenti/models"
	"github.com/wenerme/torrenti/pkg/torrenti/services"
	"github.com/wenerme/torrenti/pkg/web"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

func runSearchStat(cc *cli.Context) (err error) {
	return
}

func runSearchQuery(cc *cli.Context) (err error) {
	if cc.NArg() == 0 {
		return
	}
	return fxApp(cc, fx.Invoke(func(ws webv1.WebServiceServer, ss *search.Service, ts *torrenti.Service) (err error) {
		for _, v := range cc.Args().Slice() {
			log.Info().Msgf("searching %s", v)
			sr, err := ss.SearchTorrent(context.Background(), &search.SearchRequest{
				QueryString: v,
				Limit:       cc.Int("limit"),
				Offset:      cc.Int("offset"),
			})
			if err != nil {
				return err
			}
			fmt.Printf("Search Result: %v in %v with max score %.2f\n", sr.Count, sr.Duration, sr.MaxScore)

			var out []*models.MetaFile
			err = ts.DB.Model(models.MetaFile{}).
				Where("torrent_hash in (?)", lo.Map(sr.Docs, func(t *search.DocumentMatch, i int) string {
					return t.ID
				})).
				Select([]string{"id", "filename", "content_hash", "torrent_hash", "creation_date"}).
				Preload("Torrent", func(db *gorm.DB) *gorm.DB {
					return db.Select([]string{"name", "hash", "total_file_size", "file_count"})
				}).
				Find(&out).Error
			if err != nil {
				return err
			}
			byID := lo.KeyBy(out, func(v *models.MetaFile) string {
				return v.TorrentHash
			})
			hi := highlight.NewANSIHighlighterColor("\x1b[30;44m")
			// hi := highlight.NewHTMLHighlighter()
			for _, match := range sr.Docs {
				mod := byID[match.ID]
				if mod == nil {
					log.Warn().Str("id", match.ID).Msg("model not found")
					continue
				}
				fmt.Print("> ")
				fragment := hi.BestFragment(match.Locations[search.TorrentFieldMetaFileName], []byte(mod.Filename))
				println(fragment)
				println("  ", hi.BestFragment(match.Locations[search.TorrentFieldTorrentFileName], []byte(mod.Torrent.Name)))
				println("  ",
					time.Unix(mod.CreationDate, 0).Format("2006-01-02T15:04"),
					humanize.Bytes(uint64(mod.Torrent.TotalFileSize)),
					" - ",
					fmt.Sprintf("%.2f", match.Score),
				)
			}
		}
		return
	}))
}

func runSearchIndex(cc *cli.Context) (err error) {
	err = fxApp(cc, fx.Invoke(searchIndex))
	if err != nil {
		return err
	}
	return
}

func searchIndex(ss *search.Service, ts *torrenti.Service) (err error) {
	db := ts.DB
	var out []*models.MetaFile
	lastID := uint(0)
	n := 0
	start := time.Now()
	for {
		out = nil
		err = db.
			Model(models.MetaFile{}).Order("id").Where("id > ?", lastID).
			Select([]string{"id", "filename", "content_hash", "torrent_hash", "creation_date"}).
			Preload("Torrent", func(db *gorm.DB) *gorm.DB {
				return db.Select([]string{"name", "hash", "total_file_size", "file_count"})
			}).
			Limit(1000).Find(&out).Error
		if err != nil {
			break
		}
		if len(out) == 0 {
			break
		}
		lastID = out[len(out)-1].ID
		n += len(out)

		docs := make([]*search.TorrentDocument, 0, len(out))
		for _, v := range out {
			if v.Torrent == nil {
				log.Warn().Str("hash", v.TorrentHash).Msg("torrent not found")
				continue
			}
			docs = append(docs, &search.TorrentDocument{
				ID:              v.TorrentHash,
				MetaFileName:    v.Filename,
				TorrentFileName: v.Torrent.Name,
				Size:            v.Torrent.TotalFileSize,
				CreatedAt:       time.Unix(v.CreationDate, 0),
			})
		}
		err = ss.IndexTorrent(context.Background(), docs)
		if err != nil {
			return err
		}
	}
	log.Info().Int("count", n).Dur("duration", time.Now().Sub(start)).Msg("indexed")
	return
}

func fxApp(cc *cli.Context, opts ...fx.Option) (err error) {
	sc := &serve.Context{
		Cli:     cc,
		Context: context.Background(),
	}

	opts = append(opts,
		fx.NopLogger,
		fx.Supply(sc, cc, _conf, &_conf.GRPC, &_conf.HTTP, &_conf.Debug),
		fx.Module("torrent",
			fx.Provide(
				func(conf *Config) (svc *torrenti.Service, err error) {
					_, gdb, err := newDB(&conf.Torrent.DB)
					if err != nil {
						return
					}
					svc, err = torrenti.NewIndexer(torrenti.NewServiceOptions{DB: gdb})
					return
				},
				func(svc *torrenti.Service) (svr torrentiv1.TorrentIndexServiceServer, err error) {
					svr = &services.TorrentIndexerServer{Indexer: svc}
					serve.RegisterEndpoints(&serve.ServiceEndpoint{
						Desc:            &torrentiv1.TorrentIndexService_ServiceDesc,
						Impl:            &services.TorrentIndexerServer{Indexer: svc},
						RegisterGateway: torrentiv1.RegisterTorrentIndexServiceHandler,
					})
					return
				},
			),
		),
		fx.Module("sub",
			fx.Provide(
				func(conf *Config) (svc *subi.Indexer, err error) {
					_, gdb, err := newDB(&conf.Sub.DB)
					if err != nil {
						return
					}
					return subi.NewIndexer(subi.NewIndexerOptions{DB: gdb})
				},
			),
		),
		fx.Module("search",
			fx.Provide(func(conf *Config) (svc *search.Service, err error) {
				svc, err = search.NewService(search.NewServiceOptions{
					DataDir: filepath.Join(conf.DataDir, "search"),
				})
				if err == nil {
				}
				return
			}),
		),
		fx.Module("web",
			fx.Provide(func(ss *search.Service, ti *torrenti.Service) (svr webv1.WebServiceServer) {
				svr = web.NewWebServiceServer(web.NewWebServiceServerOptions{
					DB:     ti.DB,
					Search: ss,
				})
				serve.RegisterEndpoints(&serve.ServiceEndpoint{
					Desc:            &webv1.WebService_ServiceDesc,
					Impl:            svr,
					RegisterGateway: webv1.RegisterWebServiceHandler,
				})
				return
			}),
		),
		fx.Module("http", fx.Invoke(serveHTTP)),
		fx.Module("debug", fx.Invoke(serveDebug)),
		fx.Module("scrape", fx.Invoke(serveScrape)),
		fx.Module("grpc", fx.Provide(setupGRPC), fx.Invoke(serveGRPCGateway)),
	)
	app := fx.New(opts...)
	return app.Start(sc.Context)
}
