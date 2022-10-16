CREATE TABLE user_base (
    login varchar PRIMARY KEY,
    password varchar,
    email varchar UNIQUE,
    profile_id serial,
    first_name varchar, 
    last_name varchar
);
