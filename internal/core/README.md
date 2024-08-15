# Core

The core module serves as backbone for the entire applications. The entire premise revolves around building the common primitives for various components involves throughout the entire release orchestration. The application is built on top of temporal to provide us exeuction guarantes and to make atomic execution easy.
For classification fo different events, We have taken the inspiration from [CD Events Specifications](https://cdevents.dev/docs), and then added on top of it. These are

1. **Source Code Control Events**:

   - These events track changes in the source code repository. They include:
     - **Push Events**: Triggered when code is pushed to a repository.
     - **Pull Request Events**: Triggered when a pull request is created, updated, or closed.
     - **Branch Creation/Deletion Events**: Indicate when branches are created or deleted.

2. **Continuous Integration (CI) Events**:

   - These events are related to the integration of code changes and include:
     - **Build Events**: Indicate the start and completion of build processes.
     - **Test Events**: Capture the results of automated tests run against the codebase.

3. **Continuous Deployment (CD) Events**:

   - These events manage the deployment of applications to production environments. They include:
     - **Deployment Events**: Indicate when an application is deployed to a specific environment.
     - **Rollback Events**: Triggered when a deployment is rolled back due to issues.

4. **Quality Assurance Events**:

   - These events focus on maintaining the quality of the software and include:
     - **Code Quality Events**: Capture metrics related to code quality, such as static analysis results.
     - **Security Scan Events**: Indicate the results of security scans performed on the codebase.

5. **Operational Events**:

   - These events are related to the operational aspects of running applications in production. They include:
     - **Monitoring Events**: Capture metrics and logs from running applications.
     - **Incident Management Events**: Triggered when incidents are detected or resolved.

6. **Ticketing Events**:

   - These events relate to issue tracking and management, including:
     - **Issue Creation/Update Events**: Indicate when issues or tickets are created or updated in a tracking system.

7. **Communication Events**:
   - These events facilitate communication between team members and stakeholders, including:
     - **Notification Events**: Trigger notifications sent to communication platforms like Slack or Microsoft Teams.

This is reflected in our directory structure.

```bash
core/
├── cd/             # deployment events
├── ci/             # continuous integration events
├── code/           # source code control events
├── comm/           # slack, teams etc
├── errratic/      # errratic, get it? a palyful combinatin of error and erratic. standardizes errors and error handling.
└── kernel/         # kernel. works on DI principle
└── mutex/          # distributed mutex that can hold a lock across different worklows
└── timers/         # enhanced timer and interval functionality for temporal workflows
├── ops/            # operations events
├── quality/        # testing events
├── tickets/        # ticket events
└── utils/          # shared utilities and helpers
└── web/            # http endpoints
└── ws/             # http endpoints
```

> [!NOTE]  
> All the cdevents specs are map to folders, albeit some are missing, but we do have other functionaltiy,

## Code

The `code` module implements Temporal workflows as controllers to manage different aspects of a Git repository. The key controllers are:

- **TrunkCtrl**: Manages events for the main branch, handling push and branch creation/deletion events.
- **RepoCtrl**: Oversees events for the entire repository, including push events, pull requests, and branch management.
- **BranchCtrl**: Processes events specific to individual branches, such as push, rebase, and pull request events.
- **QueueCtrl**: Manages a queue of pull requests, ensuring they are processed sequentially and efficiently.

Each of these controllers creates a localized state that allows for the management of their respective events. They provide helper functions to mutate state as needed. Each state inherits from `BaseCtrl`, which offers common functionalities such as standardized logging and executing activities or receiving signals with logging capabilities.
