package main

import (
	"darkhanomirbay.net/aitunews/pkg/models/postgresql"
	"database/sql"
	"flag"
	"fmt"
	"github.com/golangcollege/sessions"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	articles      *postgresql.ArticleModel
	templateCache map[string]*template.Template
	session       *sessions.Session
	users         *postgresql.UserModel
}

const (
	host     = "localhost"
	port     = "5432"
	user     = "postgres"
	password = "703905"
	dbname   = "AituNews"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	//db, err := sql.Open("postgres", psqlInfo)
	//if err != nil {
	//	panic(err)
	//}
	//defer db.Close()
	//err = db.Ping()
	//if err != nil {
	//	panic(err)
	//}

	fmt.Println("Successfully connected!")
	// for example go run ./cmd/web -addr=":9999"
	// Define a new command-line flag with the name 'addr', a default value of ":4000"
	// and some short help text explaining what the flag controls. The value of the
	// flag will be stored in the addr variable at runtime.
	addr := flag.String("addr", ":4000", "HTTP network address")
	// Importantly, we use the flag.Parse() function to parse the command-line flag.
	// This reads in the command-line flag value and assigns it to the addr
	// variable. You need to call this *before* you use the addr variable
	// otherwise it will always contain the default value of ":4000". If any errors are
	// encountered during parsing the application will be terminated.
	dsn := flag.String("dsn", psqlInfo, "PostgreSQL data source name")
	flag.Parse()
	// Use the http.NewServeMux() function to initialize a new servemux, then
	// register the home function as the handler for the "/" URL pattern.
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret key")
	flag.Parse()

	// new log
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	// error log
	errorLog := log.New(os.Stderr, "INFO\t", log.Ldate|log.Ltime)
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTemplateCache("./ui/html")
	if err != nil {
		errorLog.Fatal(err)
	}

	// Use the sessions.New() function to initialize a new session manager,
	// passing in the secret key as the parameter. Then we configure it so
	// sessions always expires after 12 hours.
	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	session.Secure = true // Set the Secure flag on our session cookies
	session.SameSite = http.SameSiteStrictMode

	var app = &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		articles:      &postgresql.ArticleModel{DB: db},
		templateCache: templateCache,
		session:       session,
		users:         &postgresql.UserModel{DB: db},
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}
	// The value returned from the flag.String() function is a pointer to the flag
	// value, not the value itself. So we need to dereference the pointer (i.e.
	// prefix it with the * symbol) before using it. Note that we're using the
	// log.Printf() function to interpolate the address with the log message.
	infoLog.Printf("Starting server on %s", *addr)
	//// Use the ListenAndServeTLS() method to start the HTTPS server. We
	//// pass in the paths to the TLS certificate and corresponding private key as
	//// the two parameters.
	//err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")

	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
