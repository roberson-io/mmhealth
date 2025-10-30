package mmhealth

import (
	"archive/zip"
	"os"
	"strings"
	"testing"
)

func TestDetectV1Format(t *testing.T) {
	// Test: A V1 packet (v10.4.5) should be detected as V1
	file, err := os.Open("testdata/support_packets/v10.4.5.zip")
	if err != nil {
		t.Fatalf("Failed to open test file: %v", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			t.Errorf("Failed to close file: %v", closeErr)
		}
	}()

	fileInfo, err := file.Stat()
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}

	zipReader, err := zip.NewReader(file, fileInfo.Size())
	if err != nil {
		t.Fatalf("Failed to create zip reader: %v", err)
	}

	result := detectPacketFormat(zipReader)
	expected := "V1"
	if result != expected {
		t.Errorf("Expected format %s, got %s", expected, result)
	}
}

func TestDetectV2Format(t *testing.T) {
	// Test: A V2 packet (v11.0.4) should be detected as V2
	file, err := os.Open("testdata/support_packets/v11.0.4.zip")
	if err != nil {
		t.Fatalf("Failed to open test file: %v", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			t.Errorf("Failed to close file: %v", closeErr)
		}
	}()

	fileInfo, err := file.Stat()
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}

	zipReader, err := zip.NewReader(file, fileInfo.Size())
	if err != nil {
		t.Fatalf("Failed to create zip reader: %v", err)
	}

	result := detectPacketFormat(zipReader)
	expected := "V2"
	if result != expected {
		t.Errorf("Expected format %s, got %s", expected, result)
	}
}

func TestDetectUnknownFormat(t *testing.T) {
	// Test: A zip without support_packet.yaml or diagnostics.yaml should be unknown
	// Create a temporary zip file with only config files (no packet files)
	tmpFile, err := os.CreateTemp("", "test_unknown_*.zip")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer func() {
		if removeErr := os.Remove(tmpFile.Name()); removeErr != nil {
			t.Errorf("Failed to remove temp file: %v", removeErr)
		}
	}()
	defer func() {
		if closeErr := tmpFile.Close(); closeErr != nil {
			t.Errorf("Failed to close temp file: %v", closeErr)
		}
	}()

	// Create a minimal zip with no support packet files
	zipWriter := zip.NewWriter(tmpFile)
	_, err = zipWriter.Create("sanitized_config.json")
	if err != nil {
		t.Fatalf("Failed to create zip entry: %v", err)
	}
	if closeErr := zipWriter.Close(); closeErr != nil {
		t.Fatalf("Failed to close zip writer: %v", closeErr)
	}

	// Reopen for reading
	file, err := os.Open(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to reopen temp file: %v", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			t.Errorf("Failed to close file: %v", closeErr)
		}
	}()

	fileInfo, err := file.Stat()
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}

	zipReader, err := zip.NewReader(file, fileInfo.Size())
	if err != nil {
		t.Fatalf("Failed to create zip reader: %v", err)
	}

	result := detectPacketFormat(zipReader)
	expected := "unknown"
	if result != expected {
		t.Errorf("Expected format %s, got %s", expected, result)
	}
}

func TestParseV2Diagnostics(t *testing.T) {
	// Test: Parse a V2 diagnostics.yaml file
	file, err := os.Open("testdata/v2/diagnostics.yaml")
	if err != nil {
		t.Fatalf("Failed to open test file: %v", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			t.Errorf("Failed to close file: %v", closeErr)
		}
	}()

	result, err := parseV2Diagnostics(file)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify key fields
	if result.Version != 1 {
		t.Errorf("Expected Version 1, got %d", result.Version)
	}
	if result.Server.OS != "linux" {
		t.Errorf("Expected OS 'linux', got '%s'", result.Server.OS)
	}
	if result.Server.Version != "10.9.5" {
		t.Errorf("Expected Version '10.9.5', got '%s'", result.Server.Version)
	}
	if result.Database.Type != "postgres" {
		t.Errorf("Expected Database Type 'postgres', got '%s'", result.Database.Type)
	}
}

func TestParseV2DiagnosticsInvalidYAML(t *testing.T) {
	// Test: Parse invalid YAML should return error
	invalidYAML := "invalid: yaml: content: [unclosed"
	reader := strings.NewReader(invalidYAML)

	_, err := parseV2Diagnostics(reader)
	if err == nil {
		t.Error("Expected YAML unmarshal error, got nil")
	}
}

func TestParseV2Stats(t *testing.T) {
	// Test: Parse a V2 stats.yaml file
	file, err := os.Open("testdata/v2/stats.yaml")
	if err != nil {
		t.Fatalf("Failed to open test file: %v", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			t.Errorf("Failed to close file: %v", closeErr)
		}
	}()

	result, err := parseV2Stats(file)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify key fields
	if result.RegisteredUsers != 15000 {
		t.Errorf("Expected RegisteredUsers 15000, got %d", result.RegisteredUsers)
	}
	if result.ActiveUsers != 8500 {
		t.Errorf("Expected ActiveUsers 8500, got %d", result.ActiveUsers)
	}
	if result.Posts != 2850000 {
		t.Errorf("Expected Posts 2850000, got %d", result.Posts)
	}
	if result.Channels != 8500 {
		t.Errorf("Expected Channels 8500, got %d", result.Channels)
	}
}

