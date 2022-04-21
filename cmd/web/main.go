package main

import (
	"database/sql"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
	"html/template"
	"learn-web/snippets/pkg/models/mysql"
	"log"
	"net/http"
	"os"
	"time"
)

type Config struct {
	Addr      string
	StaticDir string
	Dsn       string
	Secret    string
}

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	sessions      *sessions.Session
	snippets      *mysql.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	config := new(Config)

	flag.StringVar(&config.Addr, "addr", ":4000", "http network address")
	flag.StringVar(&config.StaticDir, "static-dir", "./ui/static", "path to static assets")

	flag.StringVar(&config.Dsn, "dns", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")

	flag.StringVar(&config.Secret, "secret", "secretkey", "Secret key")

	flag.Parse()

	// Use log.New() to create a logger for writing information messages. This takes
	// three parameters: the destination to write the logs to (os.Stdout), a string
	// prefix for message, and flags to indicate what additional information to include.
	// log.LUTC force your logger to use UTC datetimes instead of local ones.
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.LUTC)

	// create a logger for writing error messages in the same way, but use stderr as
	// the destination and use the log.Lshortfile flag to include the relevant
	// file name and line number.
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.LUTC|log.LUTC|log.Lshortfile)

	// open the database connection pool...
	db, err := openDB(config.Dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	// initialize the template cache
	templateCache, err := createTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	// Use the sessions.New() function to init a new session manager
	session := sessions.New([]byte(*&config.Secret))
	session.Lifetime = 12 * time.Hour

	// initialize a new instance of application containing the dependencies
	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		sessions:      session,
		snippets:      &mysql.SnippetModel{Db: db},
		templateCache: templateCache,
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

// openDB function wrap sql.Open and return a sql.DB connection pool
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
