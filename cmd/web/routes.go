package main

import "net/http"

// its return a http.Handler instead of *http.ServerMux
func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// pass the servermux as the 'next' parameter to the secureHeaders middleware.
	// because secureHeaders is just a function, and the function returns a
	// http.Handler we don't need to do anything else.
	return secureHeaders(mux)
}
