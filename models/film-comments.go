package models

import (
	"database/sql"
)

type FilmCommentModel struct {
}

type FilmCommentListItem struct {
	Id           int    `json:"film_comment_id"`
	FilmId       int    `json:"film__id"`
	CustomerId   int    `json:"customer_id"`
	CustomerName string `json:"customer_name"`
	Text         string `json:"comment_text"`
	CreatedDate  string `json:"create_date"`
}

func (c *FilmCommentModel) GetList(db *sql.DB, filmId int) ([]FilmCommentListItem, error) {

	rows, err := db.Query(`SELECT 
		film_comment_id, 
		film_id,
		customer.customer_id,
		customer.first_name || ' ' || LEFT(customer.last_name, 1) || '.' as customer_name,
		comment_text,
		film_comment.create_date
	FROM film_comment
		JOIN customer ON customer.customer_id = film_comment.customer_id
	WHERE film_id = $1`, filmId)

	if err != nil {
		return nil, err
	}

	defer rows.Close() // will prevent closing until func is complete

	comments := []FilmCommentListItem{}

	for rows.Next() {
		var c FilmCommentListItem

		// copy row cols into pointers
		if err := rows.Scan(&c.Id, &c.FilmId, &c.CustomerId, &c.CustomerName, &c.Text, &c.CreatedDate); err != nil {
			return nil, err
		}

		// add c to films slice
		comments = append(comments, c)
	}

	return comments, nil
}

func (c *FilmCommentModel) Create(db *sql.DB, text string, film_id, customer_id int) (*int, error) {

	var id int

	err := db.QueryRow(
		"INSERT INTO film_comment (film_id, customer_id, comment_text) VALUES($1, $2, $3) RETURNING film_comment_id",
		&film_id, &customer_id, &text).Scan(&id)

	if err != nil {
		return nil, err
	}

	return &id, nil
}
