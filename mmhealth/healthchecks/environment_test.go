package healthchecks

import (
	"testing"

	"github.com/coltoneshaw/mmhealth/mmhealth/types"
)

func TestH006(t *testing.T) {
	p, checkStatus := setupTest(t, "environment")

	p.Config.Versions.Supported = []string{"6.0.x", "6.1.x", "6.2.x"}
	p.Config.Versions.ESR = []string{"5.7.x", "5.8.x"}

	testCases := []struct {
		name           string
		serverVersion  string
		expectedStatus types.CheckStatus
		expectedResult string
	}{
		{
			name:           "h006 - server version if out of date",
			serverVersion:  "5.0.0",
			expectedStatus: Fail,
			expectedResult: "Unsupported version: 5.0.0",
		},
		{
			name:           "h006 - server version in support",
			serverVersion:  "6.1.1",
			expectedStatus: Pass,
			expectedResult: "Supported version: 6.1.1",
		},
		{
			name:           "h006 - server version in support and ESR",
			serverVersion:  "5.8.1",
			expectedStatus: Warn,
			expectedResult: "Supported version: 5.8.1",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p.packet.Diagnostics.Server.Version = tc.serverVersion
			checkStatus(t, p.h006, tc.expectedStatus, tc.expectedResult)
		})
	}
}

func TestH007(t *testing.T) {
	p, checkStatus := setupTest(t, "environment")

	testCases := []struct {
		name           string
		databaseType   string
		expectedStatus types.CheckStatus
		expectedResult string
	}{
		{
			name:           "h007 - Database type is not postgres",
			databaseType:   "mysql",
			expectedStatus: Fail,
			expectedResult: "mysql",
		},
		{
			name:           "h007 - Database type is postgres",
			databaseType:   "postgres",
			expectedStatus: Pass,
			expectedResult: "Postgres",
		},
		{
			name:           "h007 - Database type is maria",
			databaseType:   "maria",
			expectedStatus: Fail,
			expectedResult: "maria",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p.packet.Diagnostics.Database.Type = tc.databaseType
			checkStatus(t, p.h007, tc.expectedStatus, tc.expectedResult)
		})
	}
}

func TestH008(t *testing.T) {
	p, checkStatus := setupTest(t, "environment")

	testCases := []struct {
		name           string
		license        string
		expectedStatus types.CheckStatus
		expectedResult string
	}{
		{
			name:           "h008 - Server os not linux",
			license:        "darwin",
			expectedStatus: Fail,
			expectedResult: "darwin",
		},
		{
			name:           "h008 - server is linux",
			license:        "linux",
			expectedStatus: Pass,
			expectedResult: "Linux",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p.packet.Diagnostics.Server.OS = tc.license
			checkStatus(t, p.h008, tc.expectedStatus, tc.expectedResult)
		})
	}
}

func TestH009(t *testing.T) {
	p, checkStatus := setupTest(t, "environment")

	testCases := []struct {
		name               string
		totalPosts         int
		enableIndexing     bool
		enableSearching    bool
		enableAutoComplete bool
		expectedStatus     types.CheckStatus
		expectedResult     string
	}{
		{
			name:               "h009 - total posts is less than 2.5 million",
			totalPosts:         2000000,
			enableIndexing:     false,
			enableSearching:    false,
			enableAutoComplete: false,
			expectedStatus:     Ignore,
			expectedResult:     "<2.5M posts, No Elasticsearch",
		},
		{
			name:               "h009 - Elasticsearch indexing is not enabled",
			totalPosts:         3000000,
			enableIndexing:     false,
			enableSearching:    false,
			enableAutoComplete: false,
			expectedStatus:     Fail,
			expectedResult:     ">2.5M posts, No Elasticsearch",
		},
		{
			name:               "h009 - total posts is greater than 2.5 million and Elasticsearch indexing is enabled",
			totalPosts:         3000000,
			enableIndexing:     true,
			enableSearching:    true,
			enableAutoComplete: true,
			expectedStatus:     Pass,
			expectedResult:     "Elasticsearch enabled",
		},
		{
			name:               "h009 - bad total post count",
			totalPosts:         -1,
			enableIndexing:     true,
			enableSearching:    true,
			enableAutoComplete: true,
			expectedStatus:     Error,
			expectedResult:     "Failed to get post count",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p.packet.Stats.Posts = int64(tc.totalPosts)
			p.packet.Config.ElasticsearchSettings.EnableIndexing = &tc.enableIndexing
			p.packet.Config.ElasticsearchSettings.EnableSearching = &tc.enableSearching
			p.packet.Config.ElasticsearchSettings.EnableAutocomplete = &tc.enableAutoComplete
			checkStatus(t, p.h009, tc.expectedStatus, tc.expectedResult)
		})
	}
}

func TestH011(t *testing.T) {
	p, checkStatus := setupTest(t, "environment")

	testCases := []struct {
		name                 string
		plugins              []string
		elasticsearchEnabled bool
		expectedStatus       types.CheckStatus
		expectedResult       string
	}{
		{
			name:                 "h011 - analysis-icu plugin is not installed",
			plugins:              []string{},
			elasticsearchEnabled: true,
			expectedStatus:       Fail,
			expectedResult:       "`analysis-icu` not installed",
		},
		{
			name:                 "h011 - analysis-icu plugin is installed",
			plugins:              []string{"analysis-icu"},
			elasticsearchEnabled: true,
			expectedStatus:       Pass,
			expectedResult:       "`analysis-icu` installed",
		},
		{
			name:                 "h011 - Multiple plugins installed",
			plugins:              []string{"plugin1", "analysis-icu", "plugin2"},
			elasticsearchEnabled: true,
			expectedStatus:       Pass,
			expectedResult:       "`analysis-icu` installed",
		},
		{
			name:                 "h011 - wrong plugin installed",
			plugins:              []string{"plugin1"},
			elasticsearchEnabled: true,
			expectedStatus:       Fail,
			expectedResult:       "`analysis-icu` not installed",
		},
		{
			name:                 "h011 - elasticsearch disabled",
			plugins:              []string{"plugin1"},
			elasticsearchEnabled: false,
			expectedStatus:       Ignore,
			expectedResult:       "Elasticsearch disabled",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p.packet.Diagnostics.ElasticSearch.ServerPlugins = tc.plugins
			p.packet.Config.ElasticsearchSettings.EnableIndexing = &tc.elasticsearchEnabled
			checkStatus(t, p.h011, tc.expectedStatus, tc.expectedResult)
		})
	}
}
