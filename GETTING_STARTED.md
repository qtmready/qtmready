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
├── api
│   ├── openapi
│   │   ├── auth
│   │   │   └── v1
│   │   ├── core
│   │   │   └── v1
│   │   ├── github
│   │   │   └── v1
│   │   ├── shared
│   │   │   └── v1
│   │   └── slack
│   │       └── v1
│   └── proto
│       └── ctrlplane
│           ├── auth
│           │   └── v1
│           └── common
│               └── v1
├── cmd
│   ├── api
│   │   ├── alpha
│   │   └── legacy
│   ├── e2e
│   │   ├── mutex
│   │   ├── orm
│   │   └── qtm
│   ├── jobs
│   │   ├── migrate
│   │   └── openapi
│   └── workers
│       ├── mothership
│       └── sentinel
├── deploy
│   ├── air
│   ├── cassandra
│   ├── postgres
│   │   └── init.sql
│   └── temporal
│       └── dynamicconfig
├── docs
│   └── src
│       ├── contributing
│       ├── design
│       │   └── plantuml
│       └── ref
├── internal
│   ├── auth
│   ├── core
│   │   ├── code
│   │   ├── comm
│   │   ├── defs
│   │   ├── events
│   │   ├── kernel
│   │   ├── mutex
│   │   ├── periodic
│   │   ├── timers
│   │   ├── web
│   │   └── ws
│   ├── db
│   │   ├── config
│   │   ├── entities
│   │   ├── fields
│   │   ├── migrations
│   │   │   ├── cassandra
│   │   │   └── postgres
│   │   ├── queries
│   │   └── sequel
│   ├── durable
│   │   └── connection
│   ├── erratic
│   ├── nomad
│   │   ├── convert
│   │   ├── handler
│   │   └── proto
│   │       ├── ctrlplane
│   │       │   ├── auth
│   │       │   │   └── v1
│   │       │   └── common
│   │       │       └── v1
│   │       └── google
│   │           └── protobuf
│   ├── providers
│   │   ├── gcp
│   │   │   └── cloudrun
│   │   ├── github
│   │   │   ├── config
│   │   │   ├── defs
│   │   │   │   └── testdata
│   │   │   └── errors
│   │   └── slack
│   ├── shared
│   │   ├── logger
│   │   ├── queue
│   │   ├── service
│   │   ├── templates
│   │   │   └── echo
│   │   └── temporal
│   └── testutils
└── tmp
```
