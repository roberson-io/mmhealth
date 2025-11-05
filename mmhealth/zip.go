package mmhealth

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"

	"github.com/coltoneshaw/mmhealth/mmhealth/types"
	"github.com/mattermost/mattermost/server/public/model"
	"gopkg.in/yaml.v3"
)

func hasFile(r *zip.Reader, filename string) bool {
	for _, f := range r.File {
		if f.Name == filename {
			return true
		}
	}
	return false
}

// detectPacketFormat determines if a support packet is V1 or V2 format
// based on the presence of key files.
// Returns "V1" if support_packet.yaml exists, "V2" if diagnostics.yaml exists, or "unknown"
func detectPacketFormat(zipReader *zip.Reader) string {
	if hasFile(zipReader, "diagnostics.yaml") {
		return "V2"
	}
	if hasFile(zipReader, "support_packet.yaml") {
		return "V1"
	}
	return "unknown"
}

// processV1Packet finds and processes a V1 support_packet.yaml file, converts it to V2 format
func processV1Packet(zipReader *zip.Reader) (types.SupportPacketDiagnosticsV2, types.SupportPacketStatsV2, types.SupportPacketJobListV2, error) {
	for _, file := range zipReader.File {
		if filepath.Base(file.Name) == "support_packet.yaml" {
			zippedFile, err := file.Open()
			if err != nil {
				return types.SupportPacketDiagnosticsV2{}, types.SupportPacketStatsV2{}, types.SupportPacketJobListV2{}, err
			}
			defer func() {
				if closeErr := zippedFile.Close(); closeErr != nil {
					HandleError("Error closing support_packet.yaml:", closeErr)
				}
			}()

			v1Packet, err := processPacketFile(zippedFile)
			if err != nil {
				return types.SupportPacketDiagnosticsV2{}, types.SupportPacketStatsV2{}, types.SupportPacketJobListV2{}, err
			}

			// Convert V1 to V2
			diagnostics := convertV1ToV2Diagnostics(v1Packet)
			stats := convertV1ToV2Stats(v1Packet)
			jobs := convertV1ToV2Jobs(v1Packet)

			return diagnostics, stats, jobs, nil
		}
	}
	return types.SupportPacketDiagnosticsV2{}, types.SupportPacketStatsV2{}, types.SupportPacketJobListV2{}, fmt.Errorf("support_packet.yaml not found in V1 format packet")
}

func UnzipToMemory(zipReader *zip.Reader) (*types.PacketData, error) {

	fileContents := &types.PacketData{}

	// Detect packet format and handle V1 conversion
	packetFormat := detectPacketFormat(zipReader)

	switch packetFormat {
	case "V1":
		diagnostics, stats, jobs, err := processV1Packet(zipReader)
		if err != nil {
			return nil, err
		}
		fileContents.Diagnostics = diagnostics
		fileContents.Stats = stats
		fileContents.Jobs = jobs
	case "unknown":
		return nil, fmt.Errorf("unknown support packet format: could not find diagnostics.yaml or support_packet.yaml")
	}

	// Process remaining files
	for _, file := range zipReader.File {
		// Open each file in the zip archive
		zippedFile, err := file.Open()
		if err != nil {
			return nil, err
		}

		fmt.Println("Processing file: ", file.Name)

		defer func() {
			if closeErr := zippedFile.Close(); closeErr != nil {
				HandleError("Error closing zipped file:", closeErr)
			}
		}()

		// Use full path for matching to avoid conflicts with plugin files
		switch file.Name {
		case "sanitized_config.json":
			config, err := processConfigFile(zippedFile)
			if err != nil {
				return nil, err
			}
			fileContents.Config = config
		case "plugins.json":
			plugins, err := processPluginFile(zippedFile)
			if err != nil {
				fmt.Println("Error processing plugins: ", err)
				return nil, err
			}
			fileContents.Plugins = plugins
		case "mattermost.log":
			logs, err := processMattermostLog(zippedFile)
			if err != nil {
				return nil, err
			}
			fileContents.Logs = logs
		case "notification.log":
			notifLogs, err := processNotificationLog(zippedFile)
			if err != nil {
				return nil, err
			}
			fileContents.NotificationLogs = notifLogs

		case "support_packet.yaml":
			// V1 format already processed above, skip here
			continue
		case "diagnostics.yaml":
			fileContents.Diagnostics, err = parseV2Diagnostics(zippedFile)
			if err != nil {
				return nil, err
			}
		case "stats.yaml":
			fileContents.Stats, err = parseV2Stats(zippedFile)
			if err != nil {
				return nil, err
			}
		case "jobs.yaml":
			fileContents.Jobs, err = parseV2Jobs(zippedFile)
			if err != nil {
				return nil, err
			}
		case "permissions.yaml":
			fileContents.Permissions, err = parseV2Permissions(zippedFile)
			if err != nil {
				return nil, err
			}
		default:
			fmt.Println("Ignoring file: ", file.Name)
		}

	}
	return fileContents, nil
}

func processConfigFile(file io.Reader) (model.Config, error) {
	var config model.Config
	err := json.NewDecoder(file).Decode(&config)
	if err != nil {
		return model.Config{}, err
	}
	return config, nil
}

func processPluginFile(file io.Reader) (model.PluginsResponse, error) {
	var plugins model.PluginsResponse
	err := json.NewDecoder(file).Decode(&plugins)
	if err != nil {
		return model.PluginsResponse{}, err
	}
	return plugins, nil
}

func processMattermostLog(file io.Reader) ([]types.MattermostLogEntry, error) {

	logs, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	parsedLogs, err := ParseLogs(logs)

	if err != nil {
		return nil, err
	}
	return parsedLogs, nil
}

func processNotificationLog(file io.Reader) ([]byte, error) {
	logs, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return logs, nil
}

func processPacketFile(file io.Reader) (types.SupportPacketV1, error) {
	var packet types.SupportPacketV1
	packetBytes, err := io.ReadAll(file)
	if err != nil {
		return types.SupportPacketV1{}, err
	}
	// Unmarshal the YAML into the struct
	err = yaml.Unmarshal(packetBytes, &packet)
	if err != nil {
		return types.SupportPacketV1{}, err
	}

	return packet, nil
}

func parseV2Diagnostics(file io.Reader) (types.SupportPacketDiagnosticsV2, error) {
	var diagnostics types.SupportPacketDiagnosticsV2
	err := yaml.NewDecoder(file).Decode(&diagnostics)
	if err != nil {
		return types.SupportPacketDiagnosticsV2{}, err
	}
	return diagnostics, nil
}

func parseV2Stats(file io.Reader) (types.SupportPacketStatsV2, error) {
	var stats types.SupportPacketStatsV2
	err := yaml.NewDecoder(file).Decode(&stats)
	if err != nil {
		return types.SupportPacketStatsV2{}, err
	}
	return stats, nil
}

func parseV2Jobs(file io.Reader) (types.SupportPacketJobListV2, error) {
	var jobs types.SupportPacketJobListV2
	err := yaml.NewDecoder(file).Decode(&jobs)
	if err != nil {
		return types.SupportPacketJobListV2{}, err
	}
	return jobs, nil
}

func parseV2Permissions(file io.Reader) (types.SupportPacketPermissionInfoV2, error) {
	var permissions types.SupportPacketPermissionInfoV2
	err := yaml.NewDecoder(file).Decode(&permissions)
	if err != nil {
		return types.SupportPacketPermissionInfoV2{}, err
	}
	return permissions, nil
}
