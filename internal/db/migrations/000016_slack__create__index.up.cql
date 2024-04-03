create table if not exists slack (
    id                      uuid,
    channel_id              varchar,
    channel_name            varchar,
    workspace_name          varchar,
    workspace_id            varchar,
    workspace_bot_token     varchar,
    team_id                 uuid,
    created_at              timestamp,
    updated_at              timestamp,
    primary key (id),
) with comment = 'slack channel info storage';

