CREATE TABLE teacher_to_course_base (
    profile_id serial PRIMARY KEY,
    course_id serial PRIMARY KEY,
    FOREIGN KEY(profile_id) REFERENCES user_base(profile_id),
    FOREIGN KEY(course_id) REFERENCES course_base(id)
);
