# Dockerfile Test Suite

## Overview
This test suite provides comprehensive validation for the `Dockerfile` in this directory. The tests were generated to validate the change from `alpine:3.22.2` to `alpine:3.22.1`.

## Test Coverage

### 28 Test Functions Created:

#### File Validation (2 tests)
- **TestDockerfileExists**: Verifies the Dockerfile exists
- **TestDockerfileReadable**: Ensures the file can be read and is not empty

#### Base Image Tests (4 tests)
- **TestDockerfileBaseImage**: Validates the base image is `alpine:3.22.1`
- **TestDockerfileBaseImageNotOldVersion**: Ensures old versions (3.22.2, 3.22.0) are not used
- **TestDockerfileAlpineVersionFormat**: Verifies version follows X.Y.Z format
- **TestDockerfileAlpineVersionConsistency**: Ensures version appears only once

#### Content Validation (7 tests)
- **TestDockerfileMaintainer**: Checks MAINTAINER instruction
- **TestDockerfileCAcertificatesInstalled**: Validates CA certificates installation
- **TestDockerfileLib64Setup**: Verifies lib64 directory and symlink for Go binary
- **TestDockerfileMatchetectlBinary**: Checks matchetectl binary addition
- **TestDockerfileEntrypoint**: Validates ENTRYPOINT configuration
- **TestDockerfileEntrypointExecForm**: Ensures exec form (JSON array) is used
- **TestDockerfileComments**: Verifies important comments are present

#### Structure & Order (2 tests)
- **TestDockerfileInstructionOrder**: Validates instruction ordering (FROM first, ENTRYPOINT last)
- **TestDockerfileStructure**: Checks overall completeness

#### Best Practices (8 tests)
- **TestDockerfileNoHardcodedSecrets**: Ensures no secrets in Dockerfile
- **TestDockerfileLayerOptimization**: Validates efficient layer usage (apk + cleanup in same RUN)
- **TestDockerfileNoLatestTag**: Ensures no 'latest' tags for reproducibility
- **TestDockerfileRUNCommandFormat**: Validates RUN command formatting
- **TestDockerfileADDvsCOPY**: Checks appropriate use of ADD vs COPY
- **TestDockerfileAlpineSpecificBestPractices**: Validates Alpine-specific patterns
- **TestDockerfileBinaryPathCorrect**: Ensures binary in standard PATH location
- **TestDockerfileNoRootUser**: Checks no explicit root user

#### Syntax & Format (5 tests)
- **TestDockerfileValidSyntax**: Validates Dockerfile syntax
- **TestDockerfileLineEndings**: Ensures Unix line endings (LF)
- **TestDockerfileFinalNewline**: Verifies file ends with newline
- **TestDockerfileRelativeToTestFile**: Tests relative path resolution
- **TestDockerfileNoUnusedInstructions**: Checks for commented-out instructions

## Running the Tests

From the repository root:
```bash
# Run all tests in this package
go test github.com/docker/infrakit/pkg/provider/aws/experimental/bootstrap/cmd/infrakitctl/container

# Run with verbose output
go test -v github.com/docker/infrakit/pkg/provider/aws/experimental/bootstrap/cmd/infrakitctl/container

# Run a specific test
go test -run TestDockerfileBaseImage github.com/docker/infrakit/pkg/provider/aws/experimental/bootstrap/cmd/infrakitctl/container

# Run with coverage
go test -cover github.com/docker/infrakit/pkg/provider/aws/experimental/bootstrap/cmd/infrakitctl/container
```

From this directory:
```bash
# Run all tests
go test

# Run with verbose output
go test -v

# Run with race detection
go test -race
```

## Key Changes Validated

The primary change validated by these tests is the Alpine Linux version downgrade:
- **Old version**: `alpine:3.22.2`
- **New version**: `alpine:3.22.1`

The test suite ensures:
1. The new version (3.22.1) is correctly specified
2. The old version (3.22.2) is not present anywhere
3. The version follows proper semantic versioning format
4. All other Dockerfile instructions remain intact and correct

## Dependencies

The test suite uses:
- **Standard library**: `bufio`, `io/ioutil`, `os`, `path/filepath`, `regexp`, `runtime`, `strings`, `testing`
- **Third-party**: `github.com/stretchr/testify/require` (already in project dependencies)

No new dependencies were introduced.

## Test Philosophy

These tests follow several principles:

1. **Comprehensive Coverage**: Test all aspects of the Dockerfile, not just the changed line
2. **Security-Focused**: Check for hardcoded secrets, proper user configuration, layer optimization
3. **Best Practices**: Validate Docker and Alpine-specific best practices
4. **Maintainability**: Clear, descriptive test names that explain what's being tested
5. **Regression Prevention**: Specifically test that old versions are not reintroduced

## Future Maintenance

When updating the Dockerfile:
1. Update relevant test assertions if instructions change
2. Add new tests for any new instructions or patterns
3. Keep security-focused tests even if they seem redundant
4. Update version-specific tests when Alpine version changes

## Integration with CI/CD

These tests can be integrated into CI/CD pipelines:
```bash
# In Makefile or CI script
cd pkg/provider/aws/experimental/bootstrap/cmd/infrakitctl/container && go test -v
```

The tests are designed to be:
- Fast (no Docker builds required)
- Deterministic (no external dependencies)
- Informative (clear error messages on failure)