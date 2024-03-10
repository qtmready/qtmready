create table if not exists outh_accounts
(
  id                  uuid,
  user_id             uuid,
  type                varchar,
  provider            varchar,
  provider_account_id varchar,
  expires_at          timestamp,
  created_at          timestamp,
  updated_at          timestamp,
  primary key   (id),
) with comment = 'oauth account links against users';

create index if not exists on oauth_accounts (user_id);
create index if not exists on oauth_accounts (provider);
create index if not exists on oauth_accounts (provider_account_id);
