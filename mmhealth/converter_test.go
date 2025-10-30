package mmhealth

import (
	"testing"

	"github.com/coltoneshaw/mmhealth/mmhealth/types"
	"github.com/mattermost/mattermost/server/public/model"
)

func TestConvertV1ToV2Diagnostics(t *testing.T) {
	// Test: Convert V1 SupportPacket to V2 SupportPacketDiagnostics
	v1Packet := types.SupportPacket{
		ServerOS:              "linux",
		ServerArchitecture:    "amd64",
		ServerVersion:         "10.4.5",
		BuildHash:             "abc123def456",
		DatabaseType:          "postgres",
		DatabaseVersion:       "13.22",
		DatabaseSchemaVersion: "128",
		ClusterID:             "cluster123",
		FileDriver:            "local",
		FileStatus:            "OK",
		LdapVendorName:        "OpenLDAP",
		LdapVendorVersion:     "2.4.44",
		ElasticServerVersion:  "7.10.2",
		ElasticServerPlugins:  []string{"analysis-icu"},
		LicenseTo:             "Test Company",
		LicenseSupportedUsers: 10000,
	}

	result := convertV1ToV2Diagnostics(v1Packet)

	// Verify server fields
	if result.Server.OS != "linux" {
		t.Errorf("Expected Server.OS 'linux', got '%s'", result.Server.OS)
	}
	if result.Server.Architecture != "amd64" {
		t.Errorf("Expected Server.Architecture 'amd64', got '%s'", result.Server.Architecture)
	}
	if result.Server.Version != "10.4.5" {
		t.Errorf("Expected Server.Version '10.4.5', got '%s'", result.Server.Version)
	}
	if result.Server.BuildHash != "abc123def456" {
		t.Errorf("Expected Server.BuildHash 'abc123def456', got '%s'", result.Server.BuildHash)
	}

	// Verify database fields
	if result.Database.Type != "postgres" {
		t.Errorf("Expected Database.Type 'postgres', got '%s'", result.Database.Type)
	}
	if result.Database.Version != "13.22" {
		t.Errorf("Expected Database.Version '13.22', got '%s'", result.Database.Version)
	}
	if result.Database.SchemaVersion != "128" {
		t.Errorf("Expected Database.SchemaVersion '128', got '%s'", result.Database.SchemaVersion)
	}

	// Verify file store fields
	if result.FileStore.Driver != "local" {
		t.Errorf("Expected FileStore.Driver 'local', got '%s'", result.FileStore.Driver)
	}
	if result.FileStore.Status != "OK" {
		t.Errorf("Expected FileStore.Status 'OK', got '%s'", result.FileStore.Status)
	}

	// Verify license fields
	if result.License.Company != "Test Company" {
		t.Errorf("Expected License.Company 'Test Company', got '%s'", result.License.Company)
	}
	if result.License.Users != 10000 {
		t.Errorf("Expected License.Users 10000, got %d", result.License.Users)
	}
}

func TestConvertV1ToV2Stats(t *testing.T) {
	// Test: Convert V1 SupportPacket to V2 SupportPacketStats
	v1Packet := types.SupportPacket{
		ActiveUsers:        8500,
		DailyActiveUsers:   3200,
		MonthlyActiveUsers: 7800,
		InactiveUserCount:  2500,
		TotalPosts:         2850000,
		TotalChannels:      8500,
		TotalTeams:         25,
	}

	result := convertV1ToV2Stats(v1Packet)

	// Verify stats fields
	if result.ActiveUsers != 8500 {
		t.Errorf("Expected ActiveUsers 8500, got %d", result.ActiveUsers)
	}
	if result.DailyActiveUsers != 3200 {
		t.Errorf("Expected DailyActiveUsers 3200, got %d", result.DailyActiveUsers)
	}
	if result.MonthlyActiveUsers != 7800 {
		t.Errorf("Expected MonthlyActiveUsers 7800, got %d", result.MonthlyActiveUsers)
	}
	if result.Posts != 2850000 {
		t.Errorf("Expected Posts 2850000, got %d", result.Posts)
	}
	if result.Channels != 8500 {
		t.Errorf("Expected Channels 8500, got %d", result.Channels)
	}
	if result.Teams != 25 {
		t.Errorf("Expected Teams 25, got %d", result.Teams)
	}
}

func TestConvertV1ToV2Jobs(t *testing.T) {
	// Test: Convert V1 SupportPacket to V2 SupportPacketJobList
	v1Packet := types.SupportPacket{
		DataRetentionJobs: []*model.Job{
			{Id: "job1", Type: "data_retention", Status: "success"},
		},
		MessageExportJobs: []*model.Job{
			{Id: "job2", Type: "message_export", Status: "success"},
		},
		ElasticPostIndexingJobs: []*model.Job{
			{Id: "job3", Type: "elasticsearch_post_indexing", Status: "pending"},
		},
		LdapSyncJobs: []*model.Job{
			{Id: "job4", Type: "ldap_sync", Status: "success"},
		},
		MigrationJobs: []*model.Job{
			{Id: "job5", Type: "migrations", Status: "success"},
		},
	}

	result := convertV1ToV2Jobs(v1Packet)

	// Verify job arrays
	if len(result.DataRetentionJobs) != 1 {
		t.Errorf("Expected 1 data retention job, got %d", len(result.DataRetentionJobs))
	}
	if len(result.MessageExportJobs) != 1 {
		t.Errorf("Expected 1 message export job, got %d", len(result.MessageExportJobs))
	}
	if len(result.ElasticPostIndexingJobs) != 1 {
		t.Errorf("Expected 1 elasticsearch indexing job, got %d", len(result.ElasticPostIndexingJobs))
	}
	if len(result.LDAPSyncJobs) != 1 {
		t.Errorf("Expected 1 LDAP sync job, got %d", len(result.LDAPSyncJobs))
	}
	if len(result.MigrationJobs) != 1 {
		t.Errorf("Expected 1 migration job, got %d", len(result.MigrationJobs))
	}
}
