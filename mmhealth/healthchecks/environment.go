package healthchecks

import (
	"fmt"

	"github.com/Masterminds/semver"
	"github.com/coltoneshaw/mmhealth/mmhealth/types"
)

func (p *ProcessPacket) environmentChecks() (results []CheckResult) {

	checks := map[string]CheckFunc{
		"h006": p.h006,
		"h007": p.h007,
		"h008": p.h008,
		"h009": p.h009,
		"h011": p.h011,
	}

	for id, check := range checks {
		result := check(p.Checks.Environment)
		result.ID = id
		results = append(results, result)
	}

	return p.sortResults(results)
}

// Server Version check
func (p *ProcessPacket) h006(checks map[string]types.Check) CheckResult {

	check, result := initCheckResult("h006", checks, Fail)

	result.Result = fmt.Sprintf(check.Result.Fail, p.packet.Diagnostics.Server.Version)

	serverVersion, err := semver.NewVersion(p.packet.Diagnostics.Server.Version)
	if err != nil {
		result.Status = Fail
		result.Result = fmt.Sprintf("Error parsing server version: %s", err)
		return result
	}

	for _, version := range p.Config.Versions.Supported {
		constraint, err := semver.NewConstraint(version)
		if err != nil {
			fmt.Printf("Error parsing version constraint: %s", err)
			return result
		}
		if constraint.Check(serverVersion) {
			result.Result = fmt.Sprintf(check.Result.Pass, p.packet.Diagnostics.Server.Version)
			result.Status = Pass
			return result
		}
	}

	for _, version := range p.Config.Versions.ESR {
		esrConstraint, err := semver.NewConstraint(version)
		if err != nil {
			fmt.Printf("Error parsing version constraint: %s", err)
			return result
		}

		if esrConstraint.Check(serverVersion) {
			result.Result = fmt.Sprintf(check.Result.Pass, p.packet.Diagnostics.Server.Version)
			result.Status = Warn
		}
	}

	return result
}

// Databse type is postgres
func (p *ProcessPacket) h007(checks map[string]types.Check) CheckResult {
	check, result := initCheckResult("h007", checks, Fail)

	result.Result = fmt.Sprintf(check.Result.Fail, p.packet.Diagnostics.Database.Type)

	if p.packet.Diagnostics.Database.Type == "postgres" {
		result.Result = check.Result.Pass
		result.Status = Pass
	}

	return result
}

// Server OS is linux
func (p *ProcessPacket) h008(checks map[string]types.Check) CheckResult {
	check, result := initCheckResult("h008", checks, Fail)

	result.Result = fmt.Sprintf(check.Result.Fail, p.packet.Diagnostics.Server.OS)

	if p.packet.Diagnostics.Server.OS == "linux" {
		result.Result = check.Result.Pass
		result.Status = Pass
	}

	return result
}

// Total posts is greater than 2.5 million and ES is enabled and in use
func (p *ProcessPacket) h009(checks map[string]types.Check) CheckResult {
	check, result := initCheckResult("h009", checks, Fail)

	if p.packet.Stats.Posts == -1 {
		result.Status = Error
		result.Result = check.Result.Error
		p.log("Failed to get total posts for check h009. Usually because MaxUsersForStatistics is set lower than the actual number of users.")
		return result
	}

	if *p.packet.Config.ElasticsearchSettings.EnableIndexing && *p.packet.Config.ElasticsearchSettings.EnableSearching && *p.packet.Config.ElasticsearchSettings.EnableAutocomplete {
		result.Result = check.Result.Pass
		result.Status = Pass
		return result
	}

	if p.packet.Stats.Posts < 2500000 {
		result.Status = Ignore
		result.Result = check.Result.Ignore
		return result
	}

	return result
}

// Total posts is greater than 2.5 million and ES is enabled and in use
func (p *ProcessPacket) h011(checks map[string]types.Check) CheckResult {
	check, result := initCheckResult("h011", checks, Fail)

	if !*p.packet.Config.ElasticsearchSettings.EnableIndexing {
		result.Result = check.Result.Ignore
		result.Status = Ignore
		return result
	}

	if len(p.packet.Diagnostics.ElasticSearch.ServerPlugins) > 0 {
		for _, plugin := range p.packet.Diagnostics.ElasticSearch.ServerPlugins {
			if plugin == "analysis-icu" {
				result.Result = check.Result.Pass
				result.Status = Pass
				return result
			}
		}
	}

	return result
}
