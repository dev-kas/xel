# Contributing to Xel

Thank you for your interest in contributing to Xel! This document provides guidelines and instructions for contributing to the project.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Pull Request Process](#pull-request-process)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Building](#building)
- [Documentation](#documentation)
- [Release Process](#release-process)
- [Issue Reporting](#issue-reporting)
- [Feature Requests](#feature-requests)

## Code of Conduct

By participating in this project, you are expected to uphold our Code of Conduct:

- Be respectful and inclusive
- Be patient and welcoming
- Be thoughtful
- Be collaborative
- When disagreeing, try to understand why

## Getting Started

1. **Fork the repository** on GitHub
2. **Clone your fork** locally
   ```bash
   git clone https://github.com/YOUR-USERNAME/xel.git
   cd xel
   ```
3. **Add the upstream repository** as a remote
   ```bash
   git remote add upstream https://github.com/dev-kas/xel.git
   ```
4. **Install dependencies**
   ```bash
   go mod download
   ```

## Development Workflow

1. **Create a branch** for your feature or bugfix
   ```bash
   git checkout -b feature/your-feature-name
   ```
   or
   ```bash
   git checkout -b fix/your-bugfix-name
   ```

2. **Make your changes** and commit them with clear, descriptive commit messages
   ```bash
   git commit -m "Add feature: your feature description"
   ```

3. **Keep your branch updated** with the upstream repository
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

4. **Push your changes** to your fork
   ```bash
   git push origin feature/your-feature-name
   ```

## Pull Request Process

1. **Submit a pull request** from your forked repository to the main repository
2. **Ensure your PR description** clearly describes the problem and solution
3. **Include issue numbers** in your PR description (e.g., "Fixes #123")
4. **Update documentation** if necessary
5. **Ensure all tests pass** and add new tests for new functionality
6. **Request a review** from maintainers

## Coding Standards

- Follow Go's official [Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Run `go fmt ./...` before committing to ensure consistent formatting
- Use `go vet` and `golint` to check for common issues
- Write clear, descriptive variable and function names
- Add comments for non-obvious code sections
- Keep functions small and focused on a single responsibility

## Testing

- Write tests for all new features and bug fixes
- Ensure all tests pass before submitting a pull request
- Run tests using:
  ```bash
  go test ./...
  ```
- Aim for high test coverage for critical functionality

## Building

The project includes a Makefile with several targets:

- Build for all platforms:
  ```bash
  make build
  ```

- Build for specific platforms:
  ```bash
  make build-mac
  make build-linux
  make build-windows
  ```

- Run tests:
  ```bash
  make test
  ```

## Documentation

- Update documentation for any changes to functionality
- Document all public APIs
- Keep README.md updated with the latest information
- Add examples for new features

## Versioning

Xel follows [Semantic Versioning](https://semver.org/) (SemVer):

- **MAJOR** version for incompatible API changes (X.0.0)
- **MINOR** version for backward-compatible functionality additions (0.X.0)
- **PATCH** version for backward-compatible bug fixes (0.0.X)

The version is set in the main.go file and is overridden during build using ldflags.

## Release Process

Xel uses GitHub Actions for continuous integration and automated releases:

### Continuous Integration

- Every push to the main branch and pull requests are automatically built and tested
- The CI workflow is defined in `.github/workflows/ci.yml`
- CI builds include a development version number based on the git commit hash

### Creating a New Release

Only project maintainers should create official releases. The process is:

1. Ensure all tests pass and the code is ready for release
2. Update version references in documentation if necessary
3. Create and push a new tag with the version number following SemVer:
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```
4. The release workflow (`.github/workflows/release.yml`) will automatically:
   - Build binaries for all supported platforms
   - Create a GitHub Release with the binaries attached
   - Generate release notes based on commits since the last release

### Versioning

Xel follows [Semantic Versioning](https://semver.org/):

- **MAJOR** version for incompatible API changes
- **MINOR** version for backward-compatible functionality additions
- **PATCH** version for backward-compatible bug fixes

### Release Candidates

For significant changes, consider using release candidates:
```bash
git tag v1.0.0-rc1
git push origin v1.0.0-rc1
```

## Issue Reporting

When reporting issues, please include:

1. **Description** of the issue
2. **Steps to reproduce** the issue
3. **Expected behavior**
4. **Actual behavior**
5. **Environment details** (OS, Go version, etc.)
6. **Logs or error messages**
7. **Possible solutions** (if you have any ideas)

## Feature Requests

For feature requests, please include:

1. **Description** of the feature
2. **Rationale** for the feature (why it's needed)
3. **Example usage** of the feature
4. **Possible implementation** details (if you have ideas)

Thank you for contributing to Xel!