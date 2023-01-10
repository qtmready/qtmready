drop index repos_repo_id_idx;
alter table repos drop repo_id;
alter table repos add provider_id varchar;
create index if not exists on repos (provider_id);
