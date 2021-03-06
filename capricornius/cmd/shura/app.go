package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	host     = ""
	port     = ""
	user     = ""
	password = ""
	dbname   = ""
)

type App struct {
	Router *mux.Router
	DB     *sqlx.DB
}

func (a *App) Initialize(user, password, dbname string) {
	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s", user, password, dbname)

	var err error
	a.DB, err = sqlx.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(":8000", a.Router))
}

func (a *App) getPartner(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Partner ID")
		return
	}

	p := partner{ID: id}
	if err := p.getPartner(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Product not found")

		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, p)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(responses)
}

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/partner", getPartner).Methods("GET")
	router.HandleFunc("/partner/{id}", GetPartnersfromID).Methods("GET")
	router.HandleFunc("/partner/{name}", GetPartnersformName).Methods("GET")
	router.HandleFunc("/partner/{email}", GetPartnerfromEmail).Methods("GET")
	router.HandleFunc("/partner/{token}", GetPartnerfromToken).Methods("GET")
	router.HandleFunc("/partner", CreatePartner).Methods("POST")
	router.HandleFunc("/partner/{id}", UpdatePartner).Methods("POST")
	router.HandleFunc("/partner/{id}", DeletePartner).Methods("POST")

	log.Fatal(http.ListenAndServe(":5000", router))

}
