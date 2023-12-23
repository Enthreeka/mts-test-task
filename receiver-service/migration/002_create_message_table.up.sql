CREATE TABLE IF NOT EXISTS message(
    id int generated always as identity,
    msg_uuid uuid not null,
    msg text unique,
    created_time timestamp,
    primary key (id)
);