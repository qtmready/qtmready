
create table if not exists github_org_members (
  id                        uuid,
  name                      varchar,
  email                     varchar,
  company                   varchar,
  created_at                timestamp,
  updated_at                timestamp,
  primary key               (id),
) with comment = 'github_org_members defines a groups the github orgainzation members';
