create table if not exists app_repos (
  id                uuid,
  app_id            uuid,
  repo_id           uuid,
  default_branch    varchar,
  is_monorepo       boolean,
  git_provider      varchar,
  created_at        timestamp,
  updated_at        timestamp,
  primary key       (id),
) with comment = 'repo associated with an app';

create index if not exists app_repos(app_id);
create index if not exists app_repos(repo_id);
