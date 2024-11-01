-- auth::orgs::create
create table orgs (
  id uuid primary key default uuid_generate_v7(),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  name varchar(255) not null,
  domain varchar(255) not null,
  slug varchar(255) not null,
  hooks jsonb not null default '{"repo":0, "messaging": 0}',
  constraint orgs_domain_unique unique (domain),
  constraint orgs_slug_unique unique (slug)
);

-- auth::orgs::trigger
create trigger update_orgs_updated_at
  after update on orgs
  for each row
  execute function update_updated_at();

-- auth::orgs::index
create unique index orgs_domain_idx on orgs (domain);
create unique index orgs_slug_idx on orgs (slug);

-- auth::teams::create
create table teams (
  id uuid primary key default uuid_generate_v7(),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  org_id uuid not null references orgs (id),
  name varchar(255) not null,
  slug varchar(255) not null
);

-- auth::teams::trigger
create trigger update_teams_updated_at
  after update on teams
  for each row
  execute function update_updated_at();

-- auth::users::create
create table users (
  id uuid primary key default uuid_generate_v7(),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  org_id uuid not null references orgs (id),
  email varchar(255) not null,
  first_name varchar(255) not null,
  last_name varchar(255) not null,
  password text not null,
  picture text not null,
  is_active boolean not null default true,
  is_verified boolean not null default false,
  constraint users_email_unique unique (email)
);

-- auth::users::trigger
create trigger update_users_updated_at
  after update on users
  for each row
  execute function update_updated_at();

-- auth::users::index
create unique index users_email_idx on users (email);

-- auth::team_users::create
create type team_role as enum ('member', 'admin');

create table team_users (
  id uuid primary key default uuid_generate_v7(),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  team_id uuid not null references teams (id),
  user_id uuid not null references users (id),
  role team_role not null default 'member',
  is_active boolean not null default true,
  is_admin boolean not null default false
);

-- auth::team_users::trigger
create trigger update_team_users_updated_at
  after update on team_users
  for each row
  execute function update_updated_at();

-- auth::user_roles::create
create table user_roles (
  id uuid primary key default uuid_generate_v7(),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  name varchar(63) not null,
  user_id uuid not null references users (id),
  org_id uuid not null references orgs (id)
);

-- auth::user_roles::trigger
create trigger update_user_roles_updated_at
  after update on user_roles
  for each row
  execute function update_updated_at();

-- auth::oauth_accounts::create
create table oauth_accounts (
  id uuid primary key default uuid_generate_v7(),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  user_id uuid not null references users (id),
  provider varchar(255) not null,
  provider_account_id varchar(255) not null,
  expires_at timestamptz not null,
  type varchar(255) not null,
  constraint oauth_accounts_unique_provider_account unique (provider, provider_account_id)
);

-- auth::oauth_accounts::trigger
create trigger update_oauth_accounts_updated_at
  after update on oauth_accounts
  for each row
  execute function update_updated_at();
