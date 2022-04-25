package endpoints

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jameshiii/mockbuster/models"
)

var filmModel = models.FilmModel{}

type Film struct{}

func (f *Film) GetList(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	rating := r.URL.Query().Get("rating")
	category := r.URL.Query().Get("category")

	films, err := filmModel.GetList(sqlDb, title, rating, category)

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, films)
}

func (f *Film) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid film ID")
		return
	}

	film, err := filmModel.Get(sqlDb, id)

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, film)
}
