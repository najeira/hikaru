package hikaru

import (
	"testing"
	"bytes"
	"net/http"
	_ "./testbed"
)

func TestApplication(t *testing.T) {
	app := NewApplication()
	if app == nil {
		t.Errorf("NewApplication returns nil")
	}
	if app.Routes == nil {
		t.Errorf("Application.Routes is nil")
	}
	if app.StaticDir != "static" {
		t.Errorf("Application.StaticDir is not static")
	}
	if app.TemplateDir != "templates" {
		t.Errorf("Application.TemplateDir is not templates")
	}
	if app.TemplateExt != "html" {
		t.Errorf("Application.TemplateExt is not html")
	}
	if app.LogLevel != LogLevelInfo {
		t.Errorf("Application.LogLevel is not LogLevelInfo")
	}
}

func TestApplicationRoute(t *testing.T) {
	app := NewApplication()
	f := func(c Context) Result { return nil }
	
	var r *http.Request
	var rd *RouteData
	
	app.RouteFunc("/", f)
	r, _ = http.NewRequest("GET", "http://example.com/", nil)
	rd = app.Match(r)
	if rd == nil {
		t.Errorf("/ does not match /")
	}
	if rd.Params == nil {
		t.Errorf("/ return wrong RouteData")
	}
	if len(rd.Params) != 0 {
		t.Errorf("/ return wrong RouteData")
	}
	if rd.Route.Handler() == nil {
		t.Errorf("/ matches wrong Route")
	}
	if rd.Route.Timeout() != 0 {
		t.Errorf("/ matches wrong Route")
	}
	
	r, _ = http.NewRequest("GET", "http://example.com/dummy", nil)
	rd = app.Match(r)
	if rd != nil {
		t.Errorf("/dummy matches /")
	}
	
	app.RouteFunc("/test", f)
	r, _ = http.NewRequest("GET", "http://example.com/test", nil)
	rd = app.Match(r)
	if rd == nil {
		t.Errorf("/test does not match /test")
	}
	if rd.Params == nil {
		t.Errorf("/test return wrong RouteData")
	}
	if len(rd.Params) != 0 {
		t.Errorf("/test return wrong RouteData")
	}
	if rd.Route.Handler() == nil {
		t.Errorf("/test matches wrong Route")
	}
	if rd.Route.Timeout() != 0 {
		t.Errorf("/test matches wrong Route")
	}
	
	app.RouteFunc("/test/<foo>", f)
	r, _ = http.NewRequest("GET", "http://example.com/test/bar", nil)
	rd = app.Match(r)
	if rd == nil {
		t.Errorf("/test/bar does not match /test/bar")
	}
	if rd.Params == nil {
		t.Errorf("/test/bar return wrong RouteData")
	}
	if len(rd.Params) != 1 {
		t.Errorf("/test/bar return wrong RouteData")
	}
	if rd.Params["foo"] != "bar" {
		t.Errorf("/test/bar return wrong RouteData")
	}
	if rd.Route.Handler() == nil {
		t.Errorf("/test/bar matches wrong Route")
	}
	if rd.Route.Timeout() != 0 {
		t.Errorf("/test/bar matches wrong Route")
	}
	
	app.RouteFunc("/hoge/<id:[0-9]+>", f)
	r, _ = http.NewRequest("GET", "http://example.com/hoge/fuga", nil)
	rd = app.Match(r)
	if rd != nil {
		t.Errorf("/hoge/fuga matches /hoge/<id:[0-9]+>")
	}
	
	r, _ = http.NewRequest("GET", "http://example.com/hoge/579", nil)
	rd = app.Match(r)
	if rd == nil {
		t.Errorf("/hoge/579 does not match /hoge/<id:[0-9]+>")
	}
	if rd.Params == nil {
		t.Errorf("/hoge/579 return wrong RouteData")
	}
	if len(rd.Params) != 1 {
		t.Errorf("/hoge/579 return wrong RouteData")
	}
	if rd.Params["id"] != "579" {
		t.Errorf("/hoge/579 return wrong RouteData: %s", rd.Params["id"])
	}
	if rd.Route.Handler() == nil {
		t.Errorf("/hoge/579 matches wrong Route")
	}
	if rd.Route.Timeout() != 0 {
		t.Errorf("/hoge/579 matches wrong Route")
	}
}

