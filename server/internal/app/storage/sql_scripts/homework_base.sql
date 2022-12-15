CREATE TABLE homework_base (
    id serial PRIMARY KEY,
    name text,
    deadline timestamp,
    course_id serial,
    score integer,
    FOREIGN KEY(course_id) REFERENCES course_base(id)
);
