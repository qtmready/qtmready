-- integrations/github::github_installations::create
create table github_installations (
  id uuid primary key default uuid_generate_v7(),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  org_id uuid not null references orgs (id),
  installation_id bigint not null,
  installation_login varchar(255) not null,
  installation_login_id bigint not null,
  installation_type varchar(255) not null,
  sender_id bigint not null,
  sender_login varchar(255) not null,
  is_active boolean not null default true,
  constraint github_installations_installation_id_unique unique (installation_id)
);

-- integrations/github::github_installations::index
create unique index github_installations_installation_id_idx on github_installations (installation_id);

-- integrations/github::github_installations::trigger
create trigger update_github_installations_updated_at
  after update on github_installations
  for each row
  execute function update_updated_at();

-- integrations/github::github_orgs::create
create table github_orgs (
  id uuid primary key default uuid_generate_v7(),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  org_id uuid not null references orgs (id),
  installation_id uuid not null references github_installations (id) not null,
  github_org_id bigint not null,
  name varchar(255) not null
);

-- integrations/github::github_orgs::index
create index github_orgs_installation_id_idx on github_orgs (installation_id);

-- integrations/github::github_orgs::trigger
create trigger update_github_orgs_updated_at
  after update on github_orgs
  for each row
  execute function update_updated_at();

-- integrations/github::github_users::create
create table github_users (
  id uuid primary key default uuid_generate_v7(),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  user_id uuid references users (id) not null,
  github_id bigint not null,
  github_org_id uuid not null references github_orgs (id),
  login varchar(255) not null
);

-- integrations/github::github_users::trigger
create trigger update_github_users_updated_at
  after update on github_users
  for each row
  execute function update_updated_at();

-- integrations/github::github_repos::create
create table github_repos (
  id uuid primary key default uuid_generate_v7(),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  installation_id uuid not null references github_installations (id),
  github_id bigint not null,
  name varchar(255) not null,
  full_name varchar(255) not null,
  url varchar(255) not null,
  is_active boolean not null default true
);

-- integrations/github::github_repos::index
create index github_repos_installation_id_idx on github_repos (installation_id);

-- integrations/github::github_repos::trigger
create trigger update_github_repos_updated_at
  after update on github_repos
  for each row
  execute function update_updated_at();
