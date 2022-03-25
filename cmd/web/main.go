package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type Config struct {
	Addr      string
	StaticDir string
}

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	config := new(Config)

	flag.StringVar(&config.Addr, "addr", ":4000", "http network address")
	flag.StringVar(&config.StaticDir, "static-dir", "./ui/static", "path to static assets")

	flag.Parse()

	f, err := os.OpenFile("/tmp/info.log", os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	// Use log.New() to create a logger for writing information messages. This takes
	// three parameters: the destination to write the logs to (os.Stdout), a string
	// prefix for message, and flags to indicate what additional information to include.
	// log.LUTC force your logger to use UTC datetimes instead of local ones.
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.LUTC)

	// create a logger for writing error messages in the same way, but use stderr as
	// the destination and use the log.Lshortfile flag to include the relevant
	// file name and line number.
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.LUTC|log.LUTC|log.Lshortfile)

	// initialize a new instance of application containing the dependencies
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	// Initialize a new http.Server struct. We set the Addr and Handler fields so
	// that the server uses the same network address and routes as before and set
	// the ErrorLog field so that the server now uses the custom errorLog logger in
	// the event of any problems.
	server := &http.Server{
		Addr:     config.Addr,
		ErrorLog: errorLog,
		Handler:  app.routes(), // call the app.routes method.
	}

	infoLog.Printf("starting new server on %s", config.Addr)

	// Call the listenAndServer method on our new http.Server struct.
	err = server.ListenAndServe()

	errorLog.Fatal(err)
}
