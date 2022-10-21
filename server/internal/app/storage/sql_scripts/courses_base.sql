CREATE TABLE courses_base (
    id serial PRIMARY KEY,
    title text,
    periods_id integer,
    short_description text,
    description_timestamp timestamp,
    FOREIGN KEY(periods_id) REFERENCES periods_base(id)
);