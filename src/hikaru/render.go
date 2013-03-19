package hikaru

import (
	"bytes"
	"errors"
	"html/template"
	"net/http"
	"path"
	"strings"
	"sync"
)

type Renderer interface {
	Render(...interface{}) Result
}

type HikaruRenderer struct {
	dir       string
	ext       string
	templates map[string]*template.Template
	mutex     sync.RWMutex
}

func NewRenderer(dir string, ext string) *HikaruRenderer {
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	r := &HikaruRenderer{dir: dir, ext: ext}
	r.templates = make(map[string]*template.Template)
	return r
}

func (r *HikaruRenderer) getFullnames(names []string) []string {
	rets := make([]string, len(names))
	for i, name := range names {
		rets[i] = r.getFullname(name)
	}
	return rets
}

func (r *HikaruRenderer) getFullname(name string) string {
	if r.ext != "" {
		name += r.ext
	}
	var full string
	if r.dir != "" {
		full = path.Join("./", r.dir, name)
	} else {
		full = path.Join("./", name)
	}
	return path.Clean(full)
}

func (r *HikaruRenderer) New(names []string) *template.Template {
	fullnames := r.getFullnames(names)
	tpl, err := template.ParseFiles(fullnames...)
	if err != nil {
		panic(err)
	}
	return tpl
}

func (r *HikaruRenderer) Get(names []string) *template.Template {
	key := strings.Join(names, " ")
	r.mutex.RLock()
	tpl, ok := r.templates[key]
	r.mutex.RUnlock()
	if !ok {
		tpl = r.New(names)
		r.mutex.Lock()
		defer r.mutex.Unlock()
		r.templates[key] = tpl
	}
	return tpl
}

func (r *HikaruRenderer) Render(args ...interface{}) Result {
	if len(args) <= 0 {
		panic(errors.New("no arguments"))
	}
	var data interface{}
	names := make([]string, 0)
	n := len(args)
	for i := 0; i < n; i++ {
		s, ok := args[i].(string)
		if !ok {
			data = args[i]
			break
		}
		names = append(names, s)
	}
	if len(names) <= 0 {
		panic(errors.New("no string arguments"))
	}
	text := r.rendererTemplate(r.Get(names), data)
	result := NewResult()
	result.statusCode = http.StatusOK
	result.body.WriteString(text)
	result.header.Set("Content-Type", "text/html; charset=utf-8")
	return result
}

func (r *HikaruRenderer) rendererTemplate(tpl *template.Template, data interface{}) string {
	var buf bytes.Buffer
	err := tpl.Execute(&buf, data)
	if err != nil {
		panic(err)
	}
	return string(buf.Bytes())
}
