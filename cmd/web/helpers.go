package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

// the serverError helper writes an error message and stack trace to the errorlog.
// then sends a generic 500 internal server error response to the user
func (app *application) serverError(responseWriter http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(responseWriter, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// the clientError helper sends a specific status code and corresponding description
// to the user. We'll use this later in the book to send responses like 400 'bad request'
// when there's a problem with the request that the user sent.
func (app *application) clientError(responseWriter http.ResponseWriter, status int) {
	http.Error(responseWriter, http.StatusText(status), status)
}

// For consistency, we'll also implement a notFound helper. This is simply a convenience wrapper around clientError
// which sends a 404 Not found response to the user.
func (app *application) notFound(responseWriter http.ResponseWriter) {
	app.clientError(responseWriter, http.StatusNotFound)
}

// add common and default data to the templateData.
func (app *application) addDefaultData(td *templateData, request *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}
	td.CurrentYear = time.Now().Year()

	// Add the flash message to the template data, if one exists...
	td.Flash = app.sessions.PopString(request, "flash")
	return td
}

// get the template from the template cache map
func (app *application) render(responseWriter http.ResponseWriter, request *http.Request, name string, data *templateData) {
	// retrieve the appropriate template set from the cache
	templateSet, ok := app.templateCache[name]
	if !ok {
		app.serverError(responseWriter, fmt.Errorf("the template %s does not exist", name))
		return
	}

	// init a new Buffer
	buffer := new(bytes.Buffer)

	// execute the template set passing in any dynamic data.
	//
	// write the template to the buffer, instead of straight to the http.ResponseWriter. If there's an error,
	// call our serverError helper and the return.
	//
	// Add common/default data to the template data.
	err := templateSet.Execute(buffer, app.addDefaultData(data, request))
	if err != nil {
		app.serverError(responseWriter, err)
		return
	}

	// write the content of the buffer to the http.ResponseWriter
	buffer.WriteTo(responseWriter)
}
