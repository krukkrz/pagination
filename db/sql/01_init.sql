create table if not exists books (
    book_id serial PRIMARY KEY,
    title VARCHAR ( 100 ) NOT NULL,
    author VARCHAR ( 100 ) NOT NULL,
    created_at TIMESTAMP NOT NULL
);

CREATE SEQUENCE books_sequence
    start 1
  increment 1;
