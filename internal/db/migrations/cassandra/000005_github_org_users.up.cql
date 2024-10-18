-- https://stackoverflow.com/questions/22139747/alter-composite-primary-key-in-cassandra-cql-3-0
DROP TABLE github_org_users;

create table if not exists github_org_users
(
  id                            uuid,
  user_id                       uuid,
  github_org_id                 bigint,
  github_user_id                bigint,
  github_org_name               varchar,
  created_at                    timestamp,
  updated_at                    timestamp,
  primary key                   (id),
) with comment = 'github org users, mapped to users';

