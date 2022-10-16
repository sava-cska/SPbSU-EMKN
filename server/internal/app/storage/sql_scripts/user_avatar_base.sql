CREATE TABLE user_avatar_base (
    profile_id serial PRIMARY KEY,
    avatar_url varchar,
    first_name varchar,
    last_name varchar,
    FOREIGN KEY(profile_id) REFERENCES user_base(profile_id)
);
