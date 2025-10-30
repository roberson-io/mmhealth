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
func processV1Packet(zipReader *zip.Reader) (model.SupportPacketDiagnostics, model.SupportPacketStats, model.SupportPacketJobList, error) {
	for _, file := range zipReader.File {
		if filepath.Base(file.Name) == "support_packet.yaml" {
			zippedFile, err := file.Open()
			if err != nil {
				return model.SupportPacketDiagnostics{}, model.SupportPacketStats{}, model.SupportPacketJobList{}, err
			}
			defer func() {
				if closeErr := zippedFile.Close(); closeErr != nil {
					HandleError("Error closing support_packet.yaml:", closeErr)
				}
			}()

			v1Packet, err := processPacketFile(zippedFile)
			if err != nil {
				return model.SupportPacketDiagnostics{}, model.SupportPacketStats{}, model.SupportPacketJobList{}, err
			}

			// Convert V1 to V2
			diagnostics := convertV1ToV2Diagnostics(v1Packet)
			stats := convertV1ToV2Stats(v1Packet)
			jobs := convertV1ToV2Jobs(v1Packet)

			return diagnostics, stats, jobs, nil
		}
	}
	return model.SupportPacketDiagnostics{}, model.SupportPacketStats{}, model.SupportPacketJobList{}, fmt.Errorf("support_packet.yaml not found in V1 format packet")
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

func processPacketFile(file io.Reader) (types.SupportPacket, error) {
	var packet types.SupportPacket
	packetBytes, err := io.ReadAll(file)
	if err != nil {
		return types.SupportPacket{}, err
	}
	// Unmarshal the YAML into the struct
	err = yaml.Unmarshal(packetBytes, &packet)
	if err != nil {
		return types.SupportPacket{}, err
	}

	return packet, nil
}

func parseV2Diagnostics(file io.Reader) (model.SupportPacketDiagnostics, error) {
	var diagnostics model.SupportPacketDiagnostics
	diagnosticsBytes, err := io.ReadAll(file)
	if err != nil {
		return model.SupportPacketDiagnostics{}, err
	}
	// Unmarshal the YAML into the struct
	err = yaml.Unmarshal(diagnosticsBytes, &diagnostics)
	if err != nil {
		return model.SupportPacketDiagnostics{}, err
	}

	return diagnostics, nil
}

func parseV2Stats(file io.Reader) (model.SupportPacketStats, error) {
	var stats model.SupportPacketStats
	statsBytes, err := io.ReadAll(file)
	if err != nil {
		return model.SupportPacketStats{}, err
	}
	// Unmarshal the YAML into the struct
	err = yaml.Unmarshal(statsBytes, &stats)
	if err != nil {
		return model.SupportPacketStats{}, err
	}

	return stats, nil
}

func parseV2Jobs(file io.Reader) (model.SupportPacketJobList, error) {
	var jobs model.SupportPacketJobList
	jobsBytes, err := io.ReadAll(file)
	if err != nil {
		return model.SupportPacketJobList{}, err
	}
	// Unmarshal the YAML into the struct
	err = yaml.Unmarshal(jobsBytes, &jobs)
	if err != nil {
		return model.SupportPacketJobList{}, err
	}

	return jobs, nil
}
