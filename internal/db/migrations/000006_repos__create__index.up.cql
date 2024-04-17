create table if not exists repos
(
  id                        uuid,
  team_id                   uuid,
  stack_id                  uuid,
  default_branch            varchar,
  is_monorepo               boolean,
  provider                  varchar,
  provider_id               varchar,
  name                      varchar,
  created_at                timestamp,
  updated_at                timestamp,
  primary key               (id),
) with comment = 'repo associated with the app';

create index if not exists on repos (stack_id);
create index if not exists on repos (provider);
create index if not exists on repos (provider_id);

