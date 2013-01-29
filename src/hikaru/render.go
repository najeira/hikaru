package hikaru

import (
	"html/template"
	"io"
	"path"
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
	r := &HikaruRenderer{dir: dir, ext, ext}
	r.templates = make(map[string]*html.Template)
	return r
}

func (r *HikaruRenderer) getFullname(name string) string {
	name = path.Clean(name)
	if r.ext != "" {
		if r.ext[0] != "." {
			name += "."
		}
		name += r.ext
	}
	return path.Join(r.dir, name)
}

func (r *HikaruRenderer) New(name string) *template.Template {
	fullname := r.getFullname(name)
	tpl, err := template.ParseFile(fullname)
	if err != nil || tpl == nil {
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
	err := tpl.Render(buf, data) // FIXME
	if err != nil {
		panic(err)
	}
	return buf.String()
}