func TestParseV2StatsInvalidYAML(t *testing.T) {
	// Test: Parse invalid YAML should return error
	invalidYAML := "stats: {invalid yaml structure"
	reader := strings.NewReader(invalidYAML)

	_, err := parseV2Stats(reader)
	if err == nil {
		t.Error("Expected YAML unmarshal error, got nil")
	}
}

func TestParseV2Jobs(t *testing.T) {
	// Test: Parse a V2 jobs.yaml file
	file, err := os.Open("testdata/v2/jobs.yaml")
	if err != nil {
		t.Fatalf("Failed to open test file: %v", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			t.Errorf("Failed to close file: %v", closeErr)
		}
	}()

	result, err := parseV2Jobs(file)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify LDAP sync jobs
	if len(result.LDAPSyncJobs) != 1 {
		t.Errorf("Expected 1 LDAP sync job, got %d", len(result.LDAPSyncJobs))
	}
	if len(result.LDAPSyncJobs) > 0 {
		if result.LDAPSyncJobs[0].Status != "success" {
			t.Errorf("Expected LDAP sync job status 'success', got '%s'", result.LDAPSyncJobs[0].Status)
		}
	}

	// Verify data retention jobs
	if len(result.DataRetentionJobs) != 1 {
		t.Errorf("Expected 1 data retention job, got %d", len(result.DataRetentionJobs))
	}

	// Verify message export jobs (should be empty)
	if len(result.MessageExportJobs) != 0 {
		t.Errorf("Expected 0 message export jobs, got %d", len(result.MessageExportJobs))
	}

	// Verify elasticsearch indexing jobs
	if len(result.ElasticPostIndexingJobs) != 1 {
		t.Errorf("Expected 1 elasticsearch indexing job, got %d", len(result.ElasticPostIndexingJobs))
	}

	// Verify migration jobs
	if len(result.MigrationJobs) != 1 {
		t.Errorf("Expected 1 migration job, got %d", len(result.MigrationJobs))
	}
}

func TestParseV2JobsInvalidYAML(t *testing.T) {
	// Test: Parse invalid YAML should return error
	invalidYAML := "jobs:\n  - invalid: [unclosed"
	reader := strings.NewReader(invalidYAML)

	_, err := parseV2Jobs(reader)
	if err == nil {
		t.Error("Expected YAML unmarshal error, got nil")
	}
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
			description:     "V2 format with database_schema.yaml and diagnostics version 2 (v10.11.0+)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			file, err := os.Open(tc.filename)
			if err != nil {
				t.Fatalf("Failed to open test file: %v", err)
			}
			defer func() {
				if closeErr := file.Close(); closeErr != nil {
					t.Errorf("Failed to close file: %v", closeErr)
				}
			}()

			fileInfo, err := file.Stat()
			if err != nil {
				t.Fatalf("Failed to stat file: %v", err)
			}

			zipReader, err := zip.NewReader(file, fileInfo.Size())
			if err != nil {
				t.Fatalf("Failed to create zip reader: %v", err)
			}

			result, err := UnzipToMemory(zipReader)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			// Verify Diagnostics version (for V2 packets)
			if tc.expectedDiagVer > 0 && result.Diagnostics.Version != tc.expectedDiagVer {
				t.Errorf("Expected Diagnostics.Version %d, got %d", tc.expectedDiagVer, result.Diagnostics.Version)
			}

			// Verify server version
			if result.Diagnostics.Server.Version != tc.expectedVersion {
				t.Errorf("Expected Diagnostics.Server.Version '%s', got '%s'", tc.expectedVersion, result.Diagnostics.Server.Version)
			}

			// Verify OS
			if result.Diagnostics.Server.OS != tc.expectedOS {
				t.Errorf("Expected Diagnostics.Server.OS '%s', got '%s'", tc.expectedOS, result.Diagnostics.Server.OS)
			}

			// Verify database type
			if result.Diagnostics.Database.Type != tc.expectedDBType {
				t.Errorf("Expected Diagnostics.Database.Type '%s', got '%s'", tc.expectedDBType, result.Diagnostics.Database.Type)
			}

			// Verify Config was parsed
			if result.Config.ServiceSettings.SiteURL == nil {
				t.Error("Expected Config.ServiceSettings.SiteURL to be populated")
			} else if *result.Config.ServiceSettings.SiteURL != tc.expectedSiteURL {
				t.Errorf("Expected Config.ServiceSettings.SiteURL '%s', got '%s'", tc.expectedSiteURL, *result.Config.ServiceSettings.SiteURL)
			}

			t.Logf("✓ %s: %s", tc.name, tc.description)
		})
	}
}
