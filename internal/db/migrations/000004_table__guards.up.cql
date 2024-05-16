create table if not exists guards
(
  id                        uuid,
  name                      varchar,
  hashed                    varchar,
  lookup_id                 uuid,
  lookup_type               varchar,
  created_at                timestamp,
  updated_at                timestamp,
  primary key               (id),
) with comment = 'encrypted api keys for external services';
