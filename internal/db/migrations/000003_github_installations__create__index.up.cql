 

create table if not exists github_installations
(
  id                 uuid,
  team_id            uuid,
  installation_id    bigint,
  installation_type  varchar,
  installation_login varchar,
  sender_id          bigint,
  sender_login       varchar,
  status             varchar,
  created_at         timestamp,
  updated_at         timestamp,
  primary key (id),
) with comment = 'github app installation association with the team';

create index if not exists on github_installations (team_id);
create index if not exists on github_installations (installation_id);
