package hikaru

import (
	"bytes"
	"net/http"
)

type Result interface {
	StatusCode() int
	Header() http.Header
	Execute(Context)
}

type HikaruResult struct {
	statusCode int
	header     http.Header
	body       bytes.Buffer
	err        interface{}
}

// Creates and returns a new Result.
func NewResult() *HikaruResult {
	result := new(HikaruResult)
	result.header = make(http.Header)
	return result
}

func (r *HikaruResult) StatusCode() int {
	return r.statusCode
}

func (r *HikaruResult) Header() http.Header {
	return r.header
}

func (r *HikaruResult) Execute(c Context) {
	copyHttpHeader(r.header, c.ResponseWriter().Header())
	if r.header.Get("Location") != "" {
		r.redirect(c)
	} else {
		if r.statusCode > 0 {
			c.ResponseWriter().WriteHeader(r.statusCode)
		}
		if r.body.Len() > 0 {
			r.body.WriteTo(c.ResponseWriter())
		}
	}
}

func (r *HikaruResult) SetCookie(cookie *http.Cookie) {
}

func (r *HikaruResult) redirect(c Context) {
	http.Redirect(c.ResponseWriter(), c.HttpRequest(), r.header.Get("Location"), r.statusCode)
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
