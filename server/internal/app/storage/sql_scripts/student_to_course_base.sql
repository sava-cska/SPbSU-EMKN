CREATE TABLE student_to_course_base (
    profile_id serial,
    course_id serial,
    FOREIGN KEY(profile_id) REFERENCES user_base(profile_id),
    FOREIGN KEY(course_id) REFERENCES course_base(id),
    PRIMARY KEY (profile_id, course_id)
);
