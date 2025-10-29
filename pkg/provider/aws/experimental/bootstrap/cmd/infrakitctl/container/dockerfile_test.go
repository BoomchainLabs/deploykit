package container // import "github.com/docker/infrakit/pkg/provider/aws/experimental/bootstrap/cmd/infrakitctl/container"

import (
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const dockerfilePath = "Dockerfile"

// TestDockerfileExists verifies that the Dockerfile exists in the expected location
func TestDockerfileExists(t *testing.T) {
	_, err := os.Stat(dockerfilePath)
	require.NoError(t, err, "Dockerfile should exist at %s", dockerfilePath)
}

// TestDockerfileReadable verifies that the Dockerfile can be read
func TestDockerfileReadable(t *testing.T) {
	content, err := ioutil.ReadFile(dockerfilePath)
	require.NoError(t, err, "Dockerfile should be readable")
	require.NotEmpty(t, content, "Dockerfile should not be empty")
}

// TestDockerfileBaseImage verifies the correct base image is used
func TestDockerfileBaseImage(t *testing.T) {
	content, err := ioutil.ReadFile(dockerfilePath)
	require.NoError(t, err)

	lines := strings.Split(string(content), "\n")
	var fromLine string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "FROM ") {
			fromLine = trimmed
			break
		}
	}

	require.NotEmpty(t, fromLine, "Dockerfile should contain a FROM instruction")
	require.Equal(t, "FROM alpine:3.22.1", fromLine, "Base image should be alpine:3.22.1")
}

// TestDockerfileBaseImageNotOldVersion verifies we're not using the old version
func TestDockerfileBaseImageNotOldVersion(t *testing.T) {
	content, err := ioutil.ReadFile(dockerfilePath)
	require.NoError(t, err)

	contentStr := string(content)
	require.NotContains(t, contentStr, "alpine:3.22.2", "Dockerfile should not contain the old alpine version 3.22.2")
	require.NotContains(t, contentStr, "alpine:3.22.0", "Dockerfile should not contain alpine:3.22.0")
}

// TestDockerfileAlpineVersionFormat verifies the alpine version follows proper format
func TestDockerfileAlpineVersionFormat(t *testing.T) {
	content, err := ioutil.ReadFile(dockerfilePath)
	require.NoError(t, err)

	// Alpine version should be in format X.Y.Z
	alpineVersionRegex := regexp.MustCompile(`FROM alpine:(\d+)\.(\d+)\.(\d+)`)
	matches := alpineVersionRegex.FindStringSubmatch(string(content))
	
	require.NotNil(t, matches, "Alpine version should match format X.Y.Z")
	require.Len(t, matches, 4, "Should have major, minor, and patch versions")
	
	// Verify it's a reasonable version (3.x.x)
	require.Equal(t, "3", matches[1], "Alpine major version should be 3")
	require.Equal(t, "22", matches[2], "Alpine minor version should be 22")
	require.Equal(t, "1", matches[3], "Alpine patch version should be 1")
}

// TestDockerfileMaintainer verifies the MAINTAINER instruction is present
func TestDockerfileMaintainer(t *testing.T) {
	content, err := ioutil.ReadFile(dockerfilePath)
	require.NoError(t, err)

	require.Contains(t, string(content), "MAINTAINER", "Dockerfile should contain MAINTAINER instruction")
	require.Contains(t, string(content), "David Chung <david.chung@docker.com>", 
		"Dockerfile should contain the correct maintainer information")
}

// TestDockerfileCAcertificatesInstalled verifies CA certificates are installed
func TestDockerfileCAcertificatesInstalled(t *testing.T) {
	content, err := ioutil.ReadFile(dockerfilePath)
	require.NoError(t, err)

	contentStr := string(content)
	require.Contains(t, contentStr, "apk add --update ca-certificates", 
		"Dockerfile should install ca-certificates")
	require.Contains(t, contentStr, "rm -Rf /tmp/* /var/lib/cache/apk/*", 
		"Dockerfile should clean up apk cache")
}

