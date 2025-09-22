# Contributing to FLOSS

We welcome contributions to FLOSS! This guide explains how to contribute to the project.

## Getting Started

1. Fork the repository
2. Clone your fork
3. Create a new branch for your feature or bug fix
4. Make your changes
5. Test your changes
6. Submit a pull request

## Development Environment

### Prerequisites

- Go 1.23 or later
- Node.js (for some development tools)
- Task (https://taskfile.dev/)

### Building

```bash
task build
```

### Running Tests

```bash
task test
```

### Linting

```bash
task lint
```

Or to automatically fix issues:

```bash
task lint-fix
```

### Formatting

```bash
task fmt
```

## Code Style Guidelines

See [FLOSS Development Guide](../../FLOSS.md) for detailed code style guidelines.

## Documentation

When adding new features or modifying existing ones, please update the documentation accordingly.

## Testing

All code changes should include appropriate tests:
- Unit tests for new functions
- Integration tests for new features
- Regression tests for bug fixes

## Pull Request Process

1. Ensure your code follows the style guidelines
2. Run all tests to ensure nothing is broken
3. Write a clear, descriptive commit message
4. Submit a pull request with a detailed description of your changes

## Reporting Issues

If you find a bug or have a feature request, please open an issue on GitHub with:
- A clear title
- Detailed description
- Steps to reproduce (for bugs)
- Expected vs actual behavior (for bugs)

## Code of Conduct

Please follow our Code of Conduct when contributing to FLOSS.