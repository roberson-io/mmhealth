package types

import "github.com/mattermost/mattermost/server/public/model"

type CheckType string

const (
	Proactive CheckType = "Proactive"
	Health    CheckType = "Health"
	Adoption  CheckType = "Adoption"
)

type CheckGroup string

const (
	ConfigCheckGroup          CheckGroup = "config"
	EnvironmentCheckGroup     CheckGroup = "environment"
	MattermostLogCheckGroup   CheckGroup = "mattermostLog"
	NotificationLogCheckGroup CheckGroup = "notificationLog"
	PacketCheckGroup          CheckGroup = "packet"
)

type CheckStatus string

const (
	StatusFail   CheckStatus = "fail"
	StatusPass   CheckStatus = "pass"
	StatusWarn   CheckStatus = "warn"
	StatusIgnore CheckStatus = "ignore"
	StatusError  CheckStatus = "error"
)

type CheckSeverity string

const (
	Urgent CheckSeverity = "Urgent"
	High   CheckSeverity = "High"
	Medium CheckSeverity = "Medium"
	Low    CheckSeverity = "Low"
)

type Result struct {
	Pass   string `yaml:"pass"`
	Fail   string `yaml:"fail"`
	Ignore string `yaml:"ignore"`
	Error  string `yaml:"error"`
}

type Check struct {
	Name        string        `yaml:"name"`
	Result      Result        `yaml:"result"`
	Description string        `yaml:"description"`
	Severity    CheckSeverity `yaml:"severity"`
	Type        CheckType     `yaml:"type"`
}

type ChecksFile struct {
	Config          map[string]Check `yaml:"config"`
	Packet          map[string]Check `yaml:"packet"`
	Environment     map[string]Check `yaml:"environment"`
	MattermostLog   map[string]Check `yaml:"mattermostLog"`
	NotificationLog map[string]Check `yaml:"notificationLog"`
	Plugins         map[string]Check `yaml:"plugins"`
}

type PacketData struct {
	Logs             []MattermostLogEntry
	NotificationLogs []byte
	Config           model.Config
	Plugins          model.PluginsResponse
	Diagnostics      SupportPacketDiagnosticsV2
	Stats            SupportPacketStatsV2
	Jobs             SupportPacketJobListV2
}

type ConfigFile struct {
	Versions Versions               `yaml:"versions"`
	Plugins  map[string]PluginEntry `yaml:"plugins"`
}

type Versions struct {
	Supported []string `yaml:"supported"`
	ESR       []string `yaml:"esr"`
}

type PluginEntry struct {
	Repo              string `yaml:"repo"`
	Latest            string `yaml:"latest"`
	LatestReleaseDate string `yaml:"release_date"`
	SupportLevel      string `yaml:"support_level"`
}
