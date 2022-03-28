package main

import (
	"github.com/justinas/alice"
	"net/http"
)

// its return a http.Handler instead of *http.ServerMux
func (app *application) routes() http.Handler {
	// create a middleware chain which will be used for every request our application receives.
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	mux := http.NewServeMux()

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// return the standard middleware chain followed by the servemux.
	return standardMiddleware.Then(mux)
}
