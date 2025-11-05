package mmhealth

import (
	"archive/zip"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDetectV1Format(t *testing.T) {
	// Test: A V1 packet (v10.4.5) should be detected as V1
	file, err := os.Open("testdata/support_packets/v10.4.5.zip")
	require.NoError(t, err, "Failed to open test file")
	defer func() {
		assert.NoError(t, file.Close(), "Failed to close file")
	}()

	fileInfo, err := file.Stat()
	require.NoError(t, err, "Failed to stat file")

	zipReader, err := zip.NewReader(file, fileInfo.Size())
	require.NoError(t, err, "Failed to create zip reader")

	result := detectPacketFormat(zipReader)
	assert.Equal(t, "V1", result)
}

func TestDetectV2Format(t *testing.T) {
	// Test: A V2 packet (v11.0.4) should be detected as V2
	file, err := os.Open("testdata/support_packets/v11.0.4.zip")
	require.NoError(t, err, "Failed to open test file")
	defer func() {
		assert.NoError(t, file.Close(), "Failed to close file")
	}()

	fileInfo, err := file.Stat()
	require.NoError(t, err, "Failed to stat file")

	zipReader, err := zip.NewReader(file, fileInfo.Size())
	require.NoError(t, err, "Failed to create zip reader")

	result := detectPacketFormat(zipReader)
	assert.Equal(t, "V2", result)
}

func TestDetectUnknownFormat(t *testing.T) {
	// Test: A zip without support_packet.yaml or diagnostics.yaml should be unknown
	// Create a temporary zip file with only config files (no packet files)
	tmpFile, err := os.CreateTemp("", "test_unknown_*.zip")
	require.NoError(t, err, "Failed to create temp file")
	defer func() {
		assert.NoError(t, os.Remove(tmpFile.Name()), "Failed to remove temp file")
	}()
	defer func() {
		assert.NoError(t, tmpFile.Close(), "Failed to close temp file")
	}()

	// Create a minimal zip with no support packet files
	zipWriter := zip.NewWriter(tmpFile)
	_, err = zipWriter.Create("sanitized_config.json")
	require.NoError(t, err, "Failed to create zip entry")
	require.NoError(t, zipWriter.Close(), "Failed to close zip writer")

	// Reopen for reading
	file, err := os.Open(tmpFile.Name())
	require.NoError(t, err, "Failed to reopen temp file")
	defer func() {
		assert.NoError(t, file.Close(), "Failed to close file")
	}()

	fileInfo, err := file.Stat()
	require.NoError(t, err, "Failed to stat file")

	zipReader, err := zip.NewReader(file, fileInfo.Size())
	require.NoError(t, err, "Failed to create zip reader")

	result := detectPacketFormat(zipReader)
	assert.Equal(t, "unknown", result)
}

func TestParseV2Diagnostics(t *testing.T) {
	// Test: Parse a V2 diagnostics.yaml file
	file, err := os.Open("testdata/v2/diagnostics.yaml")
	require.NoError(t, err, "Failed to open test file")
	defer func() {
		assert.NoError(t, file.Close(), "Failed to close file")
	}()

	result, err := parseV2Diagnostics(file)
	require.NoError(t, err)

	// Verify key fields
	assert.Equal(t, 1, result.Version)
	assert.Equal(t, "linux", result.Server.OS)
	assert.Equal(t, "10.9.5", result.Server.Version)
	assert.Equal(t, "postgres", result.Database.Type)
}

func TestParseV2DiagnosticsInvalidYAML(t *testing.T) {
	// Test: Parse invalid YAML should return error
	invalidYAML := "invalid: yaml: content: [unclosed"
	reader := strings.NewReader(invalidYAML)

	_, err := parseV2Diagnostics(reader)
	assert.Error(t, err, "Expected YAML unmarshal error")
}

func TestParseV2Stats(t *testing.T) {
	// Test: Parse a V2 stats.yaml file
	file, err := os.Open("testdata/v2/stats.yaml")
	require.NoError(t, err, "Failed to open test file")
	defer func() {
		assert.NoError(t, file.Close(), "Failed to close file")
	}()

	result, err := parseV2Stats(file)
	require.NoError(t, err)

	// Verify key fields
	assert.Equal(t, int64(15000), result.RegisteredUsers)
	assert.Equal(t, int64(8500), result.ActiveUsers)
	assert.Equal(t, int64(2850000), result.Posts)
	assert.Equal(t, int64(8500), result.Channels)
}

func TestParseV2StatsInvalidYAML(t *testing.T) {
	// Test: Parse invalid YAML should return error
	invalidYAML := "stats: {invalid yaml structure"
	reader := strings.NewReader(invalidYAML)

	_, err := parseV2Stats(reader)
	assert.Error(t, err, "Expected YAML unmarshal error")
}

func TestParseV2Jobs(t *testing.T) {
	// Test: Parse a V2 jobs.yaml file
	file, err := os.Open("testdata/v2/jobs.yaml")
	require.NoError(t, err, "Failed to open test file")
	defer func() {
		assert.NoError(t, file.Close(), "Failed to close file")
	}()

	result, err := parseV2Jobs(file)
	require.NoError(t, err)

	// Verify LDAP sync jobs
	assert.Len(t, result.LDAPSyncJobs, 1)
	if len(result.LDAPSyncJobs) > 0 {
		assert.Equal(t, "success", result.LDAPSyncJobs[0].Status)
	}

	// Verify data retention jobs
	assert.Len(t, result.DataRetentionJobs, 1)

	// Verify message export jobs (should be empty)
	assert.Empty(t, result.MessageExportJobs)

	// Verify elasticsearch indexing jobs
	assert.Len(t, result.ElasticPostIndexingJobs, 1)

	// Verify migration jobs
	assert.Len(t, result.MigrationJobs, 1)
}

func TestParseV2JobsInvalidYAML(t *testing.T) {
	// Test: Parse invalid YAML should return error
	invalidYAML := "jobs:\n  - invalid: [unclosed"
	reader := strings.NewReader(invalidYAML)

	_, err := parseV2Jobs(reader)
	assert.Error(t, err, "Expected YAML unmarshal error")
}

func TestParseV2Permissions(t *testing.T) {
	// Test: Parse a V2 permissions.yaml file
	file, err := os.Open("testdata/v2/permissions.yaml")
	require.NoError(t, err, "Failed to open test file")
	defer func() {
		assert.NoError(t, file.Close(), "Failed to close file")
	}()

	result, err := parseV2Permissions(file)
	require.NoError(t, err)

	// Verify roles are present
	assert.NotEmpty(t, result.Roles, "Expected roles to be present")

	// Find and verify system_admin role exists
	var foundSystemAdmin bool
	for _, role := range result.Roles {
		if role.Name == "system_admin" {
			foundSystemAdmin = true
			assert.NotEmpty(t, role.Permissions, "Expected system_admin to have permissions")
			break
		}
	}
	assert.True(t, foundSystemAdmin, "Expected to find system_admin role")
}

func TestParseV2PermissionsInvalidYAML(t *testing.T) {
	// Test: Parse invalid YAML should return error
	invalidYAML := "roles:\n  - invalid: [unclosed"
	reader := strings.NewReader(invalidYAML)

	_, err := parseV2Permissions(reader)
	assert.Error(t, err, "Expected YAML unmarshal error")
}

// Integration test: Full packet processing for all versions
func TestUnzipToMemoryAllVersions(t *testing.T) {
	testCases := []struct {
		name            string
		filename        string
		expectedVersion string
		expectedDiagVer int // 0 for V1 (converted), 1 for V2 format v1, 2 for V2 format v2
		expectedOS      string
		expectedDBType  string
		expectedSiteURL string
		hasPermissions  bool // permissions.yaml is available
		description     string
	}{
		{
			name:            "v9.10.3 - V1 baseline",
			filename:        "testdata/support_packets/v9.10.3.zip",
			expectedVersion: "9.10.3",
			expectedDiagVer: 0, // V1 packet, no version field after conversion
			expectedOS:      "linux",
			expectedDBType:  "postgres",
			expectedSiteURL: "https://mattermost.example.com",
			hasPermissions:  false,
			description:     "Pre-v10.0.0 baseline V1 format",
		},
		{
			name:            "v10.4.5 - V1 with metadata",
			filename:        "testdata/support_packets/v10.4.5.zip",
			expectedVersion: "10.4.5",
			expectedDiagVer: 0,
			expectedOS:      "linux",
			expectedDBType:  "postgres",
			expectedSiteURL: "https://mattermost.example.com",
			hasPermissions:  false,
			description:     "V1 format with metadata.yaml (v10.0.0-v10.5.x)",
		},
		{
			name:            "v10.9.5 - V2 format with version 1",
			filename:        "testdata/support_packets/v10.9.5.zip",
			expectedVersion: "10.9.5",
			expectedDiagVer: 1,
			expectedOS:      "linux",
			expectedDBType:  "postgres",
			expectedSiteURL: "https://mattermost.example.com",
			hasPermissions:  false,
			description:     "V2 format with diagnostics version 1 (v10.6.0-v10.9.x)",
		},
		{
			name:            "v10.10.3 - V2 with HA support",
			filename:        "testdata/support_packets/v10.10.3.zip",
			expectedVersion: "10.10.3",
			expectedDiagVer: 2,
			expectedOS:      "linux",
			expectedDBType:  "postgres",
			expectedSiteURL: "https://mattermost.example.com",
			hasPermissions:  false,
			description:     "V2 format with HA support and diagnostics version 2 (v10.10.x)",
		},
		{
			name:            "v11.0.4 - V2 with database_schema",
			filename:        "testdata/support_packets/v11.0.4.zip",
			expectedVersion: "11.0.4",
			expectedDiagVer: 2,
			expectedOS:      "linux",
			expectedDBType:  "postgres",
			expectedSiteURL: "https://mattermost.example.com",
			hasPermissions:  true,
			description:     "V2 format with database_schema.yaml and diagnostics version 2 (v10.11.0+)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			file, err := os.Open(tc.filename)
			require.NoError(t, err, "Failed to open test file")
			defer func() {
				assert.NoError(t, file.Close(), "Failed to close file")
			}()

			fileInfo, err := file.Stat()
			require.NoError(t, err, "Failed to stat file")

			zipReader, err := zip.NewReader(file, fileInfo.Size())
			require.NoError(t, err, "Failed to create zip reader")

			result, err := UnzipToMemory(zipReader)
			require.NoError(t, err)

			// Verify Diagnostics version (for V2 packets)
			if tc.expectedDiagVer > 0 {
				assert.Equal(t, tc.expectedDiagVer, result.Diagnostics.Version, "Diagnostics.Version")
			}

			// Verify server version
			assert.Equal(t, tc.expectedVersion, result.Diagnostics.Server.Version, "Diagnostics.Server.Version")

			// Verify OS
			assert.Equal(t, tc.expectedOS, result.Diagnostics.Server.OS, "Diagnostics.Server.OS")

			// Verify database type
			assert.Equal(t, tc.expectedDBType, result.Diagnostics.Database.Type, "Diagnostics.Database.Type")

			// Verify Config was parsed
			require.NotNil(t, result.Config.ServiceSettings.SiteURL, "Config.ServiceSettings.SiteURL should be populated")
			assert.Equal(t, tc.expectedSiteURL, *result.Config.ServiceSettings.SiteURL, "Config.ServiceSettings.SiteURL")

			// Verify Permissions are parsed for versions that have permissions.yaml
			if tc.hasPermissions {
				assert.NotEmpty(t, result.Permissions.Roles, "Expected permissions.yaml to contain roles")
			}
		})
	}
}
