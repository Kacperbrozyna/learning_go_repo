package cyoa

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"
	"strings"
)

var tpl *template.Template

const viewTemplate = `
<!DOCTYPE html>
<html>
    <head>
        <meta charset="utf-8">
        <title> Choose Your Own Adventure </title>
    </head>
    <body>
        <h1>{{.Title}}</h1>
        {{range .Paragraph}}
            <p>{{.}}</p>
        {{end}}
        <ul>    
        {{range .Options}}
            <li>
                <a href="/{{.Arc}}">{{.Text}}</a>
            </li>
        {{end}}
        </ul>
    </body>
</html>`

type Story map[string]Chapter

type Chapter struct {
	Title     string          `json:"title"`
	Paragraph []string        `json:"story"`
	Options   []ChapterOption `json:"options"`
}

type ChapterOption struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

type handler struct {
	story    Story
	template *template.Template
	pathFunc func(r *http.Request) string
}

type HandlerOption func(h *handler)

func init() {
	tpl = template.Must(template.New("").Parse(viewTemplate))
}

func defaultPathFunc(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "" || path == "/" {
		path = "/intro"
	}

	return path[1:]
}

func WithTemplate(t *template.Template) HandlerOption {
	return func(h *handler) {
		h.template = t
	}
}

func WithPath(fn func(r *http.Request) string) HandlerOption {
	return func(h *handler) {
		h.pathFunc = fn
	}
}

func NewHandler(s Story, options ...HandlerOption) http.Handler {
	h := handler{s, tpl, defaultPathFunc}

	for _, opt := range options {
		opt(&h)
	}
	return h
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := h.pathFunc(r)

	if _, okay := h.story[path]; okay {
		template_error := h.template.Execute(w, h.story[path])
		if template_error != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
		}
		return
	}

	http.Error(w, "Chapter not found", http.StatusNotFound)
}

func ReadJsonStory(filepath string) (Story, error) {

	json_file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	defer json_file.Close()

	decoder := json.NewDecoder(json_file)

	var story Story

	if err := decoder.Decode(&story); err != nil {
		return nil, err
	}

	return story, nil
}
