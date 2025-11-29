# Contributing to Announcable

Thank you for your interest in contributing to Announcable! This document provides guidelines and instructions for contributing.

## Getting Started

1. Fork the repository
2. Clone your fork locally
3. Set up the development environment following the instructions in the README

## Development Setup

### Prerequisites

- Go 1.21+
- Node.js 18+
- Docker and Docker Compose
- Make

### Running Locally

```bash
cd backend

# Terminal 1: Start Docker services
make dev-services

# Terminal 2: Start Go backend with hot-reload
make dev-air

# Terminal 3: Start Vite for CSS/JS hot-reload
npm run dev
```

## How to Contribute

### Reporting Bugs

- Check existing issues to avoid duplicates
- Use a clear, descriptive title
- Describe the steps to reproduce the issue
- Include your environment details (OS, browser, versions)
- Add screenshots if applicable

### Suggesting Features

- Open an issue describing the feature
- Explain the use case and why it would be valuable
- Be open to discussion about implementation approaches

### Submitting Code

1. Create a branch for your changes
2. Write clear, descriptive commit messages
3. Ensure your code follows the existing style
4. Add tests for new functionality
5. Run the test suite and ensure it passes
6. Submit a pull request

### Pull Request Guidelines

- Reference any related issues
- Describe what changes you made and why
- Keep PRs focused on a single concern
- Be responsive to feedback during review

## Code Style

### Go

- Follow standard Go conventions
- Use `gofmt` for formatting
- Handler files follow the naming pattern `handle-<action>.go`

### TypeScript/Lit (Widget)

- Use TypeScript for all new code
- Follow the existing Lit web component patterns
- Run `npm run lint` before submitting

### CSS

- Follow the component-based architecture in `backend/assets/css/`
- Use CSS variables from `base/variables.css`
- Keep styles scoped and modular

## Database Migrations

If your changes require database modifications:

```bash
make migrations-new name=<descriptive_name>
```

Test migrations in both directions (up and down).

## Questions?

If you have questions about contributing, feel free to open an issue for discussion.

## License

By contributing to Announcable, you agree that your contributions will be licensed under the AGPL-3.0 license.
