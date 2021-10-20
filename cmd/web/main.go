package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/freischarler/unityapp/pkg/models/mysql"
	_ "github.com/go-sql-driver/mysql" // New import
)

// Define an application struct to hold the application-wide dependencies for the
// web application. For now we'll only include fields for the two custom logger
// we'll add more to it as the build progresses.

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *mysql.SnippetModel
	templateCache map[string]*template.Template
}

//command to [go run .]  or [go run $(ls *.go)]

//main does:
//Parsing the runtime configuration settings for the application
//Establishing the dependencies for the handlers and
//Running the HTTP server
func main() {
	// Define a new command-line flag with the name 'addr', a default value of
	// and some short help text explaining what the flag controls. The value of
	// flag will be stored in the addr variable at runtime.

	//go run $(ls *.go) -addr=":9999"
	//go run $(ls *.go) -help
	//ports 0-1023 are restricted
	addr := flag.String("addr", ":8000", "HTTP network address")

	// Define a new command-line flag for the MySQL DSN string.
	dsn := flag.String("dsn", "root:root@tcp(127.0.0.1:3306)/unity", "MySQL database name")
	flag.Parse()

	// Importantly, we use the flag.Parse() function to parse the command-line
	// This reads in the command-line flag value and assigns it to the addr
	// variable. You need to call this *before* you use the addr variable
	// otherwise it will always contain the default value of ":4000". If any error
	// encountered during parsing the application will be terminated.
	flag.Parse()

	// Use log.New() to create a logger for writing information messages. This
	// three parameters: the destination to write the logs to (os.Stdout), a st
	// prefix for message (INFO followed by a tab), and flags to indicate what
	// additional information to include (local date and time). Note that the flag
	// are joined using the bitwise OR operator |.
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	// Create a logger for writing error messages in the same way, but use stde
	// the destination and use the log.Lshortfile flag to include the relevant
	// file name and line number.

	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// To keep the main() function tidy I've put the code for creating a connec
	// pool into the separate openDB() function below. We pass openDB() the DSN
	// from the command-line flag.
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	// We also defer a call to db.Close(), so that the connection pool is closed
	// before the main() function exits.
	defer db.Close()

	// Initialize a new template cache...
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	// Initialize a new instance of application containing the dependencies.
	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &mysql.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(), //Handler:  mux,
	}

	// Use the http.ListenAndServe() function to start a new web server. We pas
	// two parameters: the TCP network address to listen on (in this case ":4000
	// and the servemux we just created. If http.ListenAndServe() returns an er
	// we use the log.Fatal() function to log the error message and exit.
	//When our server receives a new HTTP request, it calls the servemux’s ServeHTTP() method.

	// Write messages using the two new loggers, instead of the standard logger
	infoLog.Printf("Starting server on %s", *addr) //log.Printf("Starting server on %s", *addr)

	err = srv.ListenAndServe() //err := http.ListenAndServe(*addr, mux)
	errorLog.Fatal(err)        //log.Fatal(err)
}

// The openDB() function wraps sql.Open() and returns a sql.DB connection pool
// for a given DSN.
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
