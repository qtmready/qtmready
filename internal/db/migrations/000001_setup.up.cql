-- auth

create table if not exists teams
(
  id                uuid,
  name              varchar,
  slug              varchar,
  created_at        timestamp,
  updated_at        timestamp,
  primary key       (id),
) with comment = 'team, primary unit to define RBAC';

create index if not exists on teams (slug);

create table if not exists users
(
  id                    uuid,
  team_id               uuid,
  email                 varchar,
  password              varchar,
  first_name            varchar,
  last_name             varchar,
  is_verified           boolean,
  is_active             boolean,
  created_at            timestamp,
  updated_at            timestamp,
  primary key           (id),
) with comment = 'user, member of the team';

create index if not exists on users (team_id);
create index if not exists on users (email);

create table if not exists team_users (
  id                    uuid,
  team_id               uuid,
  user_id               uuid,
  role                  varchar,
  is_admin              boolean,
  is_active             boolean,
  created_at            timestamp,
  updated_at            timestamp,
  primary key           (id),
) with comment = 'team user association';

create index if not exists on team_users (team_id);
create index if not exists on team_users (user_id);

create table if not exists guards
(
  id                        uuid,
  name                      varchar,
  hashed                    varchar,
  lookup_id                 uuid,
  lookup_type               varchar,
  created_at                timestamp,
  updated_at                timestamp,
  primary key               (id),
) with comment = 'encrypted api keys for external services';

create index if not exists on guards (lookup_id);

create table if not exists oauth_accounts
(
  id                        uuid,
  user_id                   uuid,
  type                      varchar,
  provider                  varchar,
  provider_account_id       varchar,
  expires_at                timestamp,
  created_at                timestamp,
  updated_at                timestamp,
  primary key               (id, user_id),
) with comment = 'oauth account links against users';

create index if not exists on oauth_accounts (provider);
create index if not exists on oauth_accounts (provider_account_id);

-- core

create table if not exists stacks (
  id                        uuid,
  team_id                   uuid,
  name                      varchar,
  slug                      varchar,
  config                    blob,
  created_at                timestamp,
  updated_at                timestamp,
  primary key               (id, team_id),
) with comment = 'stack defines a group of services that are deployed together.';

create index if not exists on stacks (slug);

create table if not exists repos
(
  id                                  uuid,
  team_id                             uuid,
  stack_id                            uuid,
  default_branch                      varchar,
  is_monorepo                         boolean,
  provider                            varchar,
  provider_id                         varchar,
  ctrl_id                             uuid,
  name                                varchar,
  message_provider                    varchar,
  message_provider_id                 varchar,
  message_provider_data               blob,
  threshold                           int,
  created_at                          timestamp,
  updated_at                          timestamp,
  primary key                         (id, team_id),
) with comment = 'repo associated with the app';

create index if not exists on repos (stack_id);
create index if not exists on repos (provider);
create index if not exists on repos (provider_id);
create index if not exists on repos (ctrl_id);

create table if not exists workloads (
  id                        uuid,
  team_id                   uuid,
	stack_id                  uuid,
	repo_id                   uuid,
	resource_id               uuid,
	repo_path                 varchar,
	name                      varchar,
	kind                      varchar,
	container                 varchar,
	builder                   varchar,
	created_at                timestamp,
	updated_at                timestamp,
  primary key               (id, team_id),
) with comment = 'workload is the container definition';

create index if not exists on workloads (stack_id);
create index if not exists on workloads (repo_id);
create index if not exists on workloads (resource_id);

create table if not exists resources (
  id                         uuid,
  team_id                    uuid,
	stack_id                   uuid,
	name                       varchar,
	provider                   varchar,
	driver                     varchar,
  config                     varchar,
	is_immutable               boolean,
	created_at                 timestamp,
	updated_at                 timestamp,
  primary key                (id, team_id),
) with comment = 'resource defines the resource on cloud e.g. sqs, s3 etc';

create index if not exists on resources (stack_id);

create table if not exists blueprints (
  id                         uuid,
  team_id                    uuid,
	stack_id                   uuid,
	name                       varchar,
	regions                    varchar,
	provider_config            varchar,
	rollout_budget             varchar,
	created_at                 timestamp,
	updated_at                 timestamp,
  primary key                (id, team_id),
) with comment = 'blueprint defines the configuration for the progressive rollout';

create index if not exists on blueprints (stack_id);

create table if not exists changesets (
  id                  uuid,
  team_id             uuid,
  stack_id            uuid,
  repo_markers        blob,
  created_by          varchar,
  created_at          timestamp,
  updated_at          timestamp,
  primary key         (id, team_id),
) with comment = 'changeset defines the state of the stack defined by calver';

create index if not exists on changesets (stack_id);

create table if not exists rollouts (
  id                  uuid,
  team_id             uuid,
	stack_id            uuid,
	blueprint_id        uuid,
	state               varchar,
	created_at          timestamp,
	updated_at          timestamp,
  primary key         (id, team_id)
) with comment = 'rollouts is the result of a blueprint being applied';

create index if not exists on rollouts (team_id);
create index if not exists on rollouts (stack_id);

-- providers/github

create table if not exists github_org_users
(
  id                            uuid,
  user_id                       uuid,
  github_org_id                 bigint,
  github_user_id                bigint,
  github_org_name               varchar,
  created_at                    timestamp,
  updated_at                    timestamp,
  primary key                   (id, user_id)
) with comment = 'github org users, mapped to users';

create index if not exists on github_org_users (github_org_id);
create index if not exists on github_org_users (github_user_id);

create table if not exists github_installations
(
  id                            uuid,
  team_id                       uuid,
  installation_id               bigint,
  installation_type             varchar,
  installation_login            varchar,
  sender_id                     bigint,
  sender_login                  varchar,
  status                        varchar,
  created_at                    timestamp,
  updated_at                    timestamp,
  primary key (id, team_id),
) with comment = 'github app installation association with the team';

create index if not exists on github_installations (installation_id);

create table if not exists github_repos (
  id                                uuid,
  team_id                           uuid,
  github_id                         bigint,
  name                              varchar,
  full_name                         varchar,
  installation_id                   bigint,
  default_branch                    varchar,
  is_active                         boolean,
  has_early_warning                 boolean,
  created_at                        timestamp,
  updated_at                        timestamp,
  primary key                       (id, team_id)
) with comment = 'repositories provided by github installation';

create index if not exists on github_repos (github_id);
create index if not exists on github_repos (installation_id);
create index if not exists on github_repos (name);
create index if not exists on github_repos (full_name);
