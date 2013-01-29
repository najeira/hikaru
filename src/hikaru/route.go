package fimika

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type Route struct {
	Pattern string
	Regexp  *regexp.Regexp
	Handler Handler
	Timeout time.Duration
}

type RouteData struct {
	Route  *Route
	Params map[string]string
}

type Handler func(*Context) Result

var (
	routeParamRegexp *regexp.Regexp = regexp.MustCompile("<[^>]+>")
)

func NewRoute(pattern string, handler Handler) *Route {
	r := &Route{Pattern: pattern, Handler: handler}
	r.Regexp = compileRoutePattern(r.Pattern)
	return r
}

func (r *Route) Match(r *http.Request) *RouteData {
	return r.MatchPath(r.URL.Path)
}

func (r *Route) MatchURL(url *url.URL) *RouteData {
	return r.MatchPath(url.Path)
}

func (r *Route) MatchPath(path string) *RouteData {
	params := NamedRegexpGroup(path, r.Regexp)
	if params == nil {
		return nil
	}
	rd := &RouteData{Route: r, Params: params}
	return rd
}

func compileRoutePattern(pattern string) *regexp.Regexp {
	// /blog/<id>           =>   /blog/(?P<id>[^\?#/]+)
	// /blog/<id:[0-9]+>    =>   /blog/(?P<id>[0-9]+)
	ret := routeParamRegexp.ReplaceAllStringFunc(pattern, func(s string) string {
		var name, reg string
		body := s[1 : len(s)-1]
		if strings.Contains(body, ":") {
			parts := strings.SplitN(body, ":", 2)
			name, reg = parts[0], parts[1]
		} else {
			name = body
			reg = "[^\\?#/]+"
		}
		return fmt.Sprintf("(?P<%s>%s)", regexp.QuoteMeta(name), reg)
	})
	return regexp.MustCompile("^" + ret + "$")
}
