package hikaru

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type Route interface {
	Match(*http.Request) *RouteData
	Handler() Handler
	Timeout() time.Duration
}

type HikaruRoute struct {
	Pattern string
	Regexp  *regexp.Regexp
	handler Handler
	timeout time.Duration
}

type RouteData struct {
	Route  Route
	Params map[string]string
}

type Handler func(Context) Result

var (
	routeParamRegexp *regexp.Regexp = regexp.MustCompile("<[^>]+>")
)

func NewRoute(pattern string, handler Handler) *HikaruRoute {
	r := &HikaruRoute{Pattern: pattern, handler: handler}
	r.Regexp = compileRoutePattern(r.Pattern)
	return r
}

func (r *HikaruRoute) Handler() Handler {
	return r.handler
}

func (r *HikaruRoute) Timeout() time.Duration {
	return r.timeout
}

func (r *HikaruRoute) Match(req *http.Request) *RouteData {
	return r.MatchPath(req.URL.Path)
}

func (r *HikaruRoute) MatchURL(url *url.URL) *RouteData {
	return r.MatchPath(url.Path)
}

func (r *HikaruRoute) MatchPath(path string) *RouteData {
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
			reg = "[^" + regexp.QuoteMeta("?#/") + "]+"
		}
		return fmt.Sprintf("(?P<%s>%s)", regexp.QuoteMeta(name), reg)
	})
	return regexp.MustCompile("^" + ret + "$")
}