// TestDockerfileLib64Setup verifies the lib64 directory and symlink setup
func TestDockerfileLib64Setup(t *testing.T) {
	content, err := ioutil.ReadFile(dockerfilePath)
	require.NoError(t, err)

	contentStr := string(content)
	require.Contains(t, contentStr, "mkdir /lib64", "Dockerfile should create /lib64 directory")
	require.Contains(t, contentStr, "ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2", 
		"Dockerfile should create the correct symlink for Go binary compatibility")
}

// TestDockerfileMatchetectlBinary verifies the matchetectl binary is added
func TestDockerfileMatchetectlBinary(t *testing.T) {
	content, err := ioutil.ReadFile(dockerfilePath)
	require.NoError(t, err)

	contentStr := string(content)
	require.Contains(t, contentStr, "ADD bin/matchetectl /usr/local/bin/", 
		"Dockerfile should add matchetectl binary to /usr/local/bin/")
}

// TestDockerfileEntrypoint verifies the ENTRYPOINT is correctly set
func TestDockerfileEntrypoint(t *testing.T) {
	content, err := ioutil.ReadFile(dockerfilePath)
	require.NoError(t, err)

	contentStr := string(content)
	require.Contains(t, contentStr, `ENTRYPOINT [ "matchetectl" ]`, 
		"Dockerfile should have matchetectl as ENTRYPOINT in exec form")
}

// TestDockerfileEntrypointExecForm verifies ENTRYPOINT uses exec form (JSON array)
func TestDockerfileEntrypointExecForm(t *testing.T) {
	content, err := ioutil.ReadFile(dockerfilePath)
	require.NoError(t, err)

	// Exec form should use JSON array syntax: ["cmd"]
	require.Regexp(t, regexp.MustCompile(`ENTRYPOINT\s+\[\s*"matchetectl"\s*\]`), 
		string(content), "ENTRYPOINT should use exec form (JSON array syntax)")
}

// TestDockerfileInstructionOrder verifies instructions are in the correct order
func TestDockerfileInstructionOrder(t *testing.T) {
	content, err := ioutil.ReadFile(dockerfilePath)
	require.NoError(t, err)

	lines := strings.Split(string(content), "\n")
	instructions := []string{}
	
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Skip empty lines and comments
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		
		// Extract instruction (first word)
		parts := strings.Fields(trimmed)
		if len(parts) > 0 {
			instructions = append(instructions, parts[0])
		}
	}

	require.NotEmpty(t, instructions, "Should have at least one instruction")
	require.Equal(t, "FROM", instructions[0], "First instruction should be FROM")
	
	// Verify MAINTAINER comes after FROM
	maintainerIdx := -1
	for i, inst := range instructions {
		if inst == "MAINTAINER" {
			maintainerIdx = i
			break
		}
	}
	require.True(t, maintainerIdx > 0, "MAINTAINER should come after FROM")
	
	// Verify ENTRYPOINT is last
	require.Equal(t, "ENTRYPOINT", instructions[len(instructions)-1], 
		"ENTRYPOINT should be the last instruction")
}

// TestDockerfileNoHardcodedSecrets verifies no secrets are hardcoded
func TestDockerfileNoHardcodedSecrets(t *testing.T) {
	content, err := ioutil.ReadFile(dockerfilePath)
	require.NoError(t, err)

	contentStr := strings.ToLower(string(content))
	
	// Check for common secret patterns
	secretPatterns := []string{
		"password",
		"secret",
		"token",
		"api_key",
		"apikey",
		"private_key",
		"privatekey",
	}
	
	for _, pattern := range secretPatterns {
		require.NotContains(t, contentStr, pattern, 
			"Dockerfile should not contain hardcoded secrets: %s", pattern)
	}
}

// TestDockerfileLayerOptimization verifies efficient layer usage
func TestDockerfileLayerOptimization(t *testing.T) {
	content, err := ioutil.ReadFile(dockerfilePath)
	require.NoError(t, err)

	// Check that apk add and cache cleanup are in the same RUN command
	lines := strings.Split(string(content), "\n")
	
	apkAddLineIdx := -1
	cleanupLineIdx := -1
	
	for i, line := range lines {
		if strings.Contains(line, "apk add") {
			apkAddLineIdx = i
		}
		if strings.Contains(line, "rm -Rf /tmp/* /var/lib/cache/apk/*") {
			cleanupLineIdx = i
		}
	}
	
	require.True(t, apkAddLineIdx >= 0, "Should have apk add command")
	require.True(t, cleanupLineIdx >= 0, "Should have cleanup command")
	require.Equal(t, apkAddLineIdx, cleanupLineIdx, 
		"apk add and cleanup should be in the same RUN command for layer optimization")
}

