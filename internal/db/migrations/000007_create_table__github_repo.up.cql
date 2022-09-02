create table if not exists github_repos (
  id                uuid,
  team_id           uuid,
  github_id         bigint,
  name              varchar,
  full_name         varchar,
  created_at        timestamp,
  updated_at        timestamp,
  primary key       (id),
);
