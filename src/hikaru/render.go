package hikaru

import (
	"bytes"
	"html/template"
	"path"
	"strings"
)

type Renderer interface {
	Render(string, interface{}) string
}

type HikaruRenderer struct {
	dir       string
	ext       string
	templates map[string]*template.Template
}

func NewRenderer(dir string, ext string) *HikaruRenderer {
	r := &HikaruRenderer{dir: dir, ext: ext}
	r.templates = make(map[string]*template.Template)
	return r
}

func (r *HikaruRenderer) getFullname(name string) string {
	name = path.Clean(name)
	if r.ext != "" {
		if !strings.HasPrefix(r.ext, ".") {
			name += "."
		}
		name += r.ext
	}
	return path.Clean(path.Join("./", r.dir, name))
}

func (r *HikaruRenderer) New(name string) *template.Template {
	fullname := r.getFullname(name)
	tpl, err := template.ParseFiles(fullname)
	if err != nil {
		panic(err)
	}
	return tpl
}

func (r *HikaruRenderer) Get(name string) *template.Template {
	tpl, ok := r.templates[name]
	if !ok || tpl == nil {
		return r.New(name)
	}
	return tpl
}

func (r *HikaruRenderer) Render(name string, data interface{}) string {
	return r.RendererTemplate(r.Get(name), data)
}

func (r *HikaruRenderer) RendererTemplate(tpl *template.Template, data interface{}) string {
	var buf bytes.Buffer
	err := tpl.Execute(&buf, data)
	if err != nil {
		panic(err)
	}
	return string(buf.Bytes())
}
