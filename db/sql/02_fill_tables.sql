INSERT INTO books (book_id, title, author, created_at)
SELECT nextval('books_sequence'),
       concat('Title ', i),
       concat('Author ', i),
       current_timestamp
FROM generate_series(1, 200) AS i;

INSERT INTO cars (car_id, brand, model, created_at)
SELECT nextval('cars_sequence'),
       concat('Brand ', i),
       concat('Model ', i),
       current_timestamp
FROM generate_series(1, 200) AS i;