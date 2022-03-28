package main

import (
	"html/template"
	"learn-web/snippets/pkg/forms"
	"learn-web/snippets/pkg/models"
	"path/filepath"
	"time"
)

type templateData struct {
	CurrentYear int
	Form        *forms.Form
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

// Create a humanDate function which returns a nicely formatted string
// representation of a time.Time objects
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

func createTemplateCache(dir string) (map[string]*template.Template, error) {
	// init a new map to act as the in memory cache.
	cache := map[string]*template.Template{}

	// Use the filepath.Glob function to get a slice of all filepaths with
	// the extension '.page.tmpl'. This essentially gives us a slice of all the
	// 'page' templates for the application.
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	// loop through the pages one-by-one
	for _, page := range pages {
		// extract the file name from the full file path
		// and assign it to the name variable.
		name := filepath.Base(page)

		// parse the page template file in to a template set.
		// also the template.FuncMap must be registered with the template set before you
		// call the ParseFiles method. This means we hace to use template.New to create an empty
		// template set, use the Funcs method to register the template.FuncMap, and then parse the file as normal
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Use ParseGlob to add any 'layout' templates to the template set
		// (in this case, it's just the base layout at the moment)
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		// Use ParseGlob to add any partial templates to the template set
		// (in this case, it's just the footer partial at the moment)
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	// return the cache map
	return cache, nil
}
