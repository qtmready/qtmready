 

create table if not exists workloads (
  	id            uuid,
		app_id        uuid,
		repo_id       uuid,
		repo_path     varchar,
		name          varchar,
		kind          varchar,
		container     varchar,
		builder       varchar,
		created_at    timestamp,
		updated_at    timestamp,
    primary key   (id),
) with comment = 'workload defines a single unit of compute';

create table if not exists resources (
  	id            uuid,
		app_id        uuid,
		name          varchar,
		provider      varchar,
		driver        varchar,
		is_immutable  boolean,
		created_at    timestamp,
		updated_at    timestamp,
    primary key   (id),
) with comment = 'resource is cloud provider specific';


create table if not exists blueprints (
    id              uuid,
		app_id          uuid,
		repo_id         uuid,
		name            varchar,
		regions         varchar,
		repo_branch     varchar,
		rollout_budget  varchar,
		created_at      timestamp,
		updated_at      timestamp,
    primary key     (id),
) with comment = 'blueprint defines the deployment and how to rollout';

create table if not exists rollouts (
  	id                  uuid,
		app_id              uuid,
		blueprint_id        uuid,
		trigger             varchar,
		state               varchar,
		created_at          timestamp,
		updated_at          timestamp,
    primary key         (id),
) with comment = 'rollouts is the result of a blueprint being applied';
