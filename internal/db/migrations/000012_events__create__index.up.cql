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

