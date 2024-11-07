-- core::repos::create
create table repos (
  id uuid primary key default uuid_generate_v7(),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  org_id uuid not null references orgs (id),
  name varchar(255) not null,
  hook varchar(255) not null,
  hook_id uuid not null,
  default_branch varchar(255) not null default 'main',
  is_monorepo boolean not null default false,
  threshold integer not null default 500,
  stale_duration interval not null default '2 days',
  url varchar(255) not null,
  is_active boolean not null default true
);

-- core::repos::trigger
create trigger update_repos_updated_at
  after update on repos
  for each row
  execute function update_updated_at();

-- messaging::messaging::create
create table messaging (
  id uuid primary key default uuid_generate_v7(),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  hook varchar(255) not null,
  kind varchar(255) not null,
  link_to uuid not null,
  data jsonb not null default '{}'
);

-- messaging::messaging::trigger
create trigger update_messaging_updated_at
  after update on messaging
  for each row
  execute function update_updated_at();
