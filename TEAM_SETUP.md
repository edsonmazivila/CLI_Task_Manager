# Git Workflow Setup Guide

**Project**: CLI Task Manager  
**Repository Owner**: edsonmazivila  
**Repository**: cli-task-manager  
**Date**: 2026-01-29

---

## ✅ Repository Setup Status

### Current Configuration

**Git Identity:**
- Name: `Edson Mazivila`
- Email: `mazivilaedson12@gmail.com`
- Username: `edsonmazivila`

**Branch Structure:**
```
* feature/initial-cli-task-manager (current)
  develop
  main (protected)
```

**Recent Commits:**
```
877fa75 (HEAD -> feature/initial-cli-task-manager) docs: add team collaboration files
13f261c (develop) docs: add comprehensive Git workflow documentation
7c14772 (tag: v1.0.0, main) chore: update .gitignore to exclude test artifacts
66928d7 feat: initial commit - production-ready CLI task manager
```

**Team Files Added:**
- ✅ CONTRIBUTING.md (comprehensive contribution guidelines)
- ✅ CODEOWNERS (code review assignments)
- ✅ Pull Request template
- ✅ GIT_WORKFLOW.md (Git Flow documentation)
- ✅ .gitignore (Go production-ready)

---

## 🚀 Complete Workflow Commands

### Step 1: Verify Current Status

```bash
# Check Git configuration
git config user.name    # Should show: Edson Mazivila
git config user.email   # Should show: mazivilaedson12@gmail.com

# Check current branch
git branch              # Should show: * feature/initial-cli-task-manager

# Check status
git status

# View commit history
git log --oneline --graph --all -n 10
```

### Step 2: Verify Remote Repository (Before Pushing)

**Option A: Remote Already Exists**
```bash
# Check existing remote
git remote -v

# If remote exists, verify URL
# Expected: https://github.com/edsonmazivila/cli-task-manager.git
```

**Option B: Add New Remote**
```bash
# Add remote repository
git remote add origin https://github.com/edsonmazivila/cli-task-manager.git

# Verify remote was added
git remote -v

# Expected output:
# origin  https://github.com/edsonmazivila/cli-task-manager.git (fetch)
# origin  https://github.com/edsonmazivila/cli-task-manager.git (push)
```

### Step 3: Push Feature Branch (NOT main)

```bash
# Ensure you're on feature branch
git checkout feature/initial-cli-task-manager

# Push feature branch with upstream tracking
git push -u origin feature/initial-cli-task-manager

# Verify push
git branch -vv
```

**Expected Output:**
```
* feature/initial-cli-task-manager 877fa75 [origin/feature/initial-cli-task-manager] docs: add team collaboration files
  develop                          13f261c docs: add comprehensive Git workflow documentation
  main                             7c14772 chore: update .gitignore to exclude test artifacts
```

### Step 4: Create Pull Request on GitHub

**DO NOT USE COMMAND LINE FOR THIS**

1. Go to GitHub: `https://github.com/edsonmazivila/cli-task-manager`
2. Click "Compare & pull request" button
3. Fill out PR template:
   - **Base**: `main`
   - **Compare**: `feature/initial-cli-task-manager`
   - **Title**: `feat: production-ready CLI task manager with team collaboration`
   - **Description**: Use the PR template provided

4. Request review from team members
5. Wait for CI checks to pass
6. Wait for approval
7. Maintainer will merge (DO NOT merge directly to main)

### Step 5: After PR is Merged

```bash
# Switch to main branch
git checkout main

# Pull latest changes (includes your merged PR)
git pull origin main

# Delete local feature branch
git branch -d feature/initial-cli-task-manager

# Delete remote feature branch
git push origin --delete feature/initial-cli-task-manager

# Verify branches
git branch -a
```

---

## 📋 Safety Checks

### Before Pushing

Run these commands to ensure repository is clean:

```bash
# 1. Verify no sensitive data
git log --all --full-history -- "**/.env"
git log --all --full-history -- "**/config.yaml"

# 2. Check .gitignore is working
git status

# Should NOT see:
# - *.db files
# - *.out files
# - /tmp/ directory
# - .env files (if created)
# - IDE files (.vscode/, .idea/)

# 3. Verify tests pass
make test-all

# 4. Run CI verification locally
./scripts/ci-verify.sh

# Expected output:
# ✓ ALL CHECKS PASSED
```

