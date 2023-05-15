create table if not exists stacks (
  id                uuid,
  team_id           uuid,
  name              varchar,
  slug              varchar,
  config            blob,
  created_at        timestamp,
  updated_at        timestamp,
  primary key       (id),
) with comment = 'stack defines a group of services that are deployed together.';

create index if not exists on stacks (team_id);
create index if not exists on stacks (slug);
