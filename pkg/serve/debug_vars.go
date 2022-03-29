package serve

import (
	_ "expvar"
	"net/http"
)

func RegisterDebugEndpoints() {
	RegisterEndpoints(&HTTPEndpoint{
		Method:  http.MethodGet,
		Path:    "/debug/vars",
		Handler: http.DefaultServeMux,
		EndpointDesc: EndpointDesc{
			Selector: "debug",
		},
	})
}
