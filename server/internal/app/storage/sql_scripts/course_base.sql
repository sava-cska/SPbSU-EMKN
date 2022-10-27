CREATE TABLE course_base (
    id serial PRIMARY KEY,
    title text,
    period_id integer,
    short_description text,
    description_timestamp timestamp,
    FOREIGN KEY(period_id) REFERENCES period_base(id)
);
