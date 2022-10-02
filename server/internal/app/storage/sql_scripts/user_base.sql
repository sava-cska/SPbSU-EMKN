CREATE TABLE user_base (
    login varchar PRIMARY KEY,
    password varchar,
    email varchar UNIQUE,
    first_name varchar, 
    last_name varchar
);
