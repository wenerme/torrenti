package scrape

import (
	"context"
	"net/url"

	"github.com/gocolly/colly/v2"
	scraperv1 "github.com/wenerme/torrenti/pkg/apis/media/scraper/v1"
)

type svc struct {
	scraperv1.UnimplementedScrapeServiceServer
	s *Service
}

func (s *svc) Scrape(ctx context.Context, req *scraperv1.ScrapeRequest) (resp *scraperv1.ScrapeResponse, err error) {
	r, err := s.s.Scrape(ctx, &ScrapeRequest{
		URL: req.GetUrl(),
	})
	if r != nil {
		resp = &scraperv1.ScrapeResponse{}
	}
	return
}

func (s *svc) State(ctx context.Context, req *scraperv1.StateRequest) (*scraperv1.StateResponse, error) {
	// TODO implement me
	panic("implement me")
}

type Service struct {
	QueueStorage *QueueStorage
}

func (s *Service) State(ctx context.Context, req *ScrapeRequest) (resp *ScrapeResponse, err error) {
	return
}

func (s *Service) Scrape(ctx context.Context, req *ScrapeRequest) (resp *ScrapeResponse, err error) {
	r := &colly.Request{
		Method: "GET",
		Ctx:    colly.NewContext(),
	}
	r.URL, err = url.Parse(req.URL)
	if err != nil {
		return
	}
	if req.Referer != "" {
		r.Ctx.Put(ctxKeyReferer, req.Referer)
	}

	_, err = s.QueueStorage.StoreRequest(r)
	return
}

type ScrapeRequest struct {
	URL     string
	Referer string
}
type ScrapeResponse struct{}

type StateRequest struct {
	URL string
}
type StateResponse struct{}
