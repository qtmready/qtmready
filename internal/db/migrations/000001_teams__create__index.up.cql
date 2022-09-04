create table if not exists teams (
    id                          uuid,
    name                        varchar,
    slug                        varchar,
    created_at                  timestamp,
    updated_at                  timestamp,
    primary key                 (id),
) with comment = 'team is the primary unit to define RBAC';

create index if not exists on teams (slug);
