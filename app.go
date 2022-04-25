package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	endpoints "github.com/jameshiii/mockbuster/endpoints"
	_ "github.com/lib/pq"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

var filmEndpoint = endpoints.Film{}
var filmCommentEndpoint = endpoints.FilmComment{}

// "Initialize" and "Run" methods are attached to the App struct so that function calls (ex: a.Run()) can happen elsewhere
func (a *App) Initialize(host, port, user, password, dbname string) {
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	var err error
	a.DB, err = sql.Open("postgres", connectionString)

	if err != nil {
		log.Fatal(err)
	}

	endpoints.New(a.DB)
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/films", filmEndpoint.GetList).Methods("GET")
	a.Router.HandleFunc("/film/{id:[0-9]+}", filmEndpoint.Get).Methods("GET")
	a.Router.HandleFunc("/film/{id:[0-9]+}/comments", filmCommentEndpoint.GetList).Methods("GET")
	a.Router.HandleFunc("/film/{id:[0-9]+}/comment", filmCommentEndpoint.Create).Methods("POST")
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}
