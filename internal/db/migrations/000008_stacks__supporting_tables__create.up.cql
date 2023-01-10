create table if not exists workloads (
  	id            uuid,
		stack_id      uuid,
		repo_id       uuid,
		repo_path     varchar,
		name          varchar,
		kind          varchar,
		container     varchar,
		builder       varchar,
		created_at    timestamp,
		updated_at    timestamp,
    primary key   (id),
) with comment = 'workload is the container definition';

create table if not exists resources (
  	id            uuid,
		stack_id      uuid,
		name          varchar,
		provider      varchar,
		driver        varchar,
		is_immutable  boolean,
		created_at    timestamp,
		updated_at    timestamp,
    primary key   (id),
) with comment = 'resource is cloud provider specific';


create table if not exists blueprints (
    id                  uuid,
		stack_id            uuid,
		repo_id             uuid,
		name                varchar,
		regions             varchar,
		repo_branch         varchar,
		rollout_budget      varchar,
		created_at          timestamp,
		updated_at          timestamp,
    primary key         (id),
) with comment = 'blueprint defines the deployment configuration';

create table if not exists changesets (
    id                  uuid,
    stack_id            uuid,
    calver              varchar,
    repo_markers        blob,
    created_by          varchar,
    created_at          timestamp,
    updated_at          timestamp,
    primary key         (id),
) with comment = 'changeset defines the state of the stack defined by calver';

create table if not exists rollouts (
  	id                  uuid,
		stack_id            uuid,
		blueprint_id        uuid,
		state               varchar,
		created_at          timestamp,
		updated_at          timestamp,
    primary key         (id),
) with comment = 'rollouts is the result of a blueprint being applied';
