package templates

import (
	"context"
	"html/template"
	"io"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c context.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewTemplates(path string) *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob(path)),
	}
}
