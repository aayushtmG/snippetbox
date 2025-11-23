package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"github.com/aayushtmG/snippetbox/internal/models"
	_ "github.com/go-sql-driver/mysql"
)


type application struct {
	errLog *log.Logger
	infoLog *log.Logger
	snippets *models.SnippetModel
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

	
	app := &application{
		infoLog: infoLog,
		errLog: errLog,
		snippets: &models.SnippetModel{DB:db },
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
