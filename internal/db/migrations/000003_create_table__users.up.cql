create table if not exists users (
    id                          uuid,
    email                       varchar,
    password                    varchar,
    first_name                  varchar,
    last_name                   varchar,
    is_verified                 boolean,
    is_active                   boolean,
    created_at                  timestamp,
    updated_at                  timestamp,
    primary key                 (id),
);