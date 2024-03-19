alter table github_events_state add repo_name varchar;
create index if not exists on github_events_state (repo_name);
