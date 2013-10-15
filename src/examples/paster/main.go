package guestbook

import (
	"appengine/datastore"
	"appengine/user"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/russross/blackfriday"
	"hikaru"
	"html/template"
	"strconv"
	"strings"
	"time"
)

const (
	FormatMarkdown = iota
	FormatRestructuredtext
)

const (
	StatusPublic = iota
	StatusPrivate
)

const (
	KindPaste = "P"
)

var (
	ErrPasteFormatWrong     = errors.New("wrong format")
	ErrPasteStatusWrong     = errors.New("status format")
	ErrPasteTitleEmpty      = errors.New("title is empty")
	ErrPasteTitleTooLoong   = errors.New("title is too long")
	ErrPasteContentEmpty    = errors.New("content is empty")
	ErrPasteContentTooLoong = errors.New("content is too long")
)

var app = hikaru.NewApplication()

func init() {
	app.Debug = true
	app.RouteFunc("/", handleWelcome)
	app.RouteFunc("/mypage", handleMypage)
	app.RouteFunc("/new", handleEdit)
	app.RouteFunc("/p/<id:\\w+>", handleShow)
	app.RouteFunc("/p/<id:\\w+>/edit", handleEdit)
	app.Start()
}

type HtmlRendererValue map[string]interface{}

func processHtmlRendererValue(arg HtmlRendererValue) HtmlRendererValue {
	return arg
}

type Paste struct {
	Key     *datastore.Key `datastore:"-"`
	User    string         `datastore:"u,"`
	Title   string         `datastore:"t,noindex"`
	Content string         `datastore:"c,noindex"`
	Format  int            `datastore:"f,noindex"`
	Status  int            `datastore:"s,"`
	Updated time.Time      `datastore:"tu,"`
	Created time.Time      `datastore:"tc,"`
}

func NewPaste(c hikaru.Context) *Paste {
	u := user.Current(c)
	if u == nil || u.ID == "" {
		panic("invalid user")
	}
	now := time.Now()
	return &Paste{
		Key:     datastore.NewIncompleteKey(c, KindPaste, nil),
		User:    u.ID,
		Format:  FormatMarkdown,
		Status:  StatusPublic,
		Updated: now,
		Created: now,
	}
}

func (p *Paste) ID() int64 {
	return p.Key.IntID()
}

func (p *Paste) EncodedId() string {
	buf := make([]byte, 8)
	binary.PutVarint(buf, p.ID())
	ret := base64.URLEncoding.EncodeToString(buf)
	return strings.TrimRight(ret, "=")
}

func (p *Paste) ContentHTML() template.HTML {
	input := []byte(p.Content)
	content := blackfriday.MarkdownCommon(input)
	return template.HTML(content)
}

func (p *Paste) Put(c hikaru.Context) error {
	if p.Key == nil {
		return errors.New("Key is nil")
	}
	key, err := datastore.Put(c, p.Key, p)
	if err == nil {
		p.Key = key
	}
	return err
}

func (p *Paste) AllowEdit(u *user.User) bool {
	if u == nil {
		return false
	}
	return p.User == u.ID
}

func (p *Paste) AllowRead(u *user.User) bool {
	if p.Status == StatusPublic {
		return true
	}
	return p.AllowEdit(u)
}

func (p *Paste) SetFormat(value string) error {
	v, err := strconv.Atoi(value)
	if err != nil {
		return ErrPasteFormatWrong
	} else if v != FormatMarkdown && v != FormatRestructuredtext {
		return ErrPasteFormatWrong
	}
	p.Format = v
	return nil
}

func (p *Paste) SetStatus(value string) error {
	v, err := strconv.Atoi(value)
	if err != nil {
		return ErrPasteStatusWrong
	} else if v != StatusPublic && v != StatusPrivate {
		return ErrPasteStatusWrong
	}
	p.Status = v
	return nil
}

func (p *Paste) SetTitle(value string) error {
	if value == "" {
		return ErrPasteTitleEmpty
	} else if len(value) > 50 {
		return ErrPasteTitleTooLoong
	}
	p.Title = value
	return nil
}

func (p *Paste) SetContent(value string) error {
	if value == "" {
		return ErrPasteContentEmpty
	} else if len(value) > 100000 {
		return ErrPasteContentTooLoong
	}
	p.Content = value
	return nil
}

