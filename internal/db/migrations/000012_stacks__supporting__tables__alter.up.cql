alter table blueprints drop repo_id;
alter table blueprints drop repo_branch;
create index if not exists on blueprints (stack_id);

create index if not exists on resources (stack_id);
create index if not exists on workloads (stack_id);