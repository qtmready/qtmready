create table if not exists slack (
    id                                  uuid,
    channel_id                          varchar,
    channel_name                        varchar,
    workspace_name                      varchar,
    team_id                             varchar,
    repo_provider                       varchar,
    repo_provider_installation_id       bigint,
    created_at                          timestamp,
    updated_at                          timestamp,
    primary key (id),
) with comment = 'github events state storage';

