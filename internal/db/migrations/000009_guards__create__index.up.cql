-- Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved. 

create table if not exists guards
(
  id            uuid,
  name          varchar,
  hash          varchar,
  lookup_id     uuid,
  lookup_type   varchar,
  created_at    timestamp,
  updated_at    timestamp,
  primary key   (id),
) with comment = 'guard stores API for teams and users';

create index if not exists on guards (lookup_id);