### Security Audit

```bash
# Check for accidentally committed secrets
git log --all --full-history -- "**/secret*"
git log --all --full-history -- "**/*key*"
git log --all --full-history -- "**/*password*"

# Run security scan
govulncheck ./...

# Expected output:
# No vulnerabilities found
```

---

## 🔄 Future Development Workflow

### Starting New Feature

```bash
# 1. Update main branch
git checkout main
git pull origin main

# 2. Create feature branch
git checkout -b feature/your-feature-name

# 3. Make changes and commit
git add .
git commit -m "feat: your feature description"

# 4. Push feature branch
git push -u origin feature/your-feature-name

# 5. Create PR on GitHub (as shown in Step 4 above)
```

### Working on Existing Feature

```bash
# 1. Pull latest changes
git pull origin feature/your-feature-name

# 2. Make changes
# ... edit files ...

# 3. Check what changed
git status
git diff

# 4. Commit changes
git add .
git commit -m "feat: continue feature development"

# 5. Push changes
git push origin feature/your-feature-name
```

### Fixing Bugs

```bash
# 1. Create bugfix branch from main
git checkout main
git pull origin main
git checkout -b bugfix/fix-description

# 2. Fix the bug
# ... make fixes ...

# 3. Commit with detailed message
git add .
git commit -m "fix: detailed bug fix description

- Describe what was wrong
- Explain how you fixed it
- Reference issue if applicable

Fixes #123"

# 4. Push bugfix branch
git push -u origin bugfix/fix-description

# 5. Create PR to main
```

### Updating Branch from Main

If your feature branch is behind main:

```bash
# 1. Ensure you have latest main
git checkout main
git pull origin main

# 2. Go back to your feature branch
git checkout feature/your-feature-name

# 3. Rebase onto main (preferred)
git rebase main

# OR merge main into feature (alternative)
git merge main

# 4. Resolve conflicts if any
# ... resolve conflicts ...
git add .
git rebase --continue  # if rebasing
# OR
git commit            # if merging

# 5. Force push if rebased (use with caution)
git push --force-with-lease origin feature/your-feature-name
```

---

## 🛡️ Branch Protection Rules

### `main` Branch (Protected)

**Rules to Configure on GitHub:**

1. **Require pull request reviews**
   - Required approvals: 1
   - Dismiss stale reviews: ✅ Yes
   - Require review from CODEOWNERS: ✅ Yes

2. **Require status checks**
   - Require branches to be up to date: ✅ Yes
   - Required checks:
     - `test` (Go 1.21, 1.22, 1.23)
     - `build`
     - `lint`
     - `security`

3. **Require conversation resolution**
   - ✅ Require all conversations resolved before merge

4. **Additional settings**
   - ❌ Allow force pushes: Never
   - ❌ Allow deletions: Never
   - ✅ Require linear history: Yes
   - ✅ Require signed commits: Optional but recommended

### Setting Up Branch Protection

**GitHub UI Path:**
1. Go to: Settings → Branches
2. Click: "Add branch protection rule"
3. Branch name pattern: `main`
4. Configure rules as above
5. Click: "Create" or "Save changes"

---

## 🔧 Configuration Files Reference

### `.gitignore`

Already configured for Go production projects:
- ✅ Binaries excluded
- ✅ Test artifacts excluded
- ✅ Coverage files excluded
- ✅ Database files excluded
- ✅ IDE files excluded
- ✅ OS files excluded

### `CONTRIBUTING.md`

Comprehensive guide including:
- ✅ Getting started instructions
- ✅ Development workflow
- ✅ Branching strategy
- ✅ Commit standards (Conventional Commits)
- ✅ Pull request process
- ✅ Code quality standards
- ✅ Testing requirements
- ✅ Review process

### `CODEOWNERS`

Defines code ownership:
- ✅ @edsonmazivila owns all code
- ✅ Specific owners for critical paths
- ✅ Automatic review assignment

### `GIT_WORKFLOW.md`

Complete Git Flow documentation:
- ✅ Branch types and naming
- ✅ Feature development workflow
- ✅ Bug fix workflow
- ✅ Hotfix workflow
- ✅ Release process
- ✅ Commit message conventions
- ✅ Best practices

