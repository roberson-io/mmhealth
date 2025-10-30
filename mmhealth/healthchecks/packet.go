package healthchecks

import (
	"github.com/coltoneshaw/mmhealth/mmhealth/types"
)

func (p *ProcessPacket) packetChecks() (results []CheckResult) {

	checks := map[string]CheckFunc{
		"h012": p.h012,
		"h013": p.h013,
		"h014": p.h014,
		"h015": p.h015,
		"h016": p.h016,
		"h017": p.h017,
	}

	testResults := []CheckResult{}

	for id, check := range checks {
		result := check(p.Checks.Packet)
		result.ID = id
		testResults = append(testResults, result)
	}

	return p.sortResults(testResults)
}

// checks to see if any of the ldap sync jobs have failed and if LDAP is enabled. If so we fail the job.
func (p *ProcessPacket) h012(checks map[string]types.Check) CheckResult {
	// check defaults to pass here because we are looking for the failure message
	check, result := initCheckResult("h012", checks, Pass)

	// check if LDAP is enabled in the config
	if !*p.packet.Config.LdapSettings.Enable {
		result.Result = check.Result.Ignore
		result.Status = Ignore
		return result
	}

	// check if the ldap_sync_jobs for any status that's not success
	for _, job := range p.packet.Jobs.LDAPSyncJobs {
		if job.Status != "success" {
			result.Result = check.Result.Fail
			result.Status = Fail
			return result
		}
	}
	return result
}

// checks to see if any of the message export jobs have failed and if export is enabled. If so we fail the job.
func (p *ProcessPacket) h013(checks map[string]types.Check) CheckResult {
	// check defaults to pass here because we are looking for the failure message
	check, result := initCheckResult("h013", checks, Pass)

	// check if LDAP is enabled in the config
	if !*p.packet.Config.MessageExportSettings.EnableExport {
		result.Result = check.Result.Ignore
		result.Status = Ignore
		return result
	}

	// check if the message_export_jobs for any status that's not success
	for _, job := range p.packet.Jobs.MessageExportJobs {
		if job.Status != "success" {
			result.Result = check.Result.Fail
			result.Status = Fail
			return result
		}
	}
	return result
}

// h014 checks if migration jobs have passed using the support packet.
func (p *ProcessPacket) h014(checks map[string]types.Check) CheckResult {
	// check defaults to pass here because we are looking for the failure message
	check, result := initCheckResult("h014", checks, Pass)

	if len(p.packet.Jobs.MigrationJobs) == 0 {
		result.Result = check.Result.Ignore
		result.Status = Ignore
		return result
	}

	// check if the message_export_jobs for any status that's not success
	for _, job := range p.packet.Jobs.MigrationJobs {
		if job.Status != "success" {
			result.Result = check.Result.Fail
			result.Status = Fail
			return result
		}
	}
	return result
}

// h015 checks if data retention jobs have passed using the support packet.
func (p *ProcessPacket) h015(checks map[string]types.Check) CheckResult {
	// check defaults to pass here because we are looking for the failure message
	check, result := initCheckResult("h015", checks, Pass)

	if !*p.packet.Config.DataRetentionSettings.EnableMessageDeletion &&
		!*p.packet.Config.DataRetentionSettings.EnableFileDeletion {
		result.Result = check.Result.Ignore
		result.Status = Ignore
		return result
	}

	// check if the message_export_jobs for any status that's not success
	for _, job := range p.packet.Jobs.DataRetentionJobs {
		if job.Status != "success" {
			result.Result = check.Result.Fail
			result.Status = Fail
			return result
		}
	}
	return result
}

// h016 checks if data retention jobs have passed using the support packet.
func (p *ProcessPacket) h016(checks map[string]types.Check) CheckResult {
	// check defaults to pass here because we are looking for the failure message
	check, result := initCheckResult("h016", checks, Pass)

	if !*p.packet.Config.ElasticsearchSettings.EnableIndexing {
		result.Result = check.Result.Ignore
		result.Status = Ignore
		return result
	}

	for _, job := range p.packet.Jobs.ElasticPostIndexingJobs {
		if job.Status != "success" {
			result.Result = check.Result.Fail
			result.Status = Fail
			return result
		}
	}
	return result
}

// h017 checks if data retention jobs have passed using the support packet.
func (p *ProcessPacket) h017(checks map[string]types.Check) CheckResult {
	// check defaults to pass here because we are looking for the failure message
	check, result := initCheckResult("h017", checks, Pass)

	if !*p.packet.Config.ElasticsearchSettings.EnableIndexing {
		result.Result = check.Result.Ignore
		result.Status = Ignore
		return result
	}

	for _, job := range p.packet.Jobs.ElasticPostAggregationJobs {
		if job.Status != "success" {
			result.Result = check.Result.Fail
			result.Status = Fail
			return result
		}
	}
	return result
}
