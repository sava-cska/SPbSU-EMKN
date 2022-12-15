CREATE TABLE checked_homework_base (
    homework_id serial,
    user_id serial,
    user_score integer,
    FOREIGN KEY(homework_id) REFERENCES homework_base(id),
    FOREIGN KEY(user_id) REFERENCES user_base(profile_id),
    PRIMARY KEY(homework_id, user_id)
);
