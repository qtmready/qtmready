# Commit Message

We are following [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/#specification). The main reason for doing is automating our release notes and change log.

As per the spec, a commit message looks like this.

```log
prefix(scope): message

[optional body]

[optional footer(s)]
```

For us, the prefixes are

- feat (for feature)
- fix (represents a code fix)
- docs (represent documentations)
- build (building & contrinous integration)
- ops (operations)
- style (code formatting or refactoring)
- revert (when we are reverting a commit)
- chore (commits that do not fall under any of the other categories)

The scope defines the area we are working on. For us, they are

- core (the core product)
- api (http api)
- providers (3rd party integration)
- cli (command line interface)
- workflows (workflows)

One example of this is,

```docs
feat(cli): added version command to represent release channel

BREAKING CHANGE: `--version is no longer available.
```
