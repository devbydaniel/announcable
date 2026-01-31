# CI/CD Pipeline Verification Guide

This guide helps you verify that the CI/CD pipeline is working correctly after setup.

## Quick Verification Checklist

After configuring CI/CD (see README.md), verify:

- [ ] Backend CI runs on backend code changes
- [ ] Widget CI runs on widget code changes
- [ ] Both CI workflows pass
- [ ] Docker image publishes after CI success on main
- [ ] Docker image can be pulled and run
- [ ] Status badges show "passing" on README

## Detailed Verification Steps

### 1. Test Backend CI

```bash
# Create test branch
git checkout -b test-backend-ci

# Make a trivial change
echo "// CI test" >> backend/main.go

# Commit and push
git add backend/main.go
git commit -m "test: verify backend CI"
git push origin test-backend-ci
```

Open a PR on GitHub and verify:
- Backend CI workflow runs automatically
- Workflow shows in PR "Checks" tab
- All steps complete successfully (lint, test, build)
- Widget CI does NOT run (path filtering works)

### 2. Test Widget CI

```bash
# Create test branch
git checkout -b test-widget-ci

# Make a trivial change
echo "/* CI test */" >> widget/src/index.ts

# Commit and push
git add widget/src/index.ts
git commit -m "test: verify widget CI"
git push origin test-widget-ci
```

Open a PR on GitHub and verify:
- Widget CI workflow runs automatically
- All steps complete successfully (lint, build with TypeScript)
- Backend CI does NOT run (path filtering works)

### 3. Test CI Failure Handling

```bash
# Create test branch with intentional error
git checkout -b test-ci-failure

# Add linting error to Go code
echo "func badCode {}" >> backend/main.go

# Commit and push
git add backend/main.go
git commit -m "test: verify CI catches errors"
git push origin test-ci-failure
```

Open a PR and verify:
- Backend CI workflow runs
- Workflow fails at linting step
- Error message is clear
- PR shows failed check
- If branch protection enabled: Cannot merge PR

Fix the error and push again to verify CI passes on the fix.

### 4. Test Docker Publishing

After merging a PR to main:

1. Go to Actions tab on GitHub
2. Verify "Docker Build and Publish" workflow runs
3. Verify it only runs AFTER Backend CI and Widget CI complete successfully
4. Check workflow steps:
   - [ ] Docker Buildx setup succeeds
   - [ ] Login to registry succeeds
   - [ ] Multi-platform build completes
   - [ ] Images pushed with correct tags

### 5. Verify Published Image

```bash
# For GHCR (default)
docker pull ghcr.io/OWNER/REPO:latest

# For Docker Hub (if configured)
docker pull announcable/announcable:latest

# Verify image works
docker run --rm -p 3000:3000 \
  -e POSTGRES_HOST=host.docker.internal \
  -e POSTGRES_PORT=5432 \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_NAME=announcable \
  ghcr.io/OWNER/REPO:latest

# Visit http://localhost:3000 to verify app runs
```

### 6. Verify Image Tags

Check your container registry (GHCR or Docker Hub) has:
- `latest` tag (most recent main build)
- `main-<sha>` tag (specific commit)
- `main` tag (branch-based)

## Troubleshooting

### Backend CI fails with "golangci-lint: command not found"

This shouldn't happen - the workflow installs golangci-lint. If it does:
- Check `.github/workflows/backend-ci.yml` has the golangci-lint-action step
- Verify action version is v6 or later

### Widget CI fails with npm errors

- Check `widget/package-lock.json` exists and is committed
- Verify `npm ci` can install dependencies locally
- Check Node version in workflow matches package.json engines field

### Docker publish runs but fails to push

**For GHCR:**
- Go to Settings → Actions → General
- Under "Workflow permissions", select "Read and write permissions"
- Re-run the workflow

**For Docker Hub:**
- Verify `DOCKER_USERNAME` secret is set correctly
- Verify `DOCKER_PASSWORD` is an access token (not your password)
- Verify `DOCKER_REGISTRY` variable is set to `dockerhub`

### Docker publish runs even when CI fails

This indicates the workflow_run trigger is not configured correctly:
- Check `.github/workflows/docker-publish.yml` uses `workflow_run` trigger
- Verify workflow names match exactly: "Backend CI" and "Widget CI"
- Verify job has condition: `if: ${{ github.event.workflow_run.conclusion == 'success' }}`

### Workflows don't trigger at all

- Verify GitHub Actions is enabled: Settings → Actions → General → "Allow all actions"
- Check workflow YAML syntax is valid
- Verify file paths match: `.github/workflows/*.yml`
- Check repository has workflows committed to default branch

## Expected Behavior Summary

| Trigger | Backend CI | Widget CI | Docker Publish |
|---------|------------|-----------|----------------|
| PR with backend changes | ✅ Runs | ❌ Skips | ❌ Skips |
| PR with widget changes | ❌ Skips | ✅ Runs | ❌ Skips |
| PR with both changes | ✅ Runs | ✅ Runs | ❌ Skips |
| Merge to main (both CI pass) | ✅ Runs | ✅ Runs | ✅ Runs after CI |
| Merge to main (CI fails) | ❌ Fails | ❌ Fails | ❌ Does not run |
| Push to develop | ✅ Runs | ✅ Runs | ❌ Skips |
| Direct push to main | ✅ Runs | ✅ Runs | ✅ Runs after CI |

## Success Criteria

The CI/CD pipeline is fully functional when:

- ✅ All workflow files exist and are valid
- ✅ Backend CI runs and passes on backend changes
- ✅ Widget CI runs and passes on widget changes
- ✅ Docker publish runs only after CI succeeds on main
- ✅ Docker images are published with correct tags
- ✅ Published images can be pulled and run
- ✅ Workflows fail fast on errors
- ✅ Path filtering prevents unnecessary workflow runs
- ✅ Status badges show correct status in README

---

*For CI/CD setup instructions, see [README.md](../README.md#cicd)*