// TestDockerfileComments verifies important comments are present
func TestDockerfileComments(t *testing.T) {
	content, err := ioutil.ReadFile(dockerfilePath)
	require.NoError(t, err)

	contentStr := string(content)
	require.Contains(t, contentStr, "# needed in order for go binary to work", 
		"Dockerfile should contain comment explaining lib64 setup")
}

// TestDockerfileNoLatestTag verifies no 'latest' tag is used
func TestDockerfileNoLatestTag(t *testing.T) {
	content, err := ioutil.ReadFile(dockerfilePath)
	require.NoError(t, err)

	contentStr := string(content)
	require.NotContains(t, contentStr, ":latest", 
		"Dockerfile should not use 'latest' tag for reproducible builds")
}

// TestDockerfileValidSyntax verifies basic Dockerfile syntax
func TestDockerfileValidSyntax(t *testing.T) {
	file, err := os.Open(dockerfilePath)
	require.NoError(t, err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0
	
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		
		// Verify known instructions
		validInstructions := []string{
			"FROM", "MAINTAINER", "RUN", "CMD", "LABEL", "EXPOSE", 
			"ENV", "ADD", "COPY", "ENTRYPOINT", "VOLUME", "USER", 
			"WORKDIR", "ARG", "ONBUILD", "STOPSIGNAL", "HEALTHCHECK", "SHELL",
		}
		
		firstWord := strings.Fields(line)[0]
		isValid := false
		for _, inst := range validInstructions {
			if firstWord == inst {
				isValid = true
				break
			}
		}
		
		require.True(t, isValid, "Line %d contains unknown instruction: %s", lineNum, firstWord)
	}
	
	require.NoError(t, scanner.Err(), "Error reading Dockerfile")
}

// TestDockerfileRUNCommandFormat verifies RUN commands follow best practices
func TestDockerfileRUNCommandFormat(t *testing.T) {
	content, err := ioutil.ReadFile(dockerfilePath)
	require.NoError(t, err)

	lines := strings.Split(string(content), "\n")
	
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "RUN ") {
			// RUN commands should not be overly complex (security consideration)
			require.True(t, len(trimmed) < 500, 
				"Line %d: RUN command should not be overly long", i+1)
			
			// If using shell operators, they should be properly formatted
			if strings.Contains(trimmed, "&&") {
				require.Contains(t, trimmed, " && ", 
					"Line %d: && operator should have spaces around it", i+1)
			}
		}
	}
}

// TestDockerfileADDvsCOPY verifies ADD is used appropriately
func TestDockerfileADDvsCOPY(t *testing.T) {
	content, err := ioutil.ReadFile(dockerfilePath)
	require.NoError(t, err)

	contentStr := string(content)
	
	// ADD is used for bin/matchetectl which is fine for local files
	// This test documents the choice
	if strings.Contains(contentStr, "ADD ") {
		// If ADD is used, ensure it's not for URLs (should use RUN + curl instead)
		lines := strings.Split(contentStr, "\n")
		for _, line := range lines {
			if strings.HasPrefix(strings.TrimSpace(line), "ADD ") {
				require.NotContains(t, line, "http://", 
					"ADD should not be used for URLs, use RUN + curl instead")
				require.NotContains(t, line, "https://", 
					"ADD should not be used for URLs, use RUN + curl instead")
			}
		}
	}
}

// TestDockerfileAlpineSpecificBestPractices verifies Alpine-specific best practices
func TestDockerfileAlpineSpecificBestPractices(t *testing.T) {
	content, err := ioutil.ReadFile(dockerfilePath)
	require.NoError(t, err)

	contentStr := string(content)
	
	// When using apk, should use --no-cache or clean cache
	if strings.Contains(contentStr, "apk add") {
		hasCleanup := strings.Contains(contentStr, "rm -Rf") || 
					  strings.Contains(contentStr, "--no-cache")
		require.True(t, hasCleanup, 
			"Alpine apk add should either use --no-cache or clean up cache")
	}
}

