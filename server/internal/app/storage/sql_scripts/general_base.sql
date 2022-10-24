CREATE TABLE general_base (
    id integer unique ,
    current_period_id integer,
    FOREIGN KEY(current_period_id) REFERENCES periods_base(id)
);