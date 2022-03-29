package scrape

import (
	"time"

	"github.com/rs/zerolog"
)

type Stat struct {
	FileCount          int
	ErrorCount         int
	RequestCount       int
	SkipVisitCount     int
	AlreadyVisitCount  int
	SkipMarkVisitCount int
	ScrapedCount       int
	LastScrapedAt      time.Time
	Extensions         map[string]int
	ExtensionCount     int
}

func (s *Stat) CountExt(ext string) {
	s.ExtensionCount++
	if s.Extensions == nil {
		s.Extensions = make(map[string]int)
	}
	s.Extensions[ext]++
}

func (s *Stat) MarshalZerologObject(e *zerolog.Event) {
	e.
		Int("file", s.FileCount).
		Int("request", s.RequestCount).
		Int("scraped", s.ScrapedCount).
		Int("skip", s.SkipVisitCount).
		Int("dup", s.AlreadyVisitCount).
		Int("skip_mark", s.SkipMarkVisitCount).
		Int("ext", s.ExtensionCount).
		Int("err", s.ErrorCount)
	for k, v := range s.Extensions {
		e.Int(k, v)
	}
}
