CREATE TABLE registration_base (
    token  varchar PRIMARY KEY,
    login varchar,
    password varchar,
    email varchar,
    first_name varchar,
    last_name varchar,
    expire_date timestamp,
    verification_code varchar
);