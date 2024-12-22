package htmlRender

import (
	"html/template"
	"io"
	"strings"
)

type Templates struct {
	templates *template.Template
}

func NewTemplates() *Templates {
	funcMap := template.FuncMap{
		"safeID": func(s string) string {
			unsafe := []string{"/", "\\", " ", ".", ",", ":", ";"}
			result := s
			for _, ch := range unsafe {
				result = strings.ReplaceAll(result, ch, "-")
			}
			return result
		},
	}

	tmpl := template.New("")

	tmpl = template.Must(tmpl.Funcs(funcMap).ParseGlob("views/*.html"))

	return &Templates{
		templates: tmpl,
	}
}

func (t *Templates) Render(w io.Writer, name string, data interface{}) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
