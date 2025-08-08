# Contributing to Mist Exporter

We welcome contributions! Please follow these guidelines to contribute.

## Development Process

1.  Fork the repository.
2.  Create a new branch for your feature or bug fix: `git checkout -b my-new-feature`.
3.  Make your changes.
4.  Add or update tests for your changes. All contributions should be tested.
5.  Ensure the test suite passes: `go test -v ./...`.
6.  Format your code: `go fmt ./...`.
7.  Commit your changes following the Commit Message Guidelines.
8.  Push to the branch: `git push origin my-new-feature`.
9.  Create a new Pull Request.

## Commit Message Guidelines

This project uses Conventional Commits for its commit messages. This format allows for automatic versioning and changelog generation by GoReleaser.

The commit message should be structured as follows:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

**Common types:**

-   **feat**: A new feature.
-   **fix**: A bug fix.
-   **docs**: Documentation only changes.
-   **style**: Changes that do not affect the meaning of the code (white-space, formatting, etc).
-   **refactor**: A code change that neither fixes a bug nor adds a feature.
-   **perf**: A code change that improves performance.
-   **test**: Adding missing tests or correcting existing tests.
-   **build**: Changes that affect the build system or external dependencies (e.g., `goreleaser`, `go.mod`, `Dockerfile`).
-   **ci**: Changes to our CI configuration files and scripts.

**Example:**

```
feat(switch): Add support for switch metrics

Adds a new WebSocket stream and Prometheus metrics for Juniper Mist switches,
including port statistics and device status.
```