create table if not exists teams (
    id                          uuid,
    name                        varchar,
    slug                        varchar,
    created_at                  timestamp,
    updated_at                  timestamp,
    primary key                 (id),
);