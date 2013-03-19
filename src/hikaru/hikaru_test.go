package hikaru

import (
	"testing"
	"bytes"
	"net/http"
	"strings"
	"errors"
)

const (
	PYTHON    = `C:\Python\27\python.exe`
	TESTBED   = `C:\software\google_appengine_go\goroot\src\pkg\github.com\najeira\testbed\testbed.py`
	APPENGINE = `C:\software\google_appengine_go`
)

func TestApplication(t *testing.T) {
	app := NewApplication()
	if app == nil {
		t.Errorf("NewApplication should return not nil")
	}
	if app.Routes != nil {
		t.Errorf("Application.Routes should be nil")
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
	if app.Debug != false {
		t.Errorf("Application.Debug should be false")
	}
	if len(app.Renderers) != 0 {
		t.Errorf("Application.Renderers should be empty")
	}
}

type RendererTester struct {}

func (r *RendererTester) Render(arg ...interface{}) Result {
	return NewResult()
}

func TestApplicationSetRenderer(t *testing.T) {
	app := NewApplication()
	r := new(RendererTester)
	app.SetRenderer("html", r)
	if app.GetRenderer("html") != r {
		t.Errorf("Application.SetRenderer failed")
	}
}

func TestApplicationRouteFunc(t *testing.T) {
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
	if c.AppEngineContext() == nil {
		t.Errorf("Context.AppEngineContext() returns nil")
	}
	if c.Method() != "GET" {
		t.Errorf("Context.Method() returns not GET")
	}
	if c.Application() != app {
		t.Errorf("Context.Application() returns wrong value")
	}
	if c.HttpRequest() != r {
		t.Errorf("Context.HttpRequest() returns wrong value")
	}
	if c.ResponseWriter() != w {
		t.Errorf("Context.ResponseWriter() returns wrong value")
	}
	if c.RouteData() != nil {
		t.Errorf("Context.RouteData() returns not nil value")
	}
	if c.Result() != nil {
		t.Errorf("Context.RouteData() returns not nil value")
	}
	if c.IsPost() {
		t.Errorf("Context.IsPost() returns true")
	}
	if !c.IsGet() {
		t.Errorf("Context.IsGet() returns false")
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

func TestContextResultRaw(t *testing.T) {
	app := NewApplication()
	f := func(c Context) Result { return nil }
	app.RouteFunc("/", f)
	
	w := NewResponseTester()
	r, _ := http.NewRequest("GET", "http://example.com/", nil)
	c := NewContext(app, w, r)
	
	res := c.Raw([]byte("This is a test value"), "text/plain")
	if res == nil {
		t.Errorf("Context.Raw() returns nil")
	}
	if res.StatusCode() != http.StatusOK {
		t.Errorf("Result.StatusCode() should be http.StatusOK")
	}
	if ct := res.Header().Get("Content-Type"); ct != "text/plain" {
		t.Errorf("Result.Header()[Content-Type] should be text/plain: %s", ct)
	}
	res.Execute(c)
	
	body := w.body.String()
	if body != "This is a test value" {
		t.Errorf("Result should be This is a test value")
	}
}

func TestContextResultText(t *testing.T) {
	app := NewApplication()
	f := func(c Context) Result { return nil }
	app.RouteFunc("/", f)
	
	w := NewResponseTester()
	r, _ := http.NewRequest("GET", "http://example.com/", nil)
	c := NewContext(app, w, r)
	
	res := c.Text("This is a test value")
	if res == nil {
		t.Errorf("Context.Text() returns nil")
	}
	if res.StatusCode() != http.StatusOK {
		t.Errorf("Result.StatusCode() should be http.StatusOK")
	}
	if ct := res.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/plain") {
		t.Errorf("Result.Header()[Content-Type] should be text/plain: %s", ct)
	}
	res.Execute(c)
	
	body := w.body.String()
	if body != "This is a test value" {
		t.Errorf("Result should be This is a test value")
	}
}

func TestContextResultRedirect(t *testing.T) {
	app := NewApplication()
	f := func(c Context) Result { return nil }
	app.RouteFunc("/", f)
	
	w := NewResponseTester()
	r, _ := http.NewRequest("GET", "http://example.com/", nil)
	c := NewContext(app, w, r)
	
	res := c.Redirect("http://example.com/redirect")
	if res == nil {
		t.Errorf("Context.Redirect() returns nil")
	}
	if res.StatusCode() != http.StatusFound {
		t.Errorf("Result.StatusCode() should be http.StatusFound")
	}
	if res.Header().Get("Location") != "http://example.com/redirect" {
		t.Errorf("Result.Header() returns wrong value")
	}
	res.Execute(c)
}

func TestContextResultNotFound(t *testing.T) {
	app := NewApplication()
	f := func(c Context) Result { return nil }
	app.RouteFunc("/", f)
	
	w := NewResponseTester()
	r, _ := http.NewRequest("GET", "http://example.com/", nil)
	c := NewContext(app, w, r)
	
	res := c.NotFound()
	if res == nil {
		t.Errorf("Context.NotFound() returns nil")
	}
	if res.StatusCode() != http.StatusNotFound {
		t.Errorf("Result.StatusCode() should be http.StatusNotFound")
	}
	res.Execute(c)
}

func TestContextResultAbort(t *testing.T) {
	app := NewApplication()
	f := func(c Context) Result { return nil }
	app.RouteFunc("/", f)
	
	w := NewResponseTester()
	r, _ := http.NewRequest("GET", "http://example.com/", nil)
	c := NewContext(app, w, r)
	
	err := errors.New("This is a test error")
	res := c.Abort(err)
	if res == nil {
		t.Errorf("Context.Abort() returns nil")
	}
	if res.StatusCode() != http.StatusInternalServerError {
		t.Errorf("Result.StatusCode() should be http.StatusInternalServerError")
	}
	res.Execute(c)
}

func TestContextResultAbortCode(t *testing.T) {
	app := NewApplication()
	f := func(c Context) Result { return nil }
	app.RouteFunc("/", f)
	
	w := NewResponseTester()
	r, _ := http.NewRequest("GET", "http://example.com/", nil)
	c := NewContext(app, w, r)
	
	res := c.AbortCode(http.StatusInternalServerError)
	if res == nil {
		t.Errorf("Context.AbortCode() returns nil")
	}
	if res.StatusCode() != http.StatusInternalServerError {
		t.Errorf("Result.StatusCode() should be http.StatusInternalServerError")
	}
	res.Execute(c)
}

func TestContextResultPanic(t *testing.T) {
	app := NewApplication()
	f := func(c Context) Result { return nil }
	app.RouteFunc("/", f)
	
	app.LogLevel = LogLevelNo
	
	w := NewResponseTester()
	r, _ := http.NewRequest("GET", "http://example.com/", nil)
	c := NewContext(app, w, r)
	
	err := errors.New("This is a test error")
	res := c.resultPanic(err)
	if res == nil {
		t.Errorf("Context.resultPanic() returns nil")
	}
	if res.StatusCode() != http.StatusInternalServerError {
		t.Errorf("Result.StatusCode() should be http.StatusInternalServerError")
	}
	res.Execute(c)
}

func TestContextExecuteMatch(t *testing.T) {
	app := NewApplication()
	f := func(c Context) Result { return c.Text("OK") }
	app.RouteFunc("/", f)
	
	w := NewResponseTester()
	r, _ := http.NewRequest("GET", "http://example.com/", nil)
	c := NewContext(app, w, r)
	
	c.Execute()
	res := c.Result()
	
	if res == nil {
		t.Errorf("Context.Result() returns nil")
	}
	if res.StatusCode() != http.StatusOK {
		t.Errorf("Result.StatusCode() should be http.StatusOK")
	}
	if ct := res.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/plain") {
		t.Errorf("Result.Header()[Content-Type] should be text/plain: %s", ct)
	}
	res.Execute(c)
	
	body := w.body.String()
	if body != "OK" {
		t.Errorf("expect OK: ", body)
	}
}
