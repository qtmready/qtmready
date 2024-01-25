create table if not exists github_action_state
(
  id                uuid,
  pull_request_id   bigint,
  repo_owner        varchar,
  repo_name         varchar,
  branch            varchar,
  created_at        timestamp,
  updated_at        timestamp,
  primary key       (id),
); with comment = 'maintain state of completed github actions in PR';

create index if not exists on pr_action_state (pull_request_id);
create index if not exists on pr_action_state (branch);
create index if not exists on pr_action_state (repo_owner);
create index if not exists on pr_action_state (repo_name);