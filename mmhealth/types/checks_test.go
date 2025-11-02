package types

import (
	"testing"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/stretchr/testify/assert"
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
	assert.Equal(t, 1, packetData.Diagnostics.Version)

	// Verify stats
	assert.Equal(t, int64(8500), packetData.Stats.ActiveUsers)
	assert.Equal(t, int64(2850000), packetData.Stats.Posts)

	// Verify jobs
	assert.Len(t, packetData.Jobs.DataRetentionJobs, 1)
	assert.Len(t, packetData.Jobs.LDAPSyncJobs, 1)
	assert.Equal(t, "job1", packetData.Jobs.DataRetentionJobs[0].Id)
}
