CREATE TABLE IF NOT EXISTS users
(
    id              UUID primary key,
    first_name      varchar,
    last_name       varchar,
    email           varchar,
    created_at      timestamp,
    updated_at      timestamp,
    organization_id uuid
);