create table flat_events__v_0_1 (
  version             text,
  id                  uuid,
  parent_id           uuid,
  provider            text,
  scope_              text,
  action              text,
  source              text,
  subject_id          uuid,
  subject_name        text,
  payload             blob,
  team_id             uuid,
  user_id             uuid,
  created_at          timestamp,
  updated_at          timestamp,
  primary key ((subject_id, team_id), created_at)
) with clustering order by (created_at asc);

create index scope_index on flat_events__v_0_1 (scope_);
create index action_index on flat_events__v_0_1 (action);
create index source_index on flat_events__v_0_1 (source);
