# Contributing to ads-go

Thank you for your interest in contributing to ads-go! This document provides guidelines and instructions for contributing.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [How to Contribute](#how-to-contribute)
- [Coding Guidelines](#coding-guidelines)
- [Testing](#testing)
- [Pull Request Process](#pull-request-process)
- [Reporting Bugs](#reporting-bugs)
- [Feature Requests](#feature-requests)
- [Questions](#questions)

## Code of Conduct

This project adheres to a [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code. Please report unacceptable behavior to the project maintainers.

## Getting Started

### Prerequisites

Before you begin, ensure you have the following installed:

- **Go**: Version 1.23 or later ([download](https://go.dev/dl/))
- **Git**: For version control
- **TwinCAT**: For testing (optional but recommended)
  - TwinCAT 3 XAE (Engineering)
  - TwinCAT 3 XAR (Runtime) or local runtime
  - ADS Router configured

### Recommended Tools

- **golangci-lint**: For linting ([installation](https://golangci-lint.run/welcome/install/))
- **VS Code** or **GoLand**: Recommended IDEs with Go support

## Development Setup

1. **Fork the repository** on GitHub

2. **Clone your fork**:
   ```bash
   git clone https://github.com/YOUR-USERNAME/ads-go.git
   cd ads-go
   ```

3. **Add upstream remote**:
   ```bash
   git remote add upstream https://github.com/JarmoCluyse/ads-go.git
   ```

4. **Install dependencies**:
   ```bash
   go mod download
   ```

5. **Verify your setup**:
   ```bash
   go test ./...
   go build ./...
   ```

## How to Contribute

### 1. Find or Create an Issue

- Check [existing issues](https://github.com/JarmoCluyse/ads-go/issues) for something to work on
- For new features or significant changes, create an issue first to discuss
- Comment on the issue to let others know you're working on it

### 2. Create a Branch

```bash
git checkout -b feature/your-feature-name
# or
git checkout -b fix/your-bug-fix
```

Use descriptive branch names:
- `feature/add-handle-management` for new features
- `fix/subscription-memory-leak` for bug fixes
- `docs/update-examples` for documentation
- `refactor/improve-error-handling` for refactoring

### 3. Make Your Changes

- Write clean, readable code
- Follow the coding guidelines below
- Add tests for new functionality
- Update documentation as needed
- Keep commits focused and atomic

### 4. Test Your Changes

```bash
# Run all tests
go test ./...

# Run tests with race detection
go test -race ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 5. Commit Your Changes

Write clear, descriptive commit messages:

```bash
git commit -m "feat: add handle management for symbols

- Implement handle caching for improved performance
- Add automatic handle lifecycle management
- Update documentation with handle examples

Closes #123"
```

Follow [Conventional Commits](https://www.conventionalcommits.org/):
- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation changes
- `test:` - Adding or updating tests
- `refactor:` - Code refactoring
- `perf:` - Performance improvements
- `chore:` - Maintenance tasks

## Coding Guidelines

### Go Style

- Follow [Effective Go](https://go.dev/doc/effective_go)
- Use `gofmt` for formatting (enforced by CI)
- Use `go vet` to check for common issues
- Follow Go naming conventions

### Code Quality

1. **Formatting**: All code must pass `gofmt -l .`
   ```bash
   gofmt -w .
   ```

2. **Linting**: Code should pass `golangci-lint`
   ```bash
   golangci-lint run
   ```

3. **Vetting**: Code must pass `go vet`
   ```bash
   go vet ./...
   ```

### Best Practices

- **Error handling**: Always handle errors explicitly
- **Comments**: 
  - Add package documentation to `doc.go` files
  - Document all exported functions, types, and constants
  - Use complete sentences with proper punctuation
- **Testing**: Aim for >80% code coverage for new code
- **Concurrency**: Use goroutines and channels safely; avoid data races
- **Context**: Use `context.Context` for operations that can be cancelled or have timeouts

### ADS-Specific Guidelines

- **Type Safety**: Ensure type conversions between Go and PLC types are safe
- **Protocol Compliance**: Follow the ADS specification strictly
- **Error Codes**: Map ADS error codes to descriptive Go errors
- **Documentation**: Explain TwinCAT-specific concepts for users unfamiliar with ADS

## Testing

### Writing Tests

- Place tests in `*_test.go` files
- Use table-driven tests where appropriate
- Test both success and error cases
- Use meaningful test names: `TestClient_ReadVariable_Success`

### Example Test Structure

```go
func TestClient_ReadVariable(t *testing.T) {
    tests := []struct {
        name    string
        symbol  string
        want    interface{}
        wantErr bool
    }{
        {
            name:    "read bool variable",
            symbol:  "MAIN.bVariable",
            want:    true,
            wantErr: false,
        },
        // More test cases...
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### Integration Tests

Some tests require a real TwinCAT PLC:
- Mark integration tests with build tags: `//go:build integration`
- Document required PLC setup in test comments
- Use the example TwinCAT project for consistent testing

## Pull Request Process

### Before Submitting

1. **Sync with upstream**:
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

2. **Run all checks**:
   ```bash
   go test ./...
   go vet ./...
   gofmt -l .
   golangci-lint run
   ```

3. **Update documentation**:
   - Update README.md if needed
   - Update CHANGELOG.md following [Keep a Changelog](https://keepachangelog.com/)
   - Add code comments/doc comments

### Submitting Your PR

1. **Push to your fork**:
   ```bash
   git push origin feature/your-feature-name
   ```

2. **Create a Pull Request** on GitHub

3. **Fill out the PR template** completely:
   - Describe what the PR does
   - Link to related issues
   - Note any breaking changes
   - Describe how you tested the changes

4. **Respond to feedback**: Address review comments promptly

### PR Review Process

- Maintainers will review your PR
- CI checks must pass (tests, linting, build)
- At least one approval is required
- Changes may be requested
- Once approved, a maintainer will merge your PR

## Reporting Bugs

Use the [Bug Report template](.github/ISSUE_TEMPLATE/bug_report.md) and include:

- **Description**: Clear description of the bug
- **Steps to reproduce**: Detailed steps to reproduce the issue
- **Expected behavior**: What you expected to happen
- **Actual behavior**: What actually happened
- **Environment**: 
  - Go version
  - ads-go version
  - Operating System
  - TwinCAT version
  - Connection details (local/remote)
- **Code sample**: Minimal reproducible example
- **Logs/errors**: Relevant error messages or logs

## Feature Requests

Use the [Feature Request template](.github/ISSUE_TEMPLATE/feature_request.md) and include:

- **Problem**: What problem does this solve?
- **Solution**: Proposed solution
- **Alternatives**: Alternative solutions considered
- **Use case**: How would this be used?
- **ADS considerations**: Any ADS protocol implications

## Questions

Have a question? Here's where to ask:

- **GitHub Discussions**: For general questions and discussions
- **Issue Tracker**: For bug reports and feature requests
- **Documentation**: Check the [README](README.md) first

## License

By contributing to ads-go, you agree that your contributions will be licensed under the [MIT License](LICENSE).

---

Thank you for contributing to ads-go! Your efforts help make industrial automation more accessible with Go.
