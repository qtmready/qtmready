create table if not exists users (
    id                          uuid,
    email                       varchar,
    password                    varchar,
    first_name                  varchar,
    last_name                   varchar,
    slug                        varchar,
    is_verified                 boolean,
    is_active                   boolean,
    created_at                  timestamp,
    updated_at                  timestamp,
    primary key                 (id),
);

create index if not exists on users (email);
create index if not exists on users (slug);

create table if not exists teams (
    id                          uuid,
    name                        varchar,
    slug                        varchar,
    created_at                  timestamp,
    updated_at                  timestamp,
    primary key                 (id),
);

create index if not exists on teams (slug);

create table if not exists github_installations
(
    id                          uuid,
    team_id                     uuid,
    github_installation_id      bigint,
    github_installation_type    varchar,
    github_sender_id            bigint,
    github_login                varchar, 
    github_sender_login         varchar,
    primary key                 (id),   
);

create index if not exists on github_installations (github_installation_id);