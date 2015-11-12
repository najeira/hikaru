package hikaru

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func responseObjects() (*httptest.ResponseRecorder, *http.Request, *Response, error) {
	wr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://example.com/", nil)
	if err != nil {
		return wr, req, nil, err
	}
	res := &Response{}
	res.init(wr, req)
	return wr, req, res, err
}

func TestResponse(t *testing.T) {
	wr, req, res, err := responseObjects()
	if err != nil {
		t.Error(err)
	}

	if res.size != -1 {
		t.Errorf("size invalid")
	}
	if res.Size() != -1 {
		t.Errorf("Size invalid")
	}
	if res.status != 200 {
		t.Errorf("status invalid")
	}
	if res.Status() != 200 {
		t.Errorf("Status invalid")
	}
	if res.ResponseWriter != wr {
		t.Errorf("ResponseWriter invalid")
	}
	if res.request != req {
		t.Errorf("request invalid")
	}
}

func TestResponseWriteHeader(t *testing.T) {
	wr, _, res, err := responseObjects()
	if err != nil {
		t.Error(err)
	}

	res.WriteHeader(300)
	if res.Written() {
		t.Errorf("Written invalid")
	}
	if res.Status() != 300 {
		t.Errorf("Status invalid")
	}
	if wr.Code != 200 {
		t.Errorf("WriteHeader invalid")
	}

	res.WriteHeader(301)
	if res.Written() {
		t.Errorf("Written invalid")
	}
	if res.Status() != 301 {
		t.Errorf("Status invalid")
	}
	if wr.Code != 200 {
		t.Errorf("WriteHeader invalid")
	}
}

func TestResponseWrite(t *testing.T) {
	wr, _, res, err := responseObjects()
	if err != nil {
		t.Error(err)
	}

	n, err := res.Write([]byte("hoge"))
	if err != nil {
		t.Error(err)
	}
	if n != 4 {
		t.Errorf("Write invalid")
	}
	if res.Size() != 4 {
		t.Errorf("Size invalid %d", res.Size())
	}
	if res.Status() != 200 {
		t.Errorf("Status invalid")
	}
	if wr.Code != 200 {
		t.Errorf("Code invalid")
	}
	if wr.Body.String() != "hoge" {
		t.Errorf("Body invalid")
	}

	n, err = res.Write([]byte(" fuga"))
	if err != nil {
		t.Error(err)
	}
	if n != 5 {
		t.Errorf("Write invalid")
	}
	if res.Size() != 9 {
		t.Errorf("Size invalid %d", res.Size())
	}
	if wr.Body.String() != "hoge fuga" {
		t.Errorf("Body invalid")
	}
}

func TestResponseRedirect(t *testing.T) {
	wr, _, res, err := responseObjects()
	if err != nil {
		t.Error(err)
	}

	res.Redirect("/hoge", 301)
	if !res.Written() {
		t.Errorf("Redirect invalid")
	}
	if res.Status() != 301 {
		t.Errorf("Status invalid")
	}
	if wr.Code != 301 {
		t.Errorf("Code invalid")
	}
}
