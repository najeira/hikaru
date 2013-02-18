package hikaru

import (
	"testing"
	//"fmt"
	"bytes"
	"net/http"
	"appengine/datastore"
	"github.com/najeira/testbed"
)

func TestApplication(t *testing.T) {
	app := NewApplication()
	if app == nil {
		t.Errorf("NewApplication should return not nil")
	}
	if app.Routes == nil {
		t.Errorf("Application.Routes should not be nil")
	}
	if app.StaticDir != "static" {
		t.Errorf("Application.StaticDir should be static")
	}
	if app.TemplateDir != "templates" {
		t.Errorf("Application.TemplateDir should be templates")
	}
	if app.TemplateExt != "html" {
		t.Errorf("Application.TemplateExt should be html")
	}
	if app.LogLevel != LogLevelInfo {
		t.Errorf("Application.LogLevel should be LogLevelInfo")
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
		t.Errorf("/ should match /")
	}
	if rd.Params == nil {
		t.Errorf("RouteData.Params should not be nil")
	}
	if len(rd.Params) != 0 {
		t.Errorf("RouteData.Params should be empty")
	}
	if rd.Route.Handler() == nil {
		t.Errorf("RouteData.Handler() should not be nil")
	}
	if rd.Route.Timeout() != 0 {
		t.Errorf("Route.Timeout() should be 0")
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
		t.Errorf("RouteData.Params should not be nil")
	}
	if len(rd.Params) != 0 {
		t.Errorf("RouteData.Params should be empty")
	}
	if rd.Route.Handler() == nil {
		t.Errorf("RouteData.Handler() should not be nil")
	}
	if rd.Route.Timeout() != 0 {
		t.Errorf("Route.Timeout() should be 0")
	}
	
	app.RouteFunc("/test/<foo>", f)
	r, _ = http.NewRequest("GET", "http://example.com/test/bar", nil)
	rd = app.Match(r)
	if rd == nil {
		t.Errorf("/test/bar does not match /test/bar")
	}
	if rd.Params == nil {
		t.Errorf("RouteData.Params should not be nil")
	}
	if len(rd.Params) != 1 {
		t.Errorf("RouteData.Params should have one param")
	}
	if rd.Params["foo"] != "bar" {
		t.Errorf("RouteData.Params[`foo`] should be `bar`: %s", rd.Params["foo"])
	}
	if rd.Route.Handler() == nil {
		t.Errorf("RouteData.Handler() should not be nil")
	}
	if rd.Route.Timeout() != 0 {
		t.Errorf("Route.Timeout() should be 0")
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
		t.Errorf("RouteData.Params should not be nil")
	}
	if len(rd.Params) != 1 {
		t.Errorf("RouteData.Params should have one param")
	}
	if rd.Params["id"] != "579" {
		t.Errorf("RouteData.Params[`id`] should be `579`: %s", rd.Params["id"])
	}
	if rd.Route.Handler() == nil {
		t.Errorf("RouteData.Handler() should not be nil")
	}
	if rd.Route.Timeout() != 0 {
		t.Errorf("Route.Timeout() should be 0")
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
		t.Errorf("Context.Application() returns wrong value")
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
	
	var r *http.Request
	var c *HikaruContext
	
	w := NewResponseTester()
	r, _ = http.NewRequest("GET", "http://example.com/", nil)
	c = NewContext(app, w, r)
	
	res := c.Raw([]byte("This is a test value"), "text/plain")
	if res == nil {
		t.Errorf("Context.Raw(...) returns nil")
	}
	if res.StatusCode() != http.StatusOK {
		t.Errorf("Result.StatusCode() should be http.StatusOK")
	}
	if ct := res.Header().Get("Content-Type"); ct != "text/plain" {
		t.Errorf("Result.Header()['Content-Type] should be text/plain: %s", ct)
	}
	res.Execute(c)
	
	body := w.body.String()
	if body != "This is a test value" {
		t.Errorf("Result should be This is a test value")
	}
}

func TestTestbed(t *testing.T) {
	testbed.Start(`C:\Python\27\python.exe`, `C:\Program Files (x86)\Google\google_appengine`)
	defer testbed.Close()
	
	r, _ := http.NewRequest("GET", "http://example.com/", nil)
	c := testbed.NewContext(r)
	
	testbed.SetUp()
	
	low, high, err := datastore.AllocateIDs(c, "Test", nil, 10)
	if err != nil {
		t.Errorf("datastore.AllocateIDs returns error: %v", err)
	}
	if high - low != 10 {
		t.Errorf("datastore.AllocateIDs returns wrong values: %d, %d", low, high)
	}
	
	// teardown
	testbed.TearDowm()
}
