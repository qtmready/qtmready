alter table github_repos add installation_id bigint;
create index if not exists on github_repos (installation_id);
