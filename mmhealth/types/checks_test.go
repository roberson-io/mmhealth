package types

import (
	"testing"

	"github.com/mattermost/mattermost/server/public/model"
)

func TestPacketDataV2Structure(t *testing.T) {
	// Test: Verify PacketData can hold V2 structures
	packetData := PacketData{
		Diagnostics: model.SupportPacketDiagnostics{
			Version: 1,
		},
		Stats: model.SupportPacketStats{
			ActiveUsers:        int64(8500),
			DailyActiveUsers:   int64(3200),
			MonthlyActiveUsers: int64(7800),
			Posts:              int64(2850000),
			Channels:           int64(8500),
			Teams:              int64(25),
		},
		Jobs: model.SupportPacketJobList{
			DataRetentionJobs: []*model.Job{
				{Id: "job1", Type: "data_retention", Status: "success"},
			},
			LDAPSyncJobs: []*model.Job{
				{Id: "job2", Type: "ldap_sync", Status: "success"},
			},
		},
	}

	// Verify diagnostics
	if packetData.Diagnostics.Version != 1 {
		t.Errorf("Expected Diagnostics.Version 1, got %d", packetData.Diagnostics.Version)
	}

	// Verify stats
	if packetData.Stats.ActiveUsers != 8500 {
		t.Errorf("Expected Stats.ActiveUsers 8500, got %d", packetData.Stats.ActiveUsers)
	}
	if packetData.Stats.Posts != 2850000 {
		t.Errorf("Expected Stats.Posts 2850000, got %d", packetData.Stats.Posts)
	}

	// Verify jobs
	if len(packetData.Jobs.DataRetentionJobs) != 1 {
		t.Errorf("Expected 1 data retention job, got %d", len(packetData.Jobs.DataRetentionJobs))
	}
	if len(packetData.Jobs.LDAPSyncJobs) != 1 {
		t.Errorf("Expected 1 LDAP sync job, got %d", len(packetData.Jobs.LDAPSyncJobs))
	}
	if packetData.Jobs.DataRetentionJobs[0].Id != "job1" {
		t.Errorf("Expected job ID 'job1', got '%s'", packetData.Jobs.DataRetentionJobs[0].Id)
	}
}
