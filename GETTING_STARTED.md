# Quantm.io - Getting Started

## Prerequisites

Before you begin, ensure you have the following software installed:

- Go (v1.23+): [https://go.dev/doc/install](https://go.dev/doc/install)
- Docker: [https://docs.docker.com/get-docker/](https://docs.docker.com/get-docker/) (for containerization)
- taskfile: [https://taskfile.dev/#](https://taskfile.dev/#/) (for task automation)
- buf: [https://docs.buf.build/installation](https://docs.buf.build/installation) (for managing protobuf files)

## Setting Environment

### Clone the Repository with submodules

```bash
git clone https://github.com/quantmHQ/quantm.git
cd quantm
```

### Configure Environment Variables

Copy the `.env.example` file to `.env` and configure it with your specific settings:

```bash
cp .env.example .env
```

### Required Services

Start the necessary services using Docker Compose. This inlcudes

- temporal-db: Postgres 16.x for required for temporal
- temporal: 1.25.x, it drives the main logic of the application
- db: Postgres 16.x for the main database

  ```bash
  docker-compose up -d
  ```

### Running the Application

Once you've completed the setup steps, you can run the Quantm.io application using the following methods:

#### With Hot Reload

To run the application in development mode with hot reload functionality, use the `air` command:

```bash
air
```

You can also use command-line flags to run specific modules during development:

- `--migrate` or `-m`: To migrate the database.
- `--run` or `-r`: To run a specific part of the application. For example:
  - `nomad`: Run the Nomad module.
  - `mothership`: Run the Mothership module.
  - `web`: Run the web server.

#### More Control

To run the application in production mode, use the `go run` command:

```bash
go run ./cmd/quantm # or go run ./cmd/quantm --help for more options
```

> [!IMPORTANT]
> For detailed instructions on running specific modules, configurations, and advanced usage, refer to the project documentation.

## Project Structure

This is the overall project hierarchy:

- `cmd`: The project executable. This project exposes a single binary, and different command-line arguments can be used to modify its behavior. The single binary itself is responsible for running different modules in separate Go routines and provides error handling and recovery capabilities if a Go routine crashes.
- `deploy`: Holds scripts or configuration related to deployment.
- `docs`: Contains documentation for the project.
- `api`: This is a separate [git repo](https://github.com/quantmHQ/api) that defines and shares protobuf definitions (`ctrlplane`) across different repositories. It serves as the external interface for other applications or services to interact with the data structures.
- `internal`: This folder contains the core logic of the application and is separated from the `api` module to provide a more structured and secure design. The `internal` folder is organized in a layered manner, as described below:

### Understanding `internal`

The `internal` folder is organized into four levels. Higher-level folders can depend on lower-level folders but not vice-versa to avoid circular import issues. This layered approach ensures that each package has a clear responsibility and promotes a clean separation of concerns.

- **Level 0:** - Foundation
- **Level 1:** - Utilities and Core Abstractions
- **Level 2:** - Core Logic
- **Level 3:** - Integration Points
- **Level 4:** - Interfaces

Below is the detailed description of each level:

#### Level 0 - Foundation

- `proto`: Contains generated code from protobuf definitions, acting as the shared data structure for communication within the application.
- `db`: Provides the persistence layer, utilizing `sqlc` for standardized SQL queries and `go-migrate` for database migrations.
- `durable`: Provides the helper function for durable execution of tasks.
- `pulse`: Configuration and utilities for analytics.
- `observe`: Standardizes observability, handling metrics, tracing, and logging across the application.

#### Level 1 - Utilities and Core Abstractions

- `cast`: Facilitates data conversion between protobuf definitions and SQL queries, serving as a bridge between data models and the database.
- `erratic`: Provides standardized error handling, leveraging the `erratic` package for automatic error logging and utilities for custom error types.

#### Level 2 - Core Logic

- `auth`: Implements authentication and authorization mechanisms for securing access to the application.
- `core`: Encompasses the main business logic, defining core functionalities and workflows.

#### Level 3 - Integration Points

- `hooks`: Defines a mechanism for registering and invoking hooks, allowing external systems to integrate with specific functionalities within the `core` package.

#### Level 4 - Interfaces

- `nomad`: Provides gRPC server and client for interacting with Nomad, enabling task management, deployment, and resource control.
- `web`: Provides the HTTP server (webhooks) for exposing API endpoints for external interactions with the application.
