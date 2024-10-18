-- custom functions

create or replace function uuid_generate_v7() returns uuid as $$ begin -- use random v4 uuid as starting point (which has the same variant we need)
    -- then overlay timestamp
    -- then set version 7 by flipping the 2 and 1 bit in the version 4 string
    return encode(
        set_bit(
            set_bit(
                overlay(
                    uuid_send(gen_random_uuid()) placing substring(
                        int8send(
                            floor(
                                extract(
                                    epoch
                                    from clock_timestamp()
                                ) * 1000
                            )::bigint
                        )
                        from 3
                    )
                    from 1 for 6
                ),
                52,
                1
            ),
            53,
            1
        ),
        'hex'
    )::uuid;
end $$ language plpgsql volatile;

-- auth

create table users (
    id uuid primary key default uuid_generate_v7(),
    created_at timestamp without time zone not null,
    updated_at timestamp without time zone not null,
    org_id uuid not null references orgs (id),
    email text not null,
    first_name text,
    last_name text,
    password text,
    is_active boolean not null default true,
    is_verified boolean not null default false
);
create table teams (
    id uuid primary key default uuid_generate_v7(),
    created_at timestamp without time zone not null,
    updated_at timestamp without time zone not null,
    org_id uuid not null references orgs (id),
    name text not null,
    slug text not null
);
create table oauth_accounts (
    id uuid primary key default uuid_generate_v7(),
    created_at timestamp without time zone not null,
    updated_at timestamp without time zone not null,
    user_id uuid not null references users (id),
    provider text not null,
    provider_account_id text not null,
    expires_at timestamp without time zone,
    type text
);
create table orgs (
    id uuid primary key default uuid_generate_v7(),
    created_at timestamp without time zone not null,
    updated_at timestamp without time zone not null,
    org_id uuid not null references orgs (id),
    name text not null,
    slug text not null
);
create table team_users (
    id uuid primary key default uuid_generate_v7(),
    created_at timestamp without time zone not null,
    updated_at timestamp without time zone not null,
    team_id uuid not null references teams (id),
    user_id uuid not null references users (id),
    role text,
    is_active boolean not null default true,
    is_admin boolean not null default false
);

-- core

create table repos (
    id uuid primary key default uuid_generate_v7(),
    created_at timestamp without time zone not null,
    updated_at timestamp without time zone not null,
    org_id uuid not null references orgs (id),
    name text not null,
    provider text not null,
    provider_id text not null,
    default_branch text,
    is_monorepo boolean,
    threshold integer,
    stale_duration interval
);

-- integrations::github

create table github_installations (
    id uuid primary key default uuid_generate_v7(),
    created_at timestamp without time zone not null,
    updated_at timestamp without time zone not null,
    org_id uuid not null references orgs (id),
    installation_id bigint not null,
    installation_login text not null,
    installation_login_id bigint not null,
    installation_type text,
    sender_id bigint not null,
    sender_login text not null,
    status text
);

create unique index github_installations_installation_id_idx on github_installations (installation_id);

create table github_orgs (
    id uuid primary key default uuid_generate_v7(),
    created_at timestamp without time zone not null,
    updated_at timestamp without time zone not null,
    installation_id uuid not null references github_installations (id),
    github_org_id bigint not null,
    name text not null
);

create index github_orgs_installation_id_idx on github_orgs (installation_id);

create table github_users (
    id uuid primary key default uuid_generate_v7(),
    created_at timestamp without time zone not null,
    updated_at timestamp without time zone not null,
    user_id uuid references users (id),
    github_id bigint not null,
    github_org_id uuid not null references github_orgs (id),
    login text not null
);

create table github_repos (
    id uuid primary key default uuid_generate_v7(),
    created_at timestamp without time zone not null,
    updated_at timestamp without time zone not null,
    repo_id uuid not null references repos (id),
    installation_id uuid not null references github_installations (id),
    github_id bigint not null,
    name text not null,
    full_name text not null,
    url text not null,
    is_active boolean
);

create index github_repos_installation_id_idx on github_repos (installation_id);

-- messaging

create table messaging (
    id uuid primary key default uuid_generate_v7(),
    created_at timestamp without time zone not null,
    updated_at timestamp without time zone not null,
    provider text not null,
    kind text not null,
    link_to uuid not null,
    data jsonb not null
);

-- events

create type event_provider as enum (
    'github',
    'slack'
);

create table flat_events (
    id UUID primary key default uuid_generate_v7(),
    version text not null,
    parent_id UUID,
    provider event_provider not null,
    scope text not null,
    action text not null,
    source text not null,
    subject_id UUID not null,
    subject_name text not null,
    payload JSONB,
    team_id uuid not null references teams (id),
    user_id uuid not null references users (id),
    created_at timestamp without time zone not null,
    updated_at timestamp without time zone not null
);
