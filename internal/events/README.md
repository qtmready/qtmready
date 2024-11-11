# QuantmEvents Primer

> [!IMPORTANT]
> This is a work in progress and is subject to change. However, with each iteration, do change the version.

This document explains the design decisions behind the specifications and serves as a guide for future extensions. It's not a formal standard.

## Motivation

Working extensively with [Temporal](https://temporal.io) across a couple of projects revealed several challenges in our payload handling and observability.  We found ourselves repeatedly creating helper functions to enrich logs with contextual information.  While Temporal offers excellent visibility for its workflows and have excellent support for OpenTelemetry, providing observability beyond infrastructure and temporal boundaries proved difficult. Even if we somehow leveraged OpenTelemetry, doing analytics was another problem that required comprehensive engineering effort.

The need for a more streamlined and efficient approach for payload management became increasingly apparent.

### Investigation

We investigated a couple of approaches and found that each approach falling short.

#### Temporal's Internal

Temporal provides excellent workflow visibility through replays. We considered leveraging this visibility in two ways:

First, we explored an "in-flight" approach, using a combination of interceptors and custom codecs.  However, this approach would have meant directly sitting in the path of workflow execution, which raised concerns about maintainability.  We were unsure if this approach would be robust enough to handle the complexities of our event management system.  Furthermore, introducing a critical dependency directly in the path of execution could potentially introduce bugs and might have slowed us down during the update cycle of Temporal.

Alternatively, we considered accessing Temporal's visibility datastore "at rest".  This approach would have allowed us to build advanced observability and analytics solutions as a separate system.  However, this would have meant losing the ability to perform contextual logging.  Contextual logging is crucial for providing rich insights into event occurrences, and we were hesitant to sacrifice this functionality.

#### CloudEvents

[CloudEvents](https://cloudevents.io), a highly flexible and widely adopted standard, offered a promising foundation.  However, its broad applicability meant that implementing it for our specific needs required custom tooling and extensive knowledge around CloudEvents existing ecosystem. Both of which are expensive to come by. Besides Custom Events

#### Conclusion

To address our unique requirements, we developed our own event schema called QuantmEvents. This schema leverages the core principles of CloudEvents while incorporating elements of distributed tracing. However, it deviates from the standard CloudEvents approach by including metadata directly within the payload, rather than within the header. While this decision accommodates our specific needs, it does potentially impact transmission efficiency compared to the established CloudEvents design.

## Design Goals

QuantmEvents is a unified event schema designed to standardize temporal workloads with the following goals in mind:

### Contextual Clarity

The event's position within a workflow should be clear directly from the event payload itself.

### Lineage Tracking

QuantmEvents simplifies tracing the flow of events to understand the progression of activities, drawing inspiration from distributed tracing principles.

### Standardized Logging

A unified structure for contextual logging enhances observability and debugging.

### Analytics-Ready Data

QuantmEvents' standardized schema facilitates the creation of analytics dashboards and reports.

### World-class Developer Experience

The tooling around consuming and producing events should be inituative and easy to use.

### Improvements over CloudEvents

QuantmEvents introduces two core concepts that differentiate it from other event schemas:

#### Subject

QuantmEvents defines a `subject` representing the entity or resource associated with the event. This clearly identifies the event's ownership, ensuring it is directly tied to the relevant entity.

#### Provider

QuantmEvents is primarily designed for integrations, so it introduces the concept of a `hooks`. The `hooks` specifies the integration platform responsible for generating or receiving the event, providing essential context for understanding its origin or destination.

## QuantmEvents Concepts

> [!NOTE]
> Work in progress


## Architecture

> [!NOTE]
> Work in progress

## Core Attributes

### Event Structure

| Attribute | Description | Required | Default Value | Validation Error |
| --------- | --------------------------------------------------------------------------------------------------------------------------- | -------- | ------------- | --------------------------- |
| `version` | Indicates the schema version for compatibility. | Yes | N/A | Missing version field |
| `id` | A globally unique identifier for the event, using UUID v7 as specified in RFC 9562. | Yes | N/A | Missing or invalid event ID |
| `context` | Contains additional contextual data, crucial for understanding and processing the event. | Yes | N/A | Missing context field |
| `subject` | Defines the entity or resource affected by the event. | Yes | N/A | Missing subject field |
| `payload` | The specific data associated with the event. Designed for flexibility and extensibility to accommodate various event types. | Yes | N/A | Missing payload field |
| `timestamp` | The event occurrence time. | Yes | N/A | Missing timestamp field |

#### Subject Attributes

| Attribute | Description | Required | Default Value | Validation Error |
| --------- | ------------------------------------------------------------------------------------ | -------- | -------------------------------------- | ----------------------------- |
| `id` | The unique identifier of the subject (e.g., repository ID). | Yes | N/A | Missing or invalid subject ID |
| `name` | The name of the ScyllaDB table representing the subject (e.g., "repos"). | Yes | N/A | Missing subject name |
| `team_id` | The team associated with the subject. It is a UUID value (v4 or v7) as per RFC 9562. | Yes | N/A | Missing team ID |
| `user_id` | The user associated with the subject. | No | `00000000-0000-0000-0000-000000000000` | |

#### Context Attributes

| Attribute | Description | Required | Default Value | Validation Error |
| ----------- | --------------------------------------------------------------------------------------------------------------------------- | -------- | -------------------------------------- | ----------------------- |
| `parent_id` | The ID of the preceding related event (for tracing chains), which can be null uuid if it is the first in a chain of events. | No | `00000000-0000-0000-0000-000000000000` | |
| `hook` | The integration hook (e.g., GitHub, GitLab, GCP). | Yes | N/A | Missing hooks field |
| `scope` | The event category (e.g., branch, pull_request). | Yes | N/A | Missing scope field |
| `action` | The triggering action (e.g., created, updated, merged). | Yes | N/A | Missing action field |
| `source` | The event source. The format is unique to each hook. | Yes | N/A | Missing source field |
| `timestamp` | The event occurrence time. | Yes | N/A | Missing timestamp field |

## What does it look like?

```json
{
  "version": "0.1.0",
  "id": "0b2a799d-40d0-4e4c-9144-11f61f233a69",
  "timestamp": "2022-01-01T00:00:00Z",
  "context": {
    "parent_id": "00112233-4455-6677-8899-aabbccddeeff",
    "hook": 1,
    "scope": "pull_request",
    "action": "opened",
    "source": "https://github.com/quantmhq/quantm"
  },
  "subject": {
    "id": "954e472d-a22f-4696-b447-2c8144d37f47",
    "name": "repos",
    "org_id": "6d1a8193-e542-4e54-817a-a7918b56236c",
    "team_id": "0879053d-41a2-4151-a4b3-5e34252eb639",
    "user_id": "1b6362aa-b556-4e68-b440-949250d44055"
  },
  "payload": {
    "number": 1,
    "branch": "main",
    "title": "Add new feature",
  }
}
```


## Versioning

QuantmEvents follows a Semantic Versioning (SemVer) strategy `MAJOR.MINOR.PATCH` to ensure stability and backward compatibility while managing schema evolution. This strategy is based on the following principles:

### Patch Version (PATCH)

- New event types are added without modifying existing structures.
- No database changes are required.

### Minor Version (MINOR)

- Backward-compatible additions or modifications to existing `subject` or `context` attributes are implemented.
- These changes may require database schema updates.

### Major Version (MAJOR)

- Backward-incompatible changes are introduced, requiring consumer code adaptations and database migrations.
- These changes are reserved for unavoidable situations and involve thorough planning, migration strategies, and communication.

### Schema Evolution Process

1. **Identify Change:** Clearly define the business or technical need for the schema modification.
2. **Impact Assessment:** Determine the SemVer bump (PATCH, MINOR, MAJOR) and analyze the impact on the database and consumers.
3. **Design & Implementation:**
   - **PATCH:** No changes to existing structures or database.
   - **MINOR/MAJOR:** Design and implement schema and database updates. Develop and test data migration plans. Update code to handle new/modified data.
4. **Documentation:** Maintain detailed version history, schema modification descriptions, database migration steps, and consumer impact assessments.
5. **Testing:** Conduct backward compatibility testing, validate database interactions, and test all new code paths.
6. **Deployment & Communication:** Deploy changes in a controlled manner and communicate updates, rationale, impact, and migration guidance to all stakeholders.

### 1.0.0 and Beyond

After reaching version 1.0.0, signifying a stable schema, the following guidelines apply:

- **MINOR:** Backward-compatible additions to existing structures or new event types without database schema changes.
- **MAJOR:** Reserved for unavoidable backward-incompatible changes with thorough planning, migration strategies, and communication.

## Prior Art

- [CloudEvents](https://cloudevents.io)
- [Distributed Tracing](https://www.servicenow.com/products/observability/what-is-distributed-tracing.html)
- Google CloudFunctions Pubsub Payload

```json
{
  "data": {
    "@type": "types.googleapis.com/google.pubsub.v1.PubsubMessage",
    "attributes": {
      "foo": "bar"
     },
     "messageId": "12345",
     "publishTime": "2017-06-05T12:00:00.000Z",
     "data": "somebase64encodedmessage"
  },
  "context": {
    "eventId": "12345",
    "timestamp": "2017-06-05T12:00:00.000Z",
    "eventTypeId": "google.pubsub.topic.publish",
    "resource": {
      "name": "projects/myProject/topics/myTopic",
      "service": "pubsub.googleapis.com"
    }
  }
}
```


## Open Questions

> [!NOTE]
> Work in Progress
