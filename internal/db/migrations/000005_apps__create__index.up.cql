
create table if not exists apps (
  id                uuid,
  team_id           uuid,
  name              varchar,
  slug              varchar,
  created_at        timestamp,
  updated_at        timestamp,
  primary key       (id),
) with comment = 'the main application construct';

create index if not exists on apps (team_id);
create index if not exists on apps (slug);
