package serve

import "net/http"

func RegisterReadyEndpoints(f func() error) {
	RegisterEndpoints(&HTTPEndpoint{
		Method: http.MethodGet,
		Path:   "/ready",
		HandlerFunc: func(w http.ResponseWriter, r *http.Request) {
			if err := f(); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
		},
		EndpointDesc: EndpointDesc{
			Selector: "debug",
		},
	})
}

func RegisterHealthEndpoints(f func() error) {
	RegisterEndpoints(&HTTPEndpoint{
		Method: http.MethodGet,
		Path:   "/health",
		HandlerFunc: func(w http.ResponseWriter, r *http.Request) {
			if err := f(); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
		},
		EndpointDesc: EndpointDesc{
			Selector: "debug",
		},
	})
}