// stub
type ResponseTester struct {
	header      http.Header
	body        bytes.Buffer
	status_code int
}

func NewResponseTester() *ResponseTester {
	r := &ResponseTester{}
	r.header = make(http.Header)
	return r
}

func (r *ResponseTester) Header() http.Header {
	return r.header
}

func (r *ResponseTester) Write(body []byte) (int, error) {
	return r.body.Write(body)
}

func (r *ResponseTester) WriteHeader(code int) {
	r.status_code = code
}

func TestContext(t *testing.T) {
	app := NewApplication()
	
	var w http.ResponseWriter
	var r *http.Request
	
	w = NewResponseTester()
	r, _ = http.NewRequest("GET", "http://example.com/", nil)
	
	c := NewContext(app, w, r)
	if c == nil {
		t.Errorf("NewContext returns nil")
	}
	if c.Application() != app {
		t.Errorf("Context.Application() returns nil")
	}
	if c.AppEngineContext() == nil {
		t.Errorf("Context.AppEngineContext() returns nil")
	}
	if c.HttpRequest() != r {
		t.Errorf("Context.HttpRequest() returns wrong value")
	}
	if c.ResponseWriter() != w {
		t.Errorf("Context.ResponseWriter() returns wrong value")
	}
	if c.Method() != "GET" {
		t.Errorf("Context.Method() returns not GET")
	}
	if c.RouteData() != nil {
		t.Errorf("Context.RouteData() returns not nil value")
	}
	if c.Result() != nil {
		t.Errorf("Context.RouteData() returns not nil value")
	}
	if c.IsMethodPost() {
		t.Errorf("Context.IsMethodPost() returns true")
	}
	if !c.IsMethodGet() {
		t.Errorf("Context.IsMethodGet() returns false")
	}
}

func TestContextRoute(t *testing.T) {
	app := NewApplication()
	f := func(c Context) Result { return nil }
	app.RouteFunc("/", f)
	app.RouteFunc("/test/<foo>", f)
	
	var w http.ResponseWriter
	var r *http.Request
	var c *HikaruContext
	
	w = NewResponseTester()
	r, _ = http.NewRequest("GET", "http://example.com/", nil)
	c = NewContext(app, w, r)
	
	matched := c.executeRoute()
	if !matched {
		t.Errorf("Context.executeRoute does not match to /")
	}
	if c.RouteData() == nil {
		t.Errorf("Context.RouteData() returns nil")
	}
	
	w = NewResponseTester()
	r, _ = http.NewRequest("GET", "http://example.com/test/bar", nil)
	c = NewContext(app, w, r)
	
	matched = c.executeRoute()
	if !matched {
		t.Errorf("Context.executeRoute does not match to /")
	}
	if c.RouteData() == nil {
		t.Errorf("Context.RouteData() returns nil")
	}
	if !c.Has("foo") {
		t.Errorf("Context.Has(`foo`) returns false")
	}
	if c.Val("foo") == "" {
		t.Errorf("Context.Val(`foo`) returns ``")
	}
	if len(c.Vals("foo")) != 1 {
		t.Errorf("Context.Vals(`foo`) returns empty value")
	}
	
	var buf bytes.Buffer
	
	buf.WriteString("hoge=fuga")
	w = NewResponseTester()
	r, _ = http.NewRequest("POST", "http://example.com/", &buf)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c = NewContext(app, w, r)
	
	matched = c.executeRoute()
	if c.Form("hoge") == "" {
		t.Errorf("Context.Form(`hoge`) returns ``")
	}
	if len(c.Forms("hoge")) != 1 {
		t.Errorf("Context.Forms(`hoge`) returns empty value")
	}
}

func TestContextResult(t *testing.T) {
	app := NewApplication()
	f := func(c Context) Result { return nil }
	app.RouteFunc("/", f)
	
	var w http.ResponseWriter
	var r *http.Request
	var c *HikaruContext
	
	w = NewResponseTester()
	r, _ = http.NewRequest("GET", "http://example.com/", nil)
	c = NewContext(app, w, r)
	
	res := r.Raw("", "text/plain")
}