func PasteLatests(c hikaru.Context, u *user.User, limit int) ([]Paste, error) {
	if u == nil || u.ID == "" {
		return nil, errors.New("invalid user")
	}
	q := datastore.NewQuery(KindPaste).
		Filter("u =", u.ID).
		Order("-tu").
		Limit(limit)
	entities := make([]Paste, 0, limit)
	it := q.Run(c)
	for {
		p := Paste{}
		key, err := it.Next(&p)
		if err != nil {
			if err == datastore.Done {
				break
			}
			return nil, err
		}
		p.Key = key
		entities = append(entities, p)
	}
	return entities, nil
}

func PasteGetById(c hikaru.Context, id int64) (*Paste, error) {
	key := datastore.NewKey(c, KindPaste, "", id, nil)
	p := new(Paste)
	if err := datastore.Get(c, key, p); err != nil {
		if err != datastore.ErrNoSuchEntity {
			return nil, err
		}
		return nil, nil
	}
	p.Key = key
	return p, nil
}

func decodeId(id string) int64 {
	mod := len(id) % 4
	if mod != 0 {
		id += strings.Repeat("=", mod)
	}
	buf, err := base64.URLEncoding.DecodeString(id)
	if err != nil {
		panic(err)
	}
	ret, _ := binary.Varint(buf)
	return ret
}

func handleWelcome(c hikaru.Context) hikaru.Result {
	return c.Html("layout", "welcome", map[string]interface{}{
		"c":    c,
		"user": user.Current(c),
	})
}

func handleMypage(c hikaru.Context) hikaru.Result {
	u, res := requireUser(c)
	if res != nil {
		return res
	}
	entities, err := PasteLatests(c, u, 10)
	if err != nil {
		return c.Abort(err)
	}
	return c.Html("layout", "mypage", map[string]interface{}{
		"c":        c,
		"user":     u,
		"entities": entities,
	})
}

func handleShow(c hikaru.Context) hikaru.Result {
	u := user.Current(c)
	p, err := PasteGetById(c, decodeId(c.Val("id")))
	if err != nil {
		return c.Abort(err)
	}
	if p == nil {
		return c.NotFound()
	}
	if !p.AllowRead(u) {
		return c.AbortCode(403)
	}
	return c.Html("layout", "show", map[string]interface{}{
		"c":      c,
		"user":   u,
		"entity": p,
	})
}

func handleEdit(c hikaru.Context) hikaru.Result {
	u, res := requireUser(c)
	if res != nil {
		return res
	}

	var err error
	var p *Paste

	if c.Has("id") {
		p, err = PasteGetById(c, decodeId(c.Val("id")))
		if err != nil {
			return c.Abort(err)
		} else if p == nil {
			return c.NotFound()
		} else if !p.AllowEdit(u) {
			return c.AbortCode(403)
		}
	} else {
		p = NewPaste(c)
	}

	statuses := map[int]string{StatusPublic: "公開", StatusPrivate: "非公開"}

	if c.IsGet() {
		return c.Html("layout", "edit", map[string]interface{}{
			"c":        c,
			"user":     u,
			"path":     c.HttpRequest().RequestURI,
			"entity":   p,
			"errs":     nil,
			"statuses": statuses,
		})
	}

	errs := make(map[string]error)
	p.Format = FormatMarkdown
	//err = p.SetFormat(c.Form("format"))
	//if err != nil {
	//	errs["format"] = err
	//}
	err = p.SetStatus(c.Form("status"))
	if err != nil {
		errs["status"] = err
	}
	err = p.SetTitle(c.Form("title"))
	if err != nil {
		errs["title"] = err
	}
	err = p.SetContent(c.Form("content"))
	if err != nil {
		errs["content"] = err
	}
	if len(errs) >= 1 {
		return c.Html("layout", "edit", map[string]interface{}{
			"c":        c,
			"user":     u,
			"path":     c.HttpRequest().RequestURI,
			"entity":   p,
			"errs":     errs,
			"statuses": statuses,
		})
	}

	err = p.Put(c)
	if err != nil {
		return c.Abort(err)
	}
	return c.Redirect(fmt.Sprintf("/p/%s", p.EncodedId()))
}

func requireUser(c hikaru.Context) (*user.User, hikaru.Result) {
	u := user.Current(c)
	if u == nil {
		url, err := user.LoginURL(c, c.HttpRequest().URL.Path)
		if err != nil {
			return nil, c.Abort(err)
		}
		return nil, c.Redirect(url)
	}
	return u, nil
}
