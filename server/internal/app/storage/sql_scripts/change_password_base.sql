CREATE TABLE change_password_base (
    token varchar PRIMARY KEY,
    login varchar,
    expire_date timestamp, 
    verification_code varchar,
    change_password_token varchar,
    FOREIGN KEY(login) REFERENCES user_base(login)
);
