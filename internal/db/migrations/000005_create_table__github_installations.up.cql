create table if not exists github_installations
(
    id                          uuid,
    team_id                     uuid,
    installation_id             bigint,
    installation_type           varchar,
    installation_login          varchar,
    sender_id                   bigint,
    sender_login                varchar,
    created_at                  timestamp,
    updated_at                  timestamp,
    primary key                 (id),
);