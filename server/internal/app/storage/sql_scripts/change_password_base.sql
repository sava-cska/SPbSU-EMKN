CREATE TABLE change_password_base (token varchar PRIMARY KEY, login varchar, expire_date timestamp, varification_code varchar, FOREIGN KEY(login) REFERENCES user_base(login));
