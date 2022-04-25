package models

import (
	"database/sql"
	"errors"
	"strings"
)

type FilmModel struct {
}

type FilmListItem struct {
	Id       int     `json:"film_id"`
	Title    string  `json:"title"`
	Rating   string  `json:"rating"`
	Category *string `json:"category"`
}

type FilmDetail struct {
	Id              int     `json:"film_id"`
	Title           string  `json:"title"`
	Rating          string  `json:"rating"`
	Category        *string `json:"category"`
	Description     string  `json:"description"`
	ReleaseYear     int     `json:"release_year"`
	Actors          string  `json:"actors"`
	Language        string  `json:"language"`
	RentalDuration  int16   `json:"rental_duration"`
	RentalRate      float32 `json:"rental_rate"`
	Length          int16   `json:"length"`
	ReplacementCost float32 `json:"replacement_cost"`
}

func (f *FilmModel) GetList(db *sql.DB, title, rating, category string) ([]FilmListItem, error) {
	// results should most likely be paged using LIMIT OFFSET in production environment
	// 1 to zero rel to categories: a film may not have a category but if it does it will only have a single category

	title = strings.ToLower(title)
	parsedRating, err := parseStringtoRating(rating)
	category = strings.ToLower(category)

	if err != nil {
		return nil, err
	}

	//TODO: research adding multiple query parameters
	//pq.Array does not work with custom types like "mpaa_rating"
	rows, err := db.Query(`SELECT 
		film.film_id, 
		film.title, 
		film.rating,
		category.name as category
	FROM film
		LEFT OUTER JOIN film_category ON film.film_id = film_category.film_id
		LEFT OUTER JOIN category ON category.category_id = film_category.category_id
	WHERE 
		(($1 = '') IS NOT FALSE OR LOWER(title) LIKE '%' || $1 || '%') AND
		(($2 = '') IS NOT FALSE OR rating = $2::mpaa_rating ) AND
		(($3 = '') IS NOT FALSE OR LOWER(category.name) LIKE '%' || $3 || '%')`,
		title, parsedRating, category)

	if err != nil {
		return nil, err
	}

	defer rows.Close() // will prevent closing until func is complete

	films := []FilmListItem{}

	for rows.Next() {
		var f FilmListItem

		// copy row cols into pointers for "f"
		if err := rows.Scan(&f.Id, &f.Title, &f.Rating, &f.Category); err != nil {
			return nil, err
		}

		// add f to films slice
		films = append(films, f)
	}

	return films, nil
}

func (f *FilmModel) Get(db *sql.DB, id int) (*FilmDetail, error) {
	//TODO: get req's on which fields should be included in the "film details"

	film := FilmDetail{}

	err := db.QueryRow(`SELECT 
	film.film_id,
	film.title,
	film.rating,
	category.name as category,
	film.description,
	film.release_year,
	string_agg(actor.first_name || ' ' || actor.last_name, ', ') as actors,
	TRIM(language.name) as language,
	film.rental_duration,
	film.rental_rate,
	film.length,
	film.replacement_cost
	FROM film 
		LEFT OUTER JOIN film_category ON film.film_id = film_category.film_id
		LEFT OUTER JOIN category ON category.category_id = film_category.category_id
		LEFT OUTER JOIN film_actor ON film.film_id = film_actor.film_id
		LEFT OUTER JOIN actor ON actor.actor_id = film_actor.actor_id
		INNER JOIN language ON language.language_id = film.language_id
	WHERE film.film_id=$1
	GROUP By film.film_id, category.name, language.name`, id).Scan(
		&film.Id, &film.Title, &film.Rating, &film.Category, &film.Description, &film.ReleaseYear, &film.Actors, &film.Language, &film.RentalDuration, &film.RentalRate, &film.Length, &film.ReplacementCost)

	if err != nil {
		return nil, err
	}

	return &film, err
}

func parseStringtoRating(rating string) (string, error) {
	rating = strings.ToUpper(rating)

	// we only want to return an error if an invalid value is provided an empty string indicates that no value was provided
	switch rating {
	case "":
		return "", nil
	case "G":
		return "G", nil
	case "PG":
		return "PG", nil
	case "PG-13":
		return "PG-13", nil
	case "R":
		return "R", nil
	case "NC-17":
		return "NC-17", nil
	default:
		return "", errors.New("invalid rating provided")
	}
}
