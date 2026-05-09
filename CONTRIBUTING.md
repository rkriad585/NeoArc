# Contributing to NeoArc

Thank you for considering contributing to NeoArc! We welcome contributions of all kinds — bug fixes, new features, documentation improvements, and more.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Coding Guidelines](#coding-guidelines)
- [Pull Request Process](#pull-request-process)
- [Reporting Issues](#reporting-issues)

## Code of Conduct

This project and everyone participating in it is governed by the [NeoArc Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code. Please report unacceptable behavior to mdriyadkhan585@gmail.com.

## Getting Started

1. Fork the repository on GitHub.
2. Clone your fork locally:
   ```bash
   git clone https://github.com/your-username/NeoArc.git
   cd NeoArc
   ```
3. Add the original repository as an upstream remote:
   ```bash
   git remote add upstream https://github.com/rkriad585/NeoArc.git
   ```

## Development Setup

### Server (Python)

```bash
cd neoarc-server
python -m venv .venv
source .venv/bin/activate  # or .venv\Scripts\activate on Windows
pip install -e .
cp .env.example .env
python main.py
```

The server starts at `http://localhost:59248`.

### CLI (Go)

```bash
cd neoarc-cli
go build -o bin/neoarc ./cmd/neoarc
./bin/neoarc config http://localhost:59248
```

### Running Tests

```bash
# Server tests (66 tests)
cd neoarc-server
python -m pytest tests/ -v

# CLI unit tests (30 tests)
cd neoarc-cli
go test ./internal/cli/ -v

# CLI integration tests (3 tests)
go test ./tests/ -v
```

## Coding Guidelines

### Python

- Follow [PEP 8](https://peps.python.org/pep-0008/).
- Use type hints where practical.
- Use `snake_case` for functions and variables, `PascalCase` for classes.
- Format imports: standard library, third-party, local (separated by blank lines).

### Go

- Run `gofmt` before committing.
- Follow [Effective Go](https://go.dev/doc/effective_go) conventions.
- Use `camelCase` for unexported, `PascalCase` for exported names.
- Handle all errors — never use `_` for error returns unless intentional.

### General

- Keep functions focused and small.
- Write tests for all new functionality.
- Use clear, descriptive variable names.
- Avoid adding comments that explain "what" the code does — the code should speak for itself. Comments should explain "why".

## Pull Request Process

1. Create a feature branch from `main`:
   ```bash
   git checkout -b feature/my-feature
   ```

2. Make your changes, ensuring:
   - All existing tests pass.
   - New tests cover your changes.
   - Code follows the style guidelines above.
   - Documentation is updated if the public API or behavior changes.

3. Commit with a descriptive message:
   ```bash
   git commit -m "Add feature: brief description of what and why"
   ```

4. Push to your fork and open a pull request:
   ```bash
   git push origin feature/my-feature
   ```

5. In your PR description, explain:
   - What the change does.
   - Why it's needed (include issue number if applicable).
   - How you tested it.

6. A maintainer will review your PR. Address any feedback and push updates as needed.

## Reporting Issues

When reporting a bug, please include:

- A clear, descriptive title.
- Steps to reproduce the issue.
- Expected behavior and what actually happened.
- Environment details (OS, Python/Go version, browser if relevant).
- Screenshots or logs if applicable.

Feature requests are welcome! Open an issue with the tag "enhancement" and describe the use case.

---

Thank you for helping make NeoArc better!
