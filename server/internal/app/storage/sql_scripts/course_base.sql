CREATE TABLE course_base (
    id serial PRIMARY KEY,
    title text,
    period_id integer,
    full_description text,
    description_timestamp timestamp,
    short_description text,
    FOREIGN KEY(period_id) REFERENCES period_base(id)
);
