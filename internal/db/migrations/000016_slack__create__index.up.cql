create table if not exists slack (
    id                      uuid,
    channel_id              varchar,
    channel_name            varchar,
    workspace_name          varchar
    team_id                 varchar,
    created_at              timestamp,
    updated_at              timestamp,
    primary key (id),
) with comment = 'github events state storage';

