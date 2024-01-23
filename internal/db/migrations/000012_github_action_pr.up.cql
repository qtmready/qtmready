create table if not exists github_action_pr
(
  id                uuid,
  repo_id           bigint,
  pr_id             bigint,
  total_num_action  bigint,
  curr_num_action   bigint,
  primary key       (id),
); with comment = 'maintain state of completed github actions in PR';

create index if not exists on pr_action_state (repo_id);
create index if not exists on pr_action_state (pr_id);