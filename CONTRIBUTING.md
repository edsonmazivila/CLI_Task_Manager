# Contributing to CLI Task Manager

Thank you for your interest in contributing to the CLI Task Manager project! This document provides guidelines and instructions for contributing.

---

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Branching Strategy](#branching-strategy)
- [Commit Standards](#commit-standards)
- [Pull Request Process](#pull-request-process)
- [Code Quality Standards](#code-quality-standards)
- [Testing Requirements](#testing-requirements)
- [Review Process](#review-process)

---

## Code of Conduct

### Our Standards

- **Be respectful**: Treat everyone with respect and professionalism
- **Be collaborative**: Work together to solve problems
- **Be constructive**: Provide helpful feedback
- **Be inclusive**: Welcome diverse perspectives

### Unacceptable Behavior

- Harassment or discriminatory language
- Personal attacks or trolling
- Publishing others' private information
- Other unprofessional conduct

---

## Getting Started

### Prerequisites

- **Go**: 1.21 or later (tested on 1.21, 1.22, 1.23+)
- **GCC**: Required for CGO (SQLite bindings)
- **Git**: Latest stable version
- **golangci-lint**: For code quality checks
- **govulncheck**: For security scanning

### Setting Up Development Environment

```bash
# Clone the repository
git clone https://github.com/edsonmazivila/cli-task-manager.git
cd cli-task-manager

# Configure Git identity
git config user.name "Your Name"
git config user.email "your.email@example.com"

# Install dependencies
go mod download

# Install development tools
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install golang.org/x/vuln/cmd/govulncheck@latest

# Build the project
make build

# Run tests
make test-all
```

### Verifying Setup

```bash
# Verify Go version
go version

# Verify CGO support
go env | grep CGO_ENABLED

# Verify build works
CGO_ENABLED=1 go build -o task ./cmd/task

# Run CI verification locally
./scripts/ci-verify.sh
```

---

## Development Workflow

### 1. Create Feature Branch

**Always** start from `main` and create a feature branch:

```bash
# Update main branch
git checkout main
git pull origin main

# Create feature branch
git checkout -b feature/your-feature-name
```

### 2. Make Changes

- Write code following project standards
- Add tests for new functionality
- Update documentation as needed
- Run tests locally before committing

### 3. Commit Changes

```bash
# Stage changes
git add .

# Commit with conventional commit message
git commit -m "feat: add your feature description"
```

### 4. Push to Remote

```bash
# Push feature branch (NOT main)
git push origin feature/your-feature-name
```

### 5. Create Pull Request

- Go to GitHub repository
- Click "New Pull Request"
- Select your feature branch
- Fill out PR template completely
- Request review from team members

---

## Branching Strategy

### Branch Types

| Branch | Purpose | Base Branch | Merge To |
|--------|---------|-------------|----------|
| `main` | Production-ready code | N/A | N/A |
| `develop` | Integration branch | `main` | `main` |
| `feature/*` | New features | `main` | `main` (via PR) |
| `bugfix/*` | Bug fixes | `main` | `main` (via PR) |
| `hotfix/*` | Urgent fixes | `main` | `main` (via PR) |
| `docs/*` | Documentation | `main` | `main` (via PR) |

### Naming Conventions

**Feature Branches**:
```
feature/add-task-export
feature/rest-api-implementation
feature/web-ui-dashboard
```

**Bug Fix Branches**:
```
bugfix/fix-date-filtering
bugfix/memory-leak-list-operation
bugfix/null-pointer-complete-task
```

**Hotfix Branches**:
```
hotfix/security-vulnerability-CVE-2026-1234
hotfix/critical-data-corruption
hotfix/database-connection-leak
```

**Documentation Branches**:
```
docs/update-api-documentation
docs/add-deployment-guide
docs/improve-contributing-guide
```

### Branch Protection Rules

**`main` branch** (PROTECTED):
- ❌ Direct commits prohibited
- ✅ Requires pull request
- ✅ Requires 1 approval
- ✅ Requires CI checks to pass
- ✅ Requires up-to-date branch
- ❌ Force push prohibited

---

## Commit Standards

### Conventional Commits Format

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Commit Types

| Type | Description | Example |
|------|-------------|---------|
| `feat` | New feature | `feat: add task export functionality` |
| `fix` | Bug fix | `fix: correct date range filtering` |
| `docs` | Documentation | `docs: update README with examples` |
| `style` | Code formatting | `style: format code with gofmt` |
| `refactor` | Code restructuring | `refactor: simplify repository layer` |
| `perf` | Performance improvement | `perf: optimize database queries` |
| `test` | Add/update tests | `test: add integration tests for export` |
| `chore` | Maintenance | `chore: update dependencies` |
| `ci` | CI/CD changes | `ci: add security scanning` |
| `build` | Build system | `build: update Makefile targets` |

### Commit Message Examples

**Good Commits**:

```bash
# Feature addition
git commit -m "feat(cli): add task export command

- Implement JSON and CSV export formats
- Add --format and --output flags
- Add comprehensive tests
- Update documentation with examples"

# Bug fix
git commit -m "fix(repository): correct date range filtering logic

- Fix off-by-one error in date comparison
- Use >= and <= instead of > and <
- Add boundary condition tests
- Update documentation

Fixes #123"

# Documentation
git commit -m "docs: add deployment guide

- Add production deployment instructions
- Document environment configuration
- Include troubleshooting section
- Add security best practices"
```

**Bad Commits** (Avoid These):

```bash
# Too vague
git commit -m "update"
git commit -m "fix bug"
git commit -m "changes"

# Not conventional commits
git commit -m "Added new feature"
git commit -m "Fixed the issue"
git commit -m "WIP"

# Too large
git commit -m "feat: add 10 new features and fix 5 bugs"
```

### Commit Best Practices

✅ **Do**:
- Commit logical units of work
- Write clear, descriptive messages
- Reference issues when applicable
- Keep commits small and focused
- Test before committing
- Use present tense ("add" not "added")

❌ **Don't**:
- Commit broken code
- Commit commented-out code
- Commit TODO markers
- Mix unrelated changes
- Commit sensitive data
- Use vague messages

---

## Pull Request Process

### Before Creating PR

**Checklist**:
- [ ] Code builds successfully
- [ ] All tests pass locally
- [ ] Code follows style guidelines
- [ ] Documentation updated
- [ ] Commit messages follow conventions
- [ ] Branch is up to date with main
- [ ] CI checks will pass

**Commands to Run**:

```bash
# Format code
go fmt ./...

# Run linter
golangci-lint run

# Run all tests
make test-all

# Run CI verification
./scripts/ci-verify.sh

# Update from main
git checkout main
git pull origin main
git checkout feature/your-feature
git rebase main
```

### Creating Pull Request

1. **Push your branch**:
```bash
git push origin feature/your-feature-name
```

2. **Go to GitHub** and click "New Pull Request"

3. **Fill out PR template**:
   - Clear title following conventional commits
   - Detailed description of changes
   - List of changes made
   - Testing performed
   - Screenshots (if UI changes)
   - Related issues

### Pull Request Template

```markdown
## Description
[Clear description of what this PR does]

## Type of Change
- [ ] New feature (non-breaking change)
- [ ] Bug fix (non-breaking change)
- [ ] Breaking change (fix or feature)
- [ ] Documentation update
- [ ] Performance improvement
- [ ] Code refactoring

## Changes Made
- Change 1
- Change 2
- Change 3

## Testing Performed
- [ ] Unit tests added/updated
- [ ] Integration tests added/updated
- [ ] Manual testing completed
- [ ] CI pipeline passing

## Test Results
```
=== Test Results ===
PASS: All tests passing
Coverage: XX.X%
```

## Checklist
- [ ] Code follows project style guidelines
- [ ] Self-review completed
- [ ] Comments added for complex logic
- [ ] Documentation updated
- [ ] No new warnings
- [ ] Tests added/updated
- [ ] All tests passing
- [ ] CI checks passing

## Related Issues
Closes #[issue number]

## Screenshots (if applicable)
[Add screenshots here]

## Additional Notes
[Any additional information]
```

### PR Review Cycle

1. **Create PR** → Automatically assigned to reviewers
2. **CI Checks** → Must pass before review
3. **Code Review** → Reviewer provides feedback
4. **Address Feedback** → Make requested changes
5. **Re-review** → Reviewer approves
6. **Merge** → Maintainer merges to main

---

## Code Quality Standards

### Go Standards

**Formatting**:
```bash
# Format all Go code
go fmt ./...

# Organize imports
goimports -w .
```

**Linting**:
```bash
# Run golangci-lint
golangci-lint run

# Run go vet
go vet ./...
```

**Code Structure**:
- Follow clean architecture principles
- Use interfaces for abstractions
- Keep functions small and focused
- Avoid code duplication
- Write self-documenting code
- Add comments for complex logic

**Error Handling**:
```go
// Good: Proper error handling
func GetTask(id string) (*Task, error) {
    task, err := repo.GetByID(id)
    if err != nil {
        return nil, fmt.Errorf("failed to get task: %w", err)
    }
    return task, nil
}

// Bad: Ignoring errors
func GetTask(id string) *Task {
    task, _ := repo.GetByID(id)
    return task
}
```

**Naming Conventions**:
- Use camelCase for variables: `taskID`, `userName`
- Use PascalCase for exported types: `Task`, `Repository`
- Use descriptive names: `getUserByID` not `get`
- Use `err` for errors: `if err != nil`

### Security Standards

**Never Commit**:
- API keys or secrets
- Database credentials
- Private keys
- Tokens or passwords
- Personal data

**Always**:
- Use parameterized queries
- Validate user input
- Handle errors properly
- Log security events
- Use HTTPS
- Follow OWASP guidelines

### Performance Standards

- Avoid unnecessary allocations
- Use appropriate data structures
- Implement proper indexing
- Cache when appropriate
- Profile before optimizing
- Document performance considerations

---

## Testing Requirements

### Minimum Requirements

- **Coverage**: 75% minimum (currently at 77.8%)
- **All tests must pass**
- **No race conditions**
- **Integration tests for new features**
- **Unit tests for complex logic**

### Running Tests

```bash
# Run all tests
make test-all

# Run integration tests only
make test-integration

# Run with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run with race detection
go test -race ./...

# Run benchmarks
make bench
```

### Writing Tests

**Test File Naming**:
- `filename_test.go` for same package tests
- `integration_test.go` for integration tests

**Test Function Naming**:
```go
func TestTaskCreation(t *testing.T)
func TestTaskFiltering_ByStatus(t *testing.T)
func BenchmarkTaskQuery(b *testing.B)
```

**Test Structure**:
```go
func TestFeatureName(t *testing.T) {
    // Arrange
    env := setupTestEnvironment(t)
    defer env.Cleanup()
    
    // Act
    result, err := env.Service.DoSomething()
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, expected, result)
}
```

**Table-Driven Tests**:
```go
func TestValidation(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr bool
    }{
        {"valid input", "valid", false},
        {"invalid input", "", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := Validate(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("got error %v, want error %v", err, tt.wantErr)
            }
        })
    }
}
```

---

## Review Process

### Reviewer Responsibilities

**Code Review Checklist**:
- [ ] Code follows project standards
- [ ] Tests are adequate and pass
- [ ] Documentation is updated
- [ ] No security issues
- [ ] Performance is acceptable
- [ ] Error handling is proper
- [ ] No breaking changes (or documented)
- [ ] CI checks pass

**Review Guidelines**:
- Review within 24 hours
- Provide constructive feedback
- Ask questions if unclear
- Suggest improvements
- Approve when satisfied

**Types of Feedback**:
- 🚫 **Required**: Must be fixed before merge
- ⚠️ **Suggestion**: Consider changing
- 💡 **Idea**: Optional improvement
- ❓ **Question**: Needs clarification
- ✅ **Approved**: Ready to merge

### Author Responsibilities

**Responding to Feedback**:
- Address all comments
- Ask for clarification if needed
- Make requested changes
- Update PR description if scope changes
- Request re-review when ready

**Merging**:
- Only maintainers can merge
- Requires 1 approval minimum
- All CI checks must pass
- Branch must be up to date

---

## Development Best Practices

### Before Starting Work

1. Check existing issues
2. Discuss approach if uncertain
3. Update your local main branch
4. Create feature branch
5. Review relevant documentation

### During Development

1. Commit frequently
2. Write tests as you go
3. Run tests regularly
4. Follow coding standards
5. Document as needed

### Before Submitting PR

1. Run full test suite
2. Run CI verification locally
3. Update documentation
4. Review your own code
5. Write clear PR description

### After PR Created

1. Monitor CI checks
2. Respond to reviews promptly
3. Make requested changes
4. Keep PR up to date
5. Merge when approved

---

## Getting Help

### Resources

- **Documentation**: [README.md](README.md)
- **Git Workflow**: [GIT_WORKFLOW.md](GIT_WORKFLOW.md)
- **Test Summary**: [TEST_SUMMARY.md](TEST_SUMMARY.md)
- **Project Overview**: [PROJECT_COMPLETION.md](PROJECT_COMPLETION.md)

### Contact

- **GitHub Issues**: For bug reports and feature requests
- **GitHub Discussions**: For questions and discussions
- **Pull Requests**: For code reviews and contributions

### Questions?

If you have questions:
1. Check existing documentation
2. Search closed issues
3. Create a GitHub Discussion
4. Mention @edsonmazivila for clarification

---

## License

By contributing, you agree that your contributions will be licensed under the same license as the project.

---

**Thank you for contributing to CLI Task Manager!**

*Last Updated: 2026-01-29*
