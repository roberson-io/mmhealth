// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package types

import "github.com/mattermost/mattermost/server/public/model"

// SupportPacketV1 represents the V1 support packet format (pre-v10.6.0)
// This struct is copied from github.com/mattermost/mattermost/server/public@v0.0.18
// to maintain compatibility with older support packets after upgrading to v0.1.21+
type SupportPacketV1 struct {
	/* Build information */

	ServerOS           string `yaml:"server_os"`
	ServerArchitecture string `yaml:"server_architecture"`
	ServerVersion      string `yaml:"server_version"`
	BuildHash          string `yaml:"build_hash"`

	/* DB */

	DatabaseType          string `yaml:"database_type"`
	DatabaseVersion       string `yaml:"database_version"`
	DatabaseSchemaVersion string `yaml:"database_schema_version"`
	WebsocketConnections  int    `yaml:"websocket_connections"`
	MasterDbConnections   int    `yaml:"master_db_connections"`
	ReplicaDbConnections  int    `yaml:"read_db_connections"`

	/* Cluster */

	ClusterID string `yaml:"cluster_id"`

	/* File store */

	FileDriver string `yaml:"file_driver"`
	FileStatus string `yaml:"file_status"`

	/* LDAP */

	LdapVendorName    string `yaml:"ldap_vendor_name,omitempty"`
	LdapVendorVersion string `yaml:"ldap_vendor_version,omitempty"`

	/* Elastic Search */

	ElasticServerVersion string   `yaml:"elastic_server_version,omitempty"`
	ElasticServerPlugins []string `yaml:"elastic_server_plugins,omitempty"`

	/* License */

	LicenseTo             string `yaml:"license_to"`
	LicenseSupportedUsers int    `yaml:"license_supported_users"`
	LicenseIsTrial        bool   `yaml:"license_is_trial,omitempty"`

	/* Server stats */

	ActiveUsers        int `yaml:"active_users"`
	DailyActiveUsers   int `yaml:"daily_active_users"`
	MonthlyActiveUsers int `yaml:"monthly_active_users"`
	InactiveUserCount  int `yaml:"inactive_user_count"`
	TotalPosts         int `yaml:"total_posts"`
	TotalChannels      int `yaml:"total_channels"`
	TotalTeams         int `yaml:"total_teams"`

	/* Jobs */

	DataRetentionJobs          []*model.Job `yaml:"data_retention_jobs"`
	MessageExportJobs          []*model.Job `yaml:"message_export_jobs"`
	ElasticPostIndexingJobs    []*model.Job `yaml:"elastic_post_indexing_jobs"`
	ElasticPostAggregationJobs []*model.Job `yaml:"elastic_post_aggregation_jobs"`
	BlevePostIndexingJobs      []*model.Job `yaml:"bleve_post_indexin_jobs"`
	LdapSyncJobs               []*model.Job `yaml:"ldap_sync_jobs"`
	MigrationJobs              []*model.Job `yaml:"migration_jobs"`
}

// V2 support packet type aliases for convenience and clarity
// These types are from github.com/mattermost/mattermost/server/public/model (v10.6.0+)

// SupportPacketDiagnosticsV2 is a type alias for model.SupportPacketDiagnostics
type SupportPacketDiagnosticsV2 = model.SupportPacketDiagnostics

// SupportPacketStatsV2 is a type alias for model.SupportPacketStats
type SupportPacketStatsV2 = model.SupportPacketStats

// SupportPacketJobListV2 is a type alias for model.SupportPacketJobList
type SupportPacketJobListV2 = model.SupportPacketJobList
