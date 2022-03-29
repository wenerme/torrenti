package serve

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/rs/zerolog/log"
	"golang.org/x/exp/slices"
)

var Endpoints []Endpoint

func RegisterEndpoints(eps ...Endpoint) {
	for _, v := range eps {
		log.Trace().
			Str("endpoint", v.String()).
			Str("type", reflect.TypeOf(v).String()).
			Msg("registering endpoint")
	}
	Endpoints = append(Endpoints, eps...)
}

type EndpointDesc struct {
	Name                   string
	Description            string
	Selector               string
	Tags                   []string
	Disabled               bool
	DeprecationDescription string
}

func (e EndpointDesc) String() string {
	var s []string
	if e.Name != "" {
		s = append(s, e.Name)
	}
	if e.Selector != "" {
		s = append(s, "@"+e.Selector)
	}
	if e.Disabled {
		s = append(s, "disabled")
	}
	if len(e.Tags) != 0 {
		s = append(s, "tags="+strings.Join(e.Tags, ","))
	}
	return fmt.Sprintf("Endpoint(%v)", strings.Join(s, ","))
}

type Endpoint interface {
	GetEndpointDesc() *EndpointDesc
	String() string
}

type SelectEndpointOptions[T Endpoint] struct {
	Selector   string
	Filter     func(e T) bool
	Comparator func(a, b T) bool
}

func HTTPEndpointSortByPathLen(a, b *HTTPEndpoint) bool {
	return len(a.Path) > len(b.Path)
}

func SelectEndpoints[T Endpoint](o SelectEndpointOptions[T], f func(e T) error) error {
	var eps []T
	for _, v := range Endpoints {
		desc := v.GetEndpointDesc()
		switch {
		case desc.Selector != o.Selector:
		case desc.Disabled:
		default:
			vv, ok := v.(T)
			if ok && ((o.Filter != nil && o.Filter(vv)) || o.Filter == nil) {
				eps = append(eps, vv)
			}
		}
	}

	if o.Comparator != nil {
		slices.SortFunc(eps, o.Comparator)
	}

	log.Trace().
		Int("count", len(eps)).
		Str("select", fmt.Sprint(o)).
		Str("type", reflect.TypeOf(new(T)).Elem().String()).
		Msg("select endpoints")

	for _, v := range eps {
		// log.Trace().Str("endpoint", v.String()).Msg("endpoint selected")
		if err := f(v); err != nil {
			return err
		}
	}
	return nil
}
