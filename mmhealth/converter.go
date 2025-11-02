package mmhealth

import (
	"github.com/coltoneshaw/mmhealth/mmhealth/types"
)

// convertV1ToV2Diagnostics converts a V1 SupportPacketV1 to V2 SupportPacketDiagnosticsV2
func convertV1ToV2Diagnostics(v1 types.SupportPacketV1) types.SupportPacketDiagnosticsV2 {
	diagnostics := types.SupportPacketDiagnosticsV2{
		Version: 1, // V1 packets map to diagnostics version 1
	}

	diagnostics.Server.OS = v1.ServerOS
	diagnostics.Server.Architecture = v1.ServerArchitecture
	diagnostics.Server.Version = v1.ServerVersion
	diagnostics.Server.BuildHash = v1.BuildHash

	diagnostics.Database.Type = v1.DatabaseType
	diagnostics.Database.Version = v1.DatabaseVersion
	diagnostics.Database.SchemaVersion = v1.DatabaseSchemaVersion

	diagnostics.FileStore.Driver = v1.FileDriver
	diagnostics.FileStore.Status = v1.FileStatus

	diagnostics.Cluster.ID = v1.ClusterID

	diagnostics.LDAP.ServerName = v1.LdapVendorName
	diagnostics.LDAP.ServerVersion = v1.LdapVendorVersion

	diagnostics.ElasticSearch.ServerVersion = v1.ElasticServerVersion
	diagnostics.ElasticSearch.ServerPlugins = v1.ElasticServerPlugins

	diagnostics.License.Company = v1.LicenseTo
	diagnostics.License.Users = v1.LicenseSupportedUsers

	return diagnostics
}

// convertV1ToV2Stats converts a V1 SupportPacketV1 to V2 SupportPacketStatsV2
func convertV1ToV2Stats(v1 types.SupportPacketV1) types.SupportPacketStatsV2 {
	stats := types.SupportPacketStatsV2{}

	stats.ActiveUsers = int64(v1.ActiveUsers)
	stats.DailyActiveUsers = int64(v1.DailyActiveUsers)
	stats.MonthlyActiveUsers = int64(v1.MonthlyActiveUsers)
	stats.Posts = int64(v1.TotalPosts)
	stats.Channels = int64(v1.TotalChannels)
	stats.Teams = int64(v1.TotalTeams)

	return stats
}

// convertV1ToV2Jobs converts a V1 SupportPacketV1 to V2 SupportPacketJobListV2
func convertV1ToV2Jobs(v1 types.SupportPacketV1) types.SupportPacketJobListV2 {
	jobs := types.SupportPacketJobListV2{}

	jobs.DataRetentionJobs = v1.DataRetentionJobs
	jobs.MessageExportJobs = v1.MessageExportJobs
	jobs.ElasticPostIndexingJobs = v1.ElasticPostIndexingJobs
	jobs.ElasticPostAggregationJobs = v1.ElasticPostAggregationJobs
	jobs.LDAPSyncJobs = v1.LdapSyncJobs
	jobs.MigrationJobs = v1.MigrationJobs

	return jobs
}
