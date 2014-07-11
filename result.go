package hikaru

import (
	"bytes"
	"net/http"
)

type Result struct {
	StatusCode int
	Header     http.Header
	Body       bytes.Buffer
}

// Creates and returns a new Result.
func NewResult() *Result {
	r := &Result{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
	}
	return r
}

// Sets a response header.
func (r *Result) SetHeader(key, value string) {
	r.Header.Set(key, value)
}

// Adds a response header.
func (r *Result) AddHeader(key, value string) {
	r.Header.Add(key, value)
}

// Adds a cookie header.
func (r *Result) SetCookie(cookie *http.Cookie) {
	r.Header.Set("Set-Cookie", cookie.String())
}

func (r *Result) Flush(w http.ResponseWriter, req *http.Request) {
	copyHttpHeader(r.Header, w.Header())
	loc := r.Header.Get("Location")
	if loc != "" {
		http.Redirect(w, req, loc, r.StatusCode)
	} else {
		if r.StatusCode > 0 {
			w.WriteHeader(r.StatusCode)
		}
		if r.Body.Len() > 0 {
			r.Body.WriteTo(w)
		}
	}
}

func copyHttpHeader(src, dst http.Header) {
	if src != nil {
		for k, vs := range src {
			if len(vs) >= 2 {
				for _, v := range vs {
					if v != "" {
						dst.Add(k, v)
					}
				}
			} else {
				v := vs[0]
				if v != "" {
					dst.Set(k, v)
				}
			}
		}
	}
}
