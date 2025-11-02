package mmhealth

import (
	"testing"

	"github.com/coltoneshaw/mmhealth/mmhealth/types"
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/stretchr/testify/assert"
)

func TestConvertV1ToV2Diagnostics(t *testing.T) {
	// Test: Convert V1 SupportPacketV1 to V2 SupportPacketDiagnostics
	v1Packet := types.SupportPacketV1{
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
	assert.Equal(t, "linux", result.Server.OS)
	assert.Equal(t, "amd64", result.Server.Architecture)
	assert.Equal(t, "10.4.5", result.Server.Version)
	assert.Equal(t, "abc123def456", result.Server.BuildHash)

	// Verify database fields
	assert.Equal(t, "postgres", result.Database.Type)
	assert.Equal(t, "13.22", result.Database.Version)
	assert.Equal(t, "128", result.Database.SchemaVersion)

	// Verify file store fields
	assert.Equal(t, "local", result.FileStore.Driver)
	assert.Equal(t, "OK", result.FileStore.Status)

	// Verify license fields
	assert.Equal(t, "Test Company", result.License.Company)
	assert.Equal(t, 10000, result.License.Users)
}

func TestConvertV1ToV2Stats(t *testing.T) {
	// Test: Convert V1 SupportPacketV1 to V2 SupportPacketStats
	v1Packet := types.SupportPacketV1{
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
	assert.Equal(t, int64(8500), result.ActiveUsers)
	assert.Equal(t, int64(3200), result.DailyActiveUsers)
	assert.Equal(t, int64(7800), result.MonthlyActiveUsers)
	assert.Equal(t, int64(2850000), result.Posts)
	assert.Equal(t, int64(8500), result.Channels)
	assert.Equal(t, int64(25), result.Teams)
}

func TestConvertV1ToV2Jobs(t *testing.T) {
	// Test: Convert V1 SupportPacketV1 to V2 SupportPacketJobList
	v1Packet := types.SupportPacketV1{
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
	assert.Len(t, result.DataRetentionJobs, 1)
	assert.Len(t, result.MessageExportJobs, 1)
	assert.Len(t, result.ElasticPostIndexingJobs, 1)
	assert.Len(t, result.LDAPSyncJobs, 1)
	assert.Len(t, result.MigrationJobs, 1)
}
