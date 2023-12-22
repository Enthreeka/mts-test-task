CREATE TABLE message(
    id int generated always as identity,
    msg text,
    created_time timestamp,
    primary key (id)
);