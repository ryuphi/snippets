package main

import (
	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
	"net/http"
)

// its return a http.Handler instead of *http.ServerMux
func (app *application) routes() http.Handler {
	// create a middleware chain which will be used for every request our application receives.
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	mux := pat.New()

	mux.Get("/", http.HandlerFunc(app.home))
	mux.Get("/snippets/:id", http.HandlerFunc(app.showSnippet))
	mux.Get("/snippets/create", http.HandlerFunc(app.createSnippet))

	// file server...
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static/", fileServer))

	// return the standard middleware chain followed by the servemux.
	return standardMiddleware.Then(mux)
}
