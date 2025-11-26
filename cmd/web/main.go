package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aayushtmG/snippetbox/internal/models"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
)


type application struct {
	errLog *log.Logger
	infoLog *log.Logger
	snippets *models.SnippetModel
	templateCache map[string]*template.Template
	formDecoder *form.Decoder
	sessionManager *scs.SessionManager
}


func main() {
	//getting values throud cmd
	addr := flag.String("addr",":4000","HTTP network address")
	dsn := flag.String("dsn","root:secret123@/snippetbox?parseTime=true","MySQL data source name")

	flag.Parse()

	infoLog := log.New(os.Stdout,"INFO\t",log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr,"ERROR\t",log.Ldate|log.Ltime|log.Lshortfile)

	db,err := openDB(*dsn)
	if err!= nil {
		errLog.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTemplateCache()	
	if err != nil {
		errLog.Fatal(err)
	}


	formDecoder := form.NewDecoder()


	//-session_manager
	sessionManager := scs.New()	
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	

	app := &application{
		infoLog: infoLog,
		errLog: errLog,
		snippets: &models.SnippetModel{DB:db },
		templateCache: templateCache,
		formDecoder: formDecoder,
		sessionManager: sessionManager,
	}

	srv := &http.Server{
		Addr: *addr,
		ErrorLog: errLog,
		Handler: app.routes(),
	}

	infoLog.Printf("Starting server on",*addr)
	err = srv.ListenAndServe()
	errLog.Fatal(err)
}


func openDB(dsn string) (*sql.DB,error){
	db, err := sql.Open("mysql",dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err!= nil {
		return nil, err
	}
	return db, nil
}
