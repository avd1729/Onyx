# Contributing to Onyx MCP Sandbox

Thank you for your interest in contributing! This guide explains how to contribute safely and efficiently to the **Onyx MCP Sandbox**, which executes code in multiple languages within a Docker sandbox.

## Table of Contents

- [How to Contribute](#how-to-contribute)
- [Code Style and Guidelines](#code-style-and-guidelines)
- [Testing](#testing)
- [Adding New Languages](#adding-new-languages)
- [Pull Request Process](#pull-request-process)
- [Reporting Issues](#reporting-issues)
- [Community](#community)

---

## How to Contribute

1. **Fork the Repository**  
   Click "Fork" on the GitHub repository page to create your own copy.

2. **Clone Your Fork**  
   ```bash
   git clone https://github.com/<your-username>/onyx.git
   cd onyx
   ```

3. **Create a Branch**  
   Always create a feature branch for your changes:
   ```bash
   git checkout -b feature/<descriptive-name>
   ```

4. **Make Your Changes**
   * Implement your feature or bug fix.
   * Update tests under the `tests/` folder if applicable.

5. **Run Tests**  
   Ensure all tests pass before submitting:
   ```bash
   go test ./tests/...
   ```

---

## Code Style and Guidelines

* **Go code**: Run `gofmt -s -w .` before committing.
* **Docker commands**: Maintain the current sandboxing conventions:
  * Use `--read-only`, `--tmpfs`, `--network none`, etc.
  * Respect CPU/memory/pids limits for each container.
* **Commit Messages**: Use clear and descriptive messages.
  ```
  Add Java executor with Docker sandboxing
  Fix memory leak in Python executor
  ```

---

## Testing

* All language executors have corresponding tests in `tests/`:
  * `python_test.go`
  * `java_test.go`
  * `cpp_test.go`, `c_test.go`, `js_test.go`, `rust_test.go`
* Use `go test ./tests/...` to run all tests or run individual files:
  ```bash
  go test ./tests/java_test.go
  ```

---

## Adding New Languages

1. Add a new executor struct in `executor/`.
2. Implement the `Execute()` method following existing patterns.
3. Add a corresponding test file in `tests/`.
4. Update `main.go` to include the new language in the runtime selection logic.
5. Add the new Docker image to `utils.IsDockerAvailable()` if needed.

---

## Pull Request Process

1. Fork and create a branch.
2. Implement changes and tests.
3. Run all tests to ensure everything works.
4. Submit a PR against the `main` branch.
5. Wait for review and CI checks.

---

## Reporting Issues

* Use the GitHub Issues page to report bugs or request features.
* Provide:
  * Steps to reproduce
  * Error messages / stack traces
  * Language and environment used

---

## Community

* Be respectful and helpful to other contributors.
* Follow the [Code of Conduct](CODE_OF_CONDUCT.md).
