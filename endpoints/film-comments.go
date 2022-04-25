package endpoints

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	models "github.com/jameshiii/mockbuster/models"
)

var commentModel = models.FilmCommentModel{}

type FilmComment struct{}

type filmCommentPostRequest struct {
	CustomerId int    `json:"customerId"`
	Text       string `json:"text"`
}

func (c *FilmComment) GetList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filmId, err := strconv.Atoi(vars["id"])

	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid film ID")
		return
	}

	films, err := commentModel.GetList(sqlDb, filmId)

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, films)
}

func (c *FilmComment) Create(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filmId, err := strconv.Atoi(vars["id"])

	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid film ID")
		return
	}

	var postArgs filmCommentPostRequest
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&postArgs); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	defer r.Body.Close()

	newId, err := commentModel.Create(sqlDb, postArgs.Text, filmId, postArgs.CustomerId)

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusCreated, newId)
}
