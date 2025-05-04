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