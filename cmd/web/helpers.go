package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
)

// the serverError helper writes an error message and stack trace to the errorlog.
// then sends a generic 500 internal server error response to the user
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// the clientError helper sends a specific status code and corresponding description
// to the user. We'll use this later in the book to send responses like 400 'bad request'
// when there's a problem with the request that the user sent.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// For consistency, we'll also implement a notFound helper. This is simply a convenience wrapper around clientError
// which sends a 404 Not found response to the user.
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

// get the template from the template cache map
func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	// retrieve the appropriate template set from the cache
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("the template %s does not exist", name))
		return
	}

	// init a new Buffer
	buffer := new(bytes.Buffer)

	// execute the template set passing in any dynamic data.
	// write the template to the buffer, instead of straight to the http.ResponseWriter. If there's an error,
	// call our serverError helper and the return.
	err := ts.Execute(buffer, td)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// write the content of the buffer to the http.ResponseWriter
	buffer.WriteTo(w)
}
