# Overview

> Quantm is an opinionated OpenGitOps frameowork to build _durable delivery_ for distributed systems using a combination of _merge queues_, _version sets_, _immuatabe infrastructure_, _progressive rollouts_ and _one click rollbacks_.

## Merge Queues: Orderly Code Flow

- **Enforce FIFO (First In, First Out):** Pull requests get merged in the order they arrive, ensuring predictability and control.
- **Automate Validation:** Rigorous testing and scrutiny catch bugs early, keeping your codebase clean.
- **Resolve Simple Conflicts Automatically:** Minimize manual intervention and save valuable time.

## Version Sets: Treat The entire stack as One Unit

- **Provide Snapshots of Code State:** Track precise code versions for each microservice, highlighting deviations from the overall codebase.
- **Enable Dependency Management:** Identify and address potential compatibility issues proactively.
- **Facilitate Structured Rollouts:** Coordinate deployments effectively and minimize disruptions.

To learn more, See [Keeping the master green at scale](https://dl.acm.org/ft_gateway.cfm?CFID=58774186&CFTOKEN=b4523e763804e44b-116C8152-E1B2-E079-0753A4F82D863A92&dwn=1&ftid=2045461&id=3303970&uclick_id=1ce424ac-ddd3-4bb2-9368-e5ef47020e2b) at Uber.

## Immutable Infrastructure: Clean Environments for Every Deployment\*\*

Inspired by [Amazon's Brazil Build System](https://gist.github.com/terabyte/15a2d3d407285b8b5a0a7964dd6283b0) and treating infrastrcuture as just another software layer, we

- **Isolate Deployments:** Each version set gets its own unique infrastructure, preventing conflicts and ensuring clean deployments.
- **Prevent Configuration Drift:** Eliminate unintended changes and maintain stable environments.

## One-Click Rollbacks: Rapid Reversal of Issues\*\*

With a combination of version sets and immutable infrastructure, since infra is isolated per version set, we can

- **Instantly Revert to Previous Versions:** Quickly address any problems that arise during or after deployment.
- **Minimize Downtime and Impact:** Restore functionality efficiently and minimize negative consequences.
