package main

import (
	"fmt"
	"net/http"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")

		next.ServeHTTP(w, r)
	})
}

// logger request middleware
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", request.RemoteAddr, request.Proto, request.Method, request.URL.RequestURI())

		next.ServeHTTP(writer, request)
	})
}

// recoverPanic method, recover the stack when trigger a panic on the application.
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		// create a deferred function (which will always be run in the event of a panic as Go unwinds the stack).
		defer func() {
			// use the builtin recover function to check if there has been a panic or not. If there has...
			if err := recover(); err != nil {
				// Set a "Connection: close" header on the response
				writer.Header().Set("Connection", "close")

				// call the app.serverError helper method to return a 500
				// internal server response.
				app.serverError(writer, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(writer, request)
	})
}
