create table if not exists repos
(
  id                        uuid,
  stack_id                  uuid,
  repo_id                   uuid,
  default_branch            varchar,
  is_monorepo               boolean,
  provider                  varchar,
  created_at                timestamp,
  updated_at                timestamp,
  primary key               (id),
) with comment = 'repo associated with the app';

create index if not exists on repos (stack_id);
create index if not exists on repos (repo_id);
create index if not exists on repos (provider);
