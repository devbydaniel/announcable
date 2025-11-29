package templates

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/devbydaniel/announcable/internal/logger"
)

var log = logger.Get()

//go:embed layouts/* pages/* partials/*
var templates embed.FS

func Construct(name string, files ...string) *template.Template {
	log.Trace().Str("name", name).Interface("files", files).Msg("Construct")
	withPartials := append(files, "partials/*")
	return template.Must(template.New(name).ParseFS(templates, withPartials...))
}

func Get(name string) (*template.Template, error) {
	log.Trace().Str("name", name).Msg("Get")
	return template.ParseFS(templates, name)
}

// Add a new function to execute templates with common data
func ExecuteTemplate(tmpl *template.Template, w http.ResponseWriter, name string, data interface{}) error {
	log.Trace().Str("name", name).Interface("data", data).Msg("ExecuteTemplate")
	return tmpl.ExecuteTemplate(w, name, data)
}
