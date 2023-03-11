# Why?

We are a big fans of linear commit history. After experimenting with multiple merge strategies, we have come to realize the best strategy for easier reverts appears to be Squash Merge. Combining the [Squash Merge](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/configuring-pull-request-merges/configuring-commit-squashing-for-pull-requests) option on Github to Pull Request Title and Description with [conventional-commits](https://www.conventionalcommits.org/en/v1.0.0/#specification), it makes a powerful combo to publish [changeLogs](https://keepachangelog.com/en/1.0.0/) and releases.

In order to leverage these tools, we have formulated a couple of strategies.

## Title

For us the title of each PR must look

```txt
prefix(scope): [shortcut story number] short description
```

### Prefix

For us, we have finalized the following prefixes

```txt
feat (for feature)
fix (represents a bug fix)
spec (specification design)
docs (represent documentations)
ops (build, release & infra)
style (code formatting or refactoring)
revert (when we are reverting a commit)
chore (commits that do not fall under any of the other categories)
```

conventional-commit tools allows us to ignore certain prefixes. Generally, feats, fix and reverts are picked.

### Scope

The scopes for us are

```txt
core (core functionality, mostly workflow related to mothership)
api (rest api)
providers/:provider (saas providers with :provider representing the provider). eg for aws it would sound like (providers/aws)

```

### Shortcut Story Number

Each PR Title must contain the shortcut story title, e.g. If the shortcut story number is 123the shortcut story number will be sc-123. By leveraging Github's Autolinks feature, we can link the merge commit back to the shortcut story. If your repository is not configured for the said feature, ask the project manager.

## Description

Since we aim to be compatible with tooling around conventional-commits, therefore, it is recommended that we follow the recommendations as provided them. i.e.

```txt
[body]

[optional footer(s)]

Notice the space between body and footers.
```

### Body

Sometimes the short description cannot provide the clear description. In this case, the body must provide a clear description of the intended purpose of the code. This is important because code reviews can already be challenging. Without a clear description, the chances of a pull request (PR) being approved are low. This is because the reviewer has three tasks:

- Understanding the purpose of the code changes in the PR by reading the code.
- Going back and forth between Shortcut & Github to gather context.
- Trying to determine if the code does what it’s supposed to do based on the reviewer’s general views on the universe and the current geopolitical climate.

Without a clear statement that explains the purpose of the code, it is impossible to review it for correctness. The reviewer must rely on assumptions about what is correct.

In addition, providing a clear description in the body of the PR makes the lives of product managers much easier as the conventional-commits changelog tools will be able to provide them with much greater context to do product announcements.

### Footer(s)

Footers comes in handy when we need there is a Breaking Change or we need to rollback a commit.

The prefix for Breaking Change is BREAKING CHANGE and for rollbacks, they are generally refs to previous commits using Refs.

## Example

For shortcut story 123, requesting to change the encryption scheme for the authentication middleware from BCRYPT to ARGON2, the contributor changes the encryption scheme to BCRYPT. In this case, we already have customers with existing API Keys, the change would be a breaking change which would require the customers for new API Keys. One example on what the title must look like is

```text
feat(api): [sc-123] changed the default encryption scheme to ARGON2
```

and the description would look like

```text
updated the echo's authentication middleware to use argon2 encryption.

BREAKING CHANGE: All our customers will need to request new API Keys.
```

Although this is a convoluted example but it gives the idea of what makes a good PR.

> This PR should never go through. The reviewer should ask the contributor to add a flag against each user to mark if he/she has newer keys. If not, force the user to update to new keys. work with the community to notify the change.
