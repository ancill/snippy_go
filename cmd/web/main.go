package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

// app-wide deps struct
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	// CLI flag
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	infoLog.Printf("Starting server on %s", *addr)
	err := http.ListenAndServe(*addr, mux)
	errorLog.Fatal(err)
}
