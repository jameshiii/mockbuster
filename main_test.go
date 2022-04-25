package main_test

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	main "github.com/jameshiii/mockbuster"
	"github.com/joho/godotenv"
)

var a main.App

//Responsible for loading settings for sub unit tests
func TestMain(m *testing.M) {

	err := godotenv.Load() // loads .env file

	if err != nil {
		log.Fatal(".env file failed to load")
	}

	// init sql connection and endpoint listener
	a.Initialize(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	code := m.Run()
	os.Exit(code)
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestTableHasRecords(t *testing.T) {
	req, _ := http.NewRequest("GET", "/films", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	bodyBytes, err := io.ReadAll(response.Body)

	if err != nil {
		t.Error(err)
	} else {
		bodyString := string(bodyBytes)
		if bodyString == "" || bodyString == "[]" {
			t.Errorf("No Records were returned in the Response.")
		}
	}
}

func TestFilmNotFound(t *testing.T) {
	req, _ := http.NewRequest("GET", "/film/-101181", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func TestGetFilmById(t *testing.T) {
	req, _ := http.NewRequest("GET", "/film/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	bodyBytes, err := io.ReadAll(response.Body)

	if err != nil {
		t.Error(err)
	} else {
		bodyString := string(bodyBytes)
		if bodyString == "" || bodyString == "[]" {
			t.Errorf("No Records were returned in the Response.")
		}
	}
}

func TestGetCommentsByFilmId(t *testing.T) {
	req, _ := http.NewRequest("GET", "/film/1/comments", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	bodyBytes, err := io.ReadAll(response.Body)

	if err != nil {
		t.Error(err)
	} else {
		bodyString := string(bodyBytes)
		if bodyString == "" || bodyString == "[]" {
			t.Errorf("No Records were returned in the Response.")
		}
	}
}

func TestCreateCommentByFilmId(t *testing.T) {

	var lastVal int
	_ = a.DB.QueryRow(`SELECT COUNT(*) - 1 FROM film_comment`).Scan(&lastVal)

	var jsonStr = []byte(`{"customerId": 1, "text": "Unit Test Value"}`)
	req, _ := http.NewRequest("POST", "/film/1/comment", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	bodyBytes, err := io.ReadAll(response.Body)

	if err != nil {
		t.Error(err)
	}

	bodyString := string(bodyBytes)
	_, err = strconv.Atoi(bodyString)

	if err != nil {
		t.Error("invalid value received for new comment identifier")
	}
}