---

## 📊 Repository Status Summary

```
Repository: cli-task-manager
Owner: edsonmazivila
Current Branch: feature/initial-cli-task-manager

Branches:
├── main (protected, tagged v1.0.0)
├── develop (integration branch)
└── feature/initial-cli-task-manager (active)

Commits: 4 total
Files: 27 tracked
Documentation: Complete

Team Files:
✅ CONTRIBUTING.md
✅ CODEOWNERS
✅ Pull Request Template
✅ GIT_WORKFLOW.md
✅ .gitignore

Ready for:
✅ Remote push
✅ Pull request creation
✅ Team collaboration
✅ CI/CD integration
```

---

## 🚨 Important Reminders

### DO ✅

- Always create feature branches from `main`
- Use conventional commit messages
- Write clear PR descriptions
- Run tests before pushing
- Request code reviews
- Keep commits small and logical
- Update documentation
- Follow CONTRIBUTING.md guidelines

### DON'T ❌

- **NEVER** commit directly to `main`
- **NEVER** force push to `main`
- **NEVER** commit sensitive data
- **NEVER** skip tests
- **NEVER** push broken code
- **NEVER** merge without approval
- **NEVER** use vague commit messages
- **NEVER** bypass CI checks

---

## 🆘 Troubleshooting

### Remote Already Exists

```bash
# Check existing remote
git remote -v

# If wrong, remove and re-add
git remote remove origin
git remote add origin https://github.com/edsonmazivila/cli-task-manager.git
```

### Push Rejected

```bash
# Pull latest changes first
git pull origin feature/initial-cli-task-manager

# Resolve conflicts if any
# Then push again
git push origin feature/initial-cli-task-manager
```

### Wrong Branch

```bash
# If you committed to wrong branch
git log  # Find commit SHA

# Switch to correct branch
git checkout feature/correct-branch

# Cherry-pick the commit
git cherry-pick <commit-sha>

# Go back to wrong branch and reset
git checkout wrong-branch
git reset --hard HEAD~1
```

### Undo Last Commit

```bash
# Keep changes, undo commit
git reset --soft HEAD~1

# Discard changes and commit
git reset --hard HEAD~1

# Amend last commit message
git commit --amend -m "New message"
```

---

## 📞 Getting Help

### Resources

- **CONTRIBUTING.md**: Contribution guidelines
- **GIT_WORKFLOW.md**: Git Flow documentation
- **README.md**: Project documentation
- **GitHub Issues**: Bug reports and features
- **GitHub Discussions**: Questions and support

### Contact

- **GitHub**: @edsonmazivila
- **Email**: mazivilaedson12@gmail.com

---

## ✅ Pre-Push Checklist

Before pushing to remote:

- [ ] Git identity configured correctly
- [ ] On feature branch (not main)
- [ ] All tests passing (`make test-all`)
- [ ] CI verification passing (`./scripts/ci-verify.sh`)
- [ ] No sensitive data in commits
- [ ] Commit messages follow conventions
- [ ] Remote repository URL correct
- [ ] Ready to create PR

**If all checked, you're ready to push!**

---

## 🎯 Next Actions

### Immediate (Now)

1. **Verify remote repository exists on GitHub**
   - Go to: https://github.com/edsonmazivila/cli-task-manager
   - If not exists, create it first

2. **Push feature branch**
   ```bash
   git push -u origin feature/initial-cli-task-manager
   ```

3. **Create Pull Request on GitHub**
   - Base: `main`
   - Compare: `feature/initial-cli-task-manager`
   - Use PR template

### After PR Merged

4. **Update local main**
   ```bash
   git checkout main
   git pull origin main
   ```

5. **Clean up feature branch**
   ```bash
   git branch -d feature/initial-cli-task-manager
   git push origin --delete feature/initial-cli-task-manager
   ```

### Ongoing

6. **Follow development workflow**
   - See [Future Development Workflow](#-future-development-workflow)
   - Follow CONTRIBUTING.md guidelines
   - Maintain code quality standards

---

**Repository setup complete! Ready for professional team collaboration.** 🎉

*Last Updated: 2026-01-29*
