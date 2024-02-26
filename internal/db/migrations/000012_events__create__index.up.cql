create table if not exists github_events_state (
    id                      uuid,
    event_type              varchar,
    status                  varchar,
    events_data             varchar,
    github_workflow_id      bigint,
    github_workflow_run_id  bigint,
    created_at              timestamp,
    updated_at              timestamp,
    primary key (id),
) with comment = 'github events state storage';

create index if not exists on github_events_state (id);
create index if not exists on github_events_state (event_type);
create index if not exists on github_events_state (status);
create index if not exists on github_events_state (github_workflow_id);
create index if not exists on github_events_state (github_workflow_run_id);
