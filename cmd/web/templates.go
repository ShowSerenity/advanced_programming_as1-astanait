package main

import (
	"html/template"
	"path/filepath"
	"showserenity.net/astanait/pkg/forms"
	"showserenity.net/astanait/pkg/models"
	"time"
)

type templateData struct {
	CSRFToken       string
	CurrentYear     int
	Flash           string
	NewsType        string
	Form            *forms.Form
	News            *models.News
	Newses          []*models.News
	Accountant      *models.Accountant
	Accountants     []*models.Accountant
	IsAuthenticated bool
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.html"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.html"))
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}
	return cache, nil
}
