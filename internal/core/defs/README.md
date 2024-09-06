# RFD: Event Schema Versioning

**Author:** Yousuf Jawwad
**Date:** 2023-10-27

## 1. Summary

This document defines the versioning strategy for our event schema, aiming for stability and backward compatibility. We will follow Semantic Versioning (SemVer) with considerations for database impact.

## 2. Motivation

A clear versioning strategy is crucial as our event system evolves to:

- Minimize breaking changes and ensure a smooth experience for data consumers.
- Guide engineering decisions regarding schema modifications.
- Manage the impact of schema changes on our database.

## 3. Versioning Strategy

We will use SemVer (`MAJOR.MINOR.PATCH`) with the following guidelines:

- **PATCH:** Add new event types without modifying existing structures. No database changes.
- **MINOR:** Backward-compatible additions or modifications to existing `subject` or `context`, potentially requiring database schema changes.
- **MAJOR:** Backward-incompatible changes, likely requiring consumer code adaptations and database migrations.

## 4. Schema Evolution Process

1. **Identify Change:** Define the business or technical need for modification.
2. **Impact Assessment:** Determine the SemVer bump (PATCH, MINOR, MAJOR) and analyze the impact on the database and consumers.
3. **Design & Implementation:**
   - **PATCH:** No changes to existing structures or database.
   - **MINOR/MAJOR:** Design and implement schema and database updates. Develop and test data migration plans. Update code to handle new/modified data.
4. **Documentation:** Maintain detailed version history, schema modification descriptions, database migration steps, and consumer impact assessments.
5. **Testing:** Conduct backward compatibility testing, validate database interactions, and test all new code paths.
6. **Deployment & Communication:** Deploy changes in a controlled manner and communicate updates, rationale, impact, and migration guidance to all stakeholders.

## 5. 1.0.0 and Beyond

After reaching 1.0.0, signifying a stable schema:

- **MINOR:** Backward-compatible additions to existing structures or new event types without database schema changes.
- **MAJOR:** Reserved for unavoidable backward-incompatible changes with thorough planning, migration strategies, and communication.

## 6. Conclusion

This strategy emphasizes careful consideration, impact analysis, and robust engineering practices for schema evolution, minimizing disruption and maintaining data integrity for our event system.
