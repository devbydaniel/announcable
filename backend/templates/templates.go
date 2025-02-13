package templates

import (
	"embed"
	"html/template"

	"github.com/devbydaniel/release-notes-go/internal/logger"
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
