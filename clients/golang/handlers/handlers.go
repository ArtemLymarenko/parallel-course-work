package handlers

import (
	"golang/app"
	htmlRender "golang/html_render"
)

type Handlers struct {
	env  app.Env
	tmpl *htmlRender.Templates
}

func New(env app.Env, tmpl *htmlRender.Templates) *Handlers {
	return &Handlers{
		env:  env,
		tmpl: tmpl,
	}
}
