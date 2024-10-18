-- Drop the materialized view that depends on the stack_id column
drop index if exists repos_stack_id_idx;

-- alter the table to drop the column stack_id
alter table repos drop stack_id;

-- drop the tables if they exist
drop table if exists workloads;
drop table if exists resources;
drop table if exists blueprints;
drop table if exists changesets;
drop table if exists rollouts;
drop table if exists stacks;


