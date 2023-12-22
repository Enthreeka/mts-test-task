CREATE TABLE message(
    id int generated always as identity,
    msg text,
    create_time timestamp,
    primary key (id)
);