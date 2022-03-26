package main

import (
	"html/template"
	"learn-web/snippets/pkg/models"
	"path/filepath"
)

type templateData struct {
	CurrentYear int
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
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
		ts, err := template.ParseFiles(page)
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
