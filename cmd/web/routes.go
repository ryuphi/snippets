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

	// Create a new middleware. for now , this chain will only contain the session middleware but we'll
	// add more to it later.
	dynamicMiddleware := alice.New(app.sessions.Enable)

	mux := pat.New()

	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Get("/snippets", dynamicMiddleware.ThenFunc(app.home))
	mux.Post("/snippets", dynamicMiddleware.ThenFunc(app.createSnippet))
	mux.Get("/snippets/create", dynamicMiddleware.ThenFunc(app.createSnippetForm))
	mux.Get("/snippets/:id", dynamicMiddleware.ThenFunc(app.showSnippet)) // it must go here after /snippets/create

	// file server...
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static/", fileServer))

	// return the standard middleware chain followed by the servemux.
	return standardMiddleware.Then(mux)
}
