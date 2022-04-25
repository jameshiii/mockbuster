CREATE TABLE film_comment (
	film_comment_id int PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	film_id int,
	customer_id int,
	comment_text varchar(500),
	create_date timestamp with time zone DEFAULT now(),
	FOREIGN KEY (film_id) REFERENCES film (film_id),
	FOREIGN KEY (customer_id) REFERENCES customer (customer_id)
)