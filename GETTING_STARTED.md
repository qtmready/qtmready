# Quantm.io

Ship 15x faster. Rollback in minutes!

A modern platform designed to accelerate developer velocity, automate workflows, and ensure smooth, fast rollbacks. With [**Quantm.io**](https://quantm.io/), developers can deploy efficiently and focus on building value, knowing that every change is trackable and recoverable within minutes.

## Prerequisites

Ensure you have the following installed before you begin:

- [**Go** (v1.23+)](https://go.dev/doc/install)
- [**Temporal**](https://docs.temporal.io/docs) (for orchestrating workflows)
- [**Postgres**](https://www.postgresql.org/download/) (as the primary database)
- [**Docker**](https://docs.docker.com/get-docker/) (for containerization)
- [**taskfile**](https://taskfile.dev/#/) (for task automation)
- [**buf**](https://docs.buf.build/installation) (for managing protobuf files)

## Getting Started

1. Clone the repository:

   ```bash
   git clone https://github.com/quantmHQ/quantm.git
   cd quantm
   ```

2. Set up environment variables:
   Create a .env file based on the example:

   ```bash
   cp .env.example .env
   ```

3. Start the services using Docker Compose:

   ```bash
   docker-compose up -d
   ```

### Cmd

task: available tasks for this project:

```bash
sqlc: Generate Go code from SQL files
api:go: Generate Go code from Protobuf files
api:lint: Lint Protobuf files
api:ts: Generate TypeScript code from Protobuf files
```

## Folder Structure

Understanding the folder structure of quantm can help you navigate the project effectively. Here’s a brief overview:

```
├── CODE_OF_CONDUCT.md
├── Dockerfile
├── GETTING_STARTED.md
├── LICENSE
├── README.md -> ./docs/src/ref/overview.md
├── add-resource.py
├── api
│   ├── buf.gen.go.yaml
│   ├── buf.gen.ts.yaml
│   ├── buf.lock
│   ├── buf.yaml
│   ├── openapi
│   │   ├── auth
│   │   │   └── v1
│   │   │       └── schema.yaml
│   │   ├── core
│   │   │   └── v1
│   │   │       ├── components.yaml
│   │   │       └── paths.yaml
│   │   ├── github
│   │   │   └── v1
│   │   │       └── schema.yaml
│   │   ├── shared
│   │   │   └── v1
│   │   │       └── schema.yaml
│   │   └── slack
│   │       └── v1
│   │           └── schema.yaml
│   ├── proto
│   │   └── ctrlplane
│   │       ├── auth
│   │       │   └── v1
│   │       │       ├── accounts.proto
│   │       │       └── orgs.proto
│   │       └── common
│   │           └── v1
│   │               └── uuid.proto
│   └── taskfile.yml
├── clean.py
├── cmd
│   ├── api
│   │   ├── alpha
│   │   │   └── main.go
│   │   └── legacy
│   │       ├── auth.go
│   │       ├── echo.go
│   │       ├── handlers.go
│   │       ├── main.go
│   │       └── otel.go
│   ├── e2e
│   │   ├── mutex
│   │   │   └── main.go
│   │   ├── orm
│   │   │   └── main.go
│   │   └── qtm
│   │       └── main.go
│   ├── jobs
│   │   ├── migrate
│   │   │   └── main.go
│   │   └── openapi
│   │       ├── custom.go
│   │       ├── main.go
│   │       ├── orgs.go
│   │       └── shared.go
│   └── workers
│       ├── mothership
│       │   ├── main.go
│       │   └── queues.go
│       └── sentinel
│           └── main.go
├── deploy
│   ├── air
│   │   ├── api.toml
│   │   ├── migrate.toml
│   │   └── mothership.toml
│   ├── cassandra
│   │   ├── cassandra.yaml
│   │   └── keyspace.cql
│   ├── postgres
│   │   └── init.sql
│   └── temporal
│       └── dynamicconfig
│           ├── development-cass.yaml
│           ├── development.yaml
│           └── docker.yaml
├── do-config.py
├── docker-compose.yaml
├── docs
│   ├── book.toml
│   └── src
│       ├── SUMMARY.md
│       ├── contributing
│       │   ├── contributor.md
│       │   ├── creating-pr.md
│       │   ├── getting-started.md
│       │   ├── guidelines.md
│       │   ├── merge-responsibility.md
│       │   ├── project-management.md
│       │   ├── pull-requests.md
│       │   └── reviewer.md
│       ├── design
│       │   └── plantuml
│       │       └── create_resources_changeset.md
│       ├── gitOps_design.puml
│       ├── introduction.md
│       └── ref
│           ├── overview.md
│           └── raison-detre.md
├── go.mod
├── go.sum
├── internal
│   ├── auth
│   │   ├── crypto.go
│   │   ├── crypto_test.go
│   │   ├── doc.go
│   │   ├── entity.go
│   │   ├── errors.go
│   │   ├── guards.go
│   │   ├── guards_test.go
│   │   ├── handler.go
│   │   ├── handler_test.go
│   │   ├── io_team.go
│   │   ├── io_team_user.go
│   │   ├── io_user.go
│   │   ├── middleware.go
│   │   ├── openapi.codegen.yaml
│   │   ├── openapi.gen.go
│   │   ├── teams_test.go
│   │   └── users_test.go
│   ├── core
│   │   ├── README.md
│   │   ├── code
│   │   │   ├── activities.go
│   │   │   ├── base_state.go
│   │   │   ├── branch_ctrl.go
│   │   │   ├── branch_ctrl_state.go
│   │   │   ├── errors.go
│   │   │   ├── helpers.go
│   │   │   ├── io_repo.go
│   │   │   ├── queries.go
│   │   │   ├── queue_ctrl.go
│   │   │   ├── queue_ctrl_state.go
│   │   │   ├── repo_ctrl.go
│   │   │   ├── repo_ctrl_state.go
│   │   │   ├── trunk_ctrl.go
│   │   │   ├── trunk_ctrl_state.go
│   │   │   ├── utils.go
│   │   │   ├── workflow_logger.go
│   │   │   └── workflow_options.go
│   │   ├── comm
│   │   │   └── fns.go
│   │   ├── defs
│   │   │   ├── README.md
│   │   │   ├── alias.go
│   │   │   ├── defs.gen.go
│   │   │   ├── docs.go
│   │   │   ├── entity.go
│   │   │   ├── errors.go
│   │   │   ├── events.go
│   │   │   ├── events_test.go
│   │   │   ├── message_io.go
│   │   │   ├── oapi-codegen.yaml
│   │   │   └── repo_io.go
│   │   ├── events
│   │   │   ├── README.md
│   │   │   ├── actions.go
│   │   │   ├── alias.go
│   │   │   ├── doc.go
│   │   │   ├── errors.go
│   │   │   ├── events.go
│   │   │   ├── flat.go
│   │   │   ├── payloads.go
│   │   │   ├── scopes.go
│   │   │   ├── subjects.go
│   │   │   ├── validations.go
│   │   │   └── versions.go
│   │   ├── kernel
│   │   │   ├── kernel.go
│   │   │   ├── message.go
│   │   │   └── repo.go
│   │   ├── mutex
│   │   │   ├── activities.go
│   │   │   ├── doc.go
│   │   │   ├── errors.go
│   │   │   ├── log.go
│   │   │   ├── mutex.go
│   │   │   ├── pool.go
│   │   │   ├── state.go
│   │   │   ├── workflow.go
│   │   │   └── workflow_options.go
│   │   ├── periodic
│   │   │   └── interval.go
│   │   ├── timers
│   │   │   ├── interval.go
│   │   │   └── utils.go
│   │   ├── web
│   │   │   ├── docs.go
│   │   │   ├── handlers.go
│   │   │   ├── oapi-codegen.yaml
│   │   │   └── web.gen.go
│   │   └── ws
│   │       ├── activities.go
│   │       ├── auth.go
│   │       ├── doc.go
│   │       ├── errors.go
│   │       ├── hub.go
│   │       ├── ids.go
│   │       ├── ids_bench_test.go
│   │       ├── queue.go
│   │       ├── state.go
│   │       └── workflows.go
│   ├── db
│   │   ├── cassandra.go
│   │   ├── config
│   │   │   ├── config.go
│   │   │   └── queries.go
│   │   ├── db.go
│   │   ├── entities
│   │   │   ├── db.go
│   │   │   ├── github_installation.sql.go
│   │   │   ├── github_org.sql.go
│   │   │   ├── github_repo.sql.go
│   │   │   ├── github_user.sql.go
│   │   │   ├── models.go
│   │   │   ├── oauth_account.sql.go
│   │   │   ├── org.sql.go
│   │   │   ├── repo.sql.go
│   │   │   ├── team.sql.go
│   │   │   ├── team_user.sql.go
│   │   │   └── user.sql.go
│   │   ├── fields
│   │   │   ├── duration.go
│   │   │   ├── int64.go
│   │   │   ├── sensitive.go
│   │   │   └── sensitive_test.go
│   │   ├── fields.go
│   │   ├── mapper.go
│   │   ├── migrations
│   │   │   ├── cassandra
│   │   │   │   ├── 000001_setup.up.cql
│   │   │   │   ├── 000002_alter_github_installations.down.cql
│   │   │   │   ├── 000002_alter_github_installations.up.cql
│   │   │   │   ├── 000003_drop_unused_tables.down.cql
│   │   │   │   ├── 000003_drop_unused_tables.up.cql
│   │   │   │   ├── 000004_message_provider_linked_user.down.cql
│   │   │   │   ├── 000004_message_provider_linked_user.up.cql
│   │   │   │   ├── 000005_github_org_users.down.cql
│   │   │   │   ├── 000005_github_org_users.up.cql
│   │   │   │   ├── 000006_alter_core_repo.down.cql
│   │   │   │   ├── 000006_alter_core_repo.up.cql
│   │   │   │   ├── 000007_table__flat_event_0_1.down.cql
│   │   │   │   └── 000007_table__flat_event_0_1.up.cql
│   │   │   └── postgres
│   │   │       ├── 000001_setup.down.sql
│   │   │       ├── 000001_setup.up.sql
│   │   │       ├── 000002_create_tables.down.sql
│   │   │       └── 000002_create_tables.up.sql
│   │   ├── orm.go
│   │   ├── queries
│   │   │   ├── github_installation.sql
│   │   │   ├── github_org.sql
│   │   │   ├── github_repo.sql
│   │   │   ├── github_user.sql
│   │   │   ├── messaging.sql
│   │   │   ├── oauth_account.sql
│   │   │   ├── org.sql
│   │   │   ├── repo.sql
│   │   │   ├── team.sql
│   │   │   ├── team_user.sql
│   │   │   └── user.sql
│   │   ├── sequel
│   │   │   ├── org.wrap.go
│   │   │   ├── prompt.md
│   │   │   ├── team.wrap.go
│   │   │   └── user.wrap.go
│   │   ├── sqlc.yml
│   │   ├── utils.go
│   │   ├── uuid.go
│   │   └── validations.go
│   ├── durable
│   │   └── connection
│   │       └── config.go
│   ├── erratic
│   │   ├── details.go
│   │   ├── http.go
│   │   └── quantm.go
│   ├── nomad
│   │   ├── convert
│   │   │   ├── accounts.go
│   │   │   └── uuid.go
│   │   ├── handler
│   │   │   └── accounts.go
│   │   └── proto
│   │       ├── ctrlplane
│   │       │   ├── auth
│   │       │   │   └── v1
│   │       │   │       ├── accounts.pb.go
│   │       │   │       ├── accounts_grpc.pb.go
│   │       │   │       └── orgs.pb.go
│   │       │   └── common
│   │       │       └── v1
│   │       │           └── uuid.pb.go
│   │       └── google
│   │           └── protobuf
│   │               └── timestamp.pb.go
│   ├── providers
│   │   ├── gcp
│   │   │   └── cloudrun
│   │   │       ├── activities.go
│   │   │       ├── constructor.go
│   │   │       ├── resource.go
│   │   │       └── workflows.go
│   │   ├── github
│   │   │   ├── README.md
│   │   │   ├── activities.go
│   │   │   ├── config
│   │   │   │   ├── config.go
│   │   │   │   └── new.go
│   │   │   ├── defs
│   │   │   │   ├── events.go
│   │   │   │   ├── events_test.go
│   │   │   │   ├── github.go
│   │   │   │   ├── openapi.gen.go
│   │   │   │   ├── testdata
│   │   │   │   │   └── create-branch.json
│   │   │   │   ├── timestamp.go
│   │   │   │   └── workflows.go
│   │   │   ├── doc.go
│   │   │   ├── entity.go
│   │   │   ├── errors
│   │   │   │   └── errors.go
│   │   │   ├── errors.go
│   │   │   ├── github.go
│   │   │   ├── handler.go
│   │   │   ├── openapi.codegen.yaml
│   │   │   ├── openapi.gen.go
│   │   │   ├── repo__io.go
│   │   │   ├── types.github.go
│   │   │   ├── types.go
│   │   │   ├── webhook.go
│   │   │   ├── workflow_options.go
│   │   │   └── workflows.go
│   │   └── slack
│   │       ├── activities.go
│   │       ├── blocksets.go
│   │       ├── client.go
│   │       ├── config.go
│   │       ├── crypto.go
│   │       ├── doc.go
│   │       ├── errors.go
│   │       ├── handler.go
│   │       ├── helpers.go
│   │       ├── logger.go
│   │       ├── openapi.codegen.yaml
│   │       └── openapi.gen.go
│   ├── shared
│   │   ├── doc.go
│   │   ├── errors.go
│   │   ├── ids.go
│   │   ├── logger
│   │   │   ├── echo.go
│   │   │   ├── gcp.go
│   │   │   └── utils.go
│   │   ├── queue
│   │   │   ├── core.go
│   │   │   ├── init.go
│   │   │   ├── mutex.go
│   │   │   └── providers.go
│   │   ├── service
│   │   │   └── service.go
│   │   ├── singletons.go
│   │   ├── templates
│   │   │   ├── constants.tmpl
│   │   │   ├── echo
│   │   │   │   ├── echo-interface.tmpl
│   │   │   │   ├── echo-register.tmpl
│   │   │   │   └── echo-wrappers.tmpl
│   │   │   ├── imports.tmpl
│   │   │   └── typedef.tmpl
│   │   ├── temporal
│   │   │   └── temporal.go
│   │   ├── testing.go
│   │   ├── types.gen.go
│   │   ├── types.go
│   │   └── utils.go
│   └── testutils
│       ├── containers.go
│       ├── db.go
│       └── types.go
├── taskfile.yml
└── tmp
    ├── api
    ├── api-build.log
    ├── jobs-migrate-build.log
    ├── migrate
    ├── mothership
    └── worker-mothership-build.log
```
