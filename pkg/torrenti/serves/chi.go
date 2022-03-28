package serves

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func ChiRoute(r chi.Router, e *HTTPEndpoint) (err error) {
	if e == nil || e.Disabled {
		return nil
	}
	if err = e.Validate(); err != nil {
		return errors.Wrap(err, "invalid HTTP endpoint")
	}
	if len(e.Middlewares) > 0 {
		r = r.With(e.Middlewares...)
	}
	if e.Handler != nil {
		if len(e.Methods) > 0 {
			for _, v := range e.Methods {
				log.Trace().Str("path", e.Path).Str("method", v).Msg("handle")
				r.Method(v, e.Path, e.Handler)
			}
		} else {
			log.Trace().Str("path", e.Path).Str("method", "ANY").Msg("handle")
			r.Handle(e.Path, e.Handler)
		}
	}
	for _, v := range e.Children {
		if err = ChiRoute(r, v); err != nil {
			return err
		}
	}
	return
}

func LogRouter(log zerolog.Logger, r chi.Router) {
	path := map[string][]string{}

	err := chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		path[route] = append(path[route], method)
		return nil
	})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to walk router")
	}

	for p, m := range path {
		if len(m) < 9 {
			for _, v := range m {
				log.Info().Str("method", v).Str("path", p).Msg("route")
			}
		} else {
			log.Info().Str("method", "ANY").Str("path", p).Msg("route")
		}
	}
}