// TestDockerfileStructure verifies overall structure and completeness
func TestDockerfileStructure(t *testing.T) {
	content, err := ioutil.ReadFile(dockerfilePath)
	require.NoError(t, err)

	contentStr := string(content)
	
	// Essential instructions should be present
	essentialInstructions := map[string]string{
		"FROM":       "base image declaration",
		"ENTRYPOINT": "container entrypoint",
	}
	
	for instruction, description := range essentialInstructions {
		require.Contains(t, contentStr, instruction, 
			"Dockerfile should contain %s (%s)", instruction, description)
	}
}

// TestDockerfileBinaryPathCorrect verifies binary is added to standard location
func TestDockerfileBinaryPathCorrect(t *testing.T) {
	content, err := ioutil.ReadFile(dockerfilePath)
	require.NoError(t, err)

	contentStr := string(content)
	require.Contains(t, contentStr, "/usr/local/bin/", 
		"Binary should be added to /usr/local/bin/ which is in PATH by default")
}

// TestDockerfileNoRootUser verifies container doesn't explicitly run as root
func TestDockerfileNoRootUser(t *testing.T) {
	content, err := ioutil.ReadFile(dockerfilePath)
	require.NoError(t, err)

	contentStr := string(content)
	// Good practice: no explicit USER root
	require.NotContains(t, contentStr, "USER root", 
		"Should not explicitly set USER to root")
	
	// Note: This Dockerfile doesn't set a non-root user, which is acceptable
	// for tools like this, but we document the check
}

// TestDockerfileLineEndings verifies Unix line endings
func TestDockerfileLineEndings(t *testing.T) {
	content, err := ioutil.ReadFile(dockerfilePath)
	require.NoError(t, err)

	// Check for Windows line endings (CRLF)
	require.NotContains(t, string(content), "\r\n", 
		"Dockerfile should use Unix line endings (LF), not Windows (CRLF)")
}

// TestDockerfileFinalNewline verifies file ends with newline
func TestDockerfileFinalNewline(t *testing.T) {
	content, err := ioutil.ReadFile(dockerfilePath)
	require.NoError(t, err)
	require.NotEmpty(t, content)

	// File should end with newline
	lastByte := content[len(content)-1]
	require.Equal(t, byte('\n'), lastByte, 
		"Dockerfile should end with a newline character")
}

// TestDockerfileRelativeToTestFile verifies we can find Dockerfile from test location
func TestDockerfileRelativeToTestFile(t *testing.T) {
	// Get the directory of this test file
	_, filename, _, ok := runtime.Caller(0)
	require.True(t, ok, "Should be able to get caller information")
	
	dir := filepath.Dir(filename)
	dockerfilePath := filepath.Join(dir, "Dockerfile")
	
	_, err := os.Stat(dockerfilePath)
	require.NoError(t, err, "Dockerfile should exist relative to test file at %s", dockerfilePath)
}

// TestDockerfileAlpineVersionConsistency verifies version is used consistently
func TestDockerfileAlpineVersionConsistency(t *testing.T) {
	content, err := ioutil.ReadFile(dockerfilePath)
	require.NoError(t, err)

	contentStr := string(content)
	
	// Count occurrences of alpine version
	alpineCount := strings.Count(contentStr, "alpine:")
	require.Equal(t, 1, alpineCount, 
		"Alpine version should appear exactly once (in FROM instruction)")
}

// TestDockerfileNoUnusedInstructions verifies no commented out instructions
func TestDockerfileNoUnusedInstructions(t *testing.T) {
	content, err := ioutil.ReadFile(dockerfilePath)
	require.NoError(t, err)

	lines := strings.Split(string(content), "\n")
	
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "#") {
			// Comments should not look like commented-out instructions
			commentContent := strings.TrimPrefix(trimmed, "#")
			commentContent = strings.TrimSpace(commentContent)
			
			instructionsToCheck := []string{"FROM", "RUN", "CMD", "ENTRYPOINT", "ADD", "COPY"}
			for _, inst := range instructionsToCheck {
				if strings.HasPrefix(commentContent, inst+" ") {
					t.Logf("Line %d appears to be a commented-out instruction: %s", i+1, trimmed)
					// This is informational, not a hard failure
				}
			}
		}
	}
}