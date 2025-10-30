package cmd

import (
	"log"
	"os"

	mmhealth "github.com/coltoneshaw/mmhealth/mmhealth"
	"github.com/coltoneshaw/mmhealth/mmhealth/types"
	"github.com/gocarina/gocsv"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var ParseCSVCmd = &cobra.Command{
	Use:   "parsecsv",
	Short: "parses the csv healthcheck file and resolves any differences in the checks.yaml",
	RunE:  parseCSVCmdF,
}

func init() {

	AddCmd.Hidden = mmhealth.BuildVersion != "(devel)"

	RootCmd.AddCommand(
		ParseCSVCmd,
	)
}

type Record struct {
	ID                  string              `csv:"ID"`
	Status              string              `csv:"Status"`
	Type                types.CheckType     `csv:"Type"`
	InHealthcheckTool   string              `csv:"In Healthcheck Tool"`
	Severity            types.CheckSeverity `csv:"Severity"`
	CheckGroup          types.CheckGroup    `csv:"Check Group"`
	CustomerReportTitle string              `csv:"Title"`
	Description         string              `csv:"Description"`
	RelatedText         string              `csv:"Related Text"`
	LogMessage          string              `csv:"Log message"`
	PassMessage         string              `csv:"Pass Message"`
	FailMessage         string              `csv:"Fail Message"`
	IgnoreMessage       string              `csv:"Ignore Message"`
	InternalNotes       string              `csv:"Internal Notes"`
}

func parseCSVCmdF(cmd *cobra.Command, args []string) error {

	file, err := os.Open("healthchecks.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			mmhealth.HandleError("Error closing CSV file:", closeErr)
		}
	}()

	// Read the CSV file
	var records []*Record
	if err := gocsv.UnmarshalFile(file, &records); err != nil { // Load clients from file
		log.Fatal(err)
	}

	checksFile := types.ChecksFile{
		Config:          make(map[string]types.Check),
		Environment:     make(map[string]types.Check),
		MattermostLog:   make(map[string]types.Check),
		NotificationLog: make(map[string]types.Check),
		Packet:          make(map[string]types.Check),
	}

	for _, record := range records {
		if record.ID == "" {
			continue
		}
		check := types.Check{
			Name:        record.CustomerReportTitle,
			Description: record.Description,
			Severity:    record.Severity,
			Type:        record.Type,
			Result: types.Result{
				Pass:   record.PassMessage,
				Fail:   record.FailMessage,
				Ignore: record.IgnoreMessage,
			},
		}

		switch record.CheckGroup {
		case types.ConfigCheckGroup:
			checksFile.Config[record.ID] = check
		case types.EnvironmentCheckGroup:
			checksFile.Environment[record.ID] = check
		case types.MattermostLogCheckGroup:
			checksFile.MattermostLog[record.ID] = check
		case types.NotificationLogCheckGroup:
			checksFile.NotificationLog[record.ID] = check
		case types.PacketCheckGroup:
			checksFile.Packet[record.ID] = check
		}
	}

	checksFile.Config = sortGroup(checksFile.Config)
	checksFile.Environment = sortGroup(checksFile.Environment)
	checksFile.MattermostLog = sortGroup(checksFile.MattermostLog)
	checksFile.NotificationLog = sortGroup(checksFile.NotificationLog)
	checksFile.Packet = sortGroup(checksFile.Packet)
	checksFile.Plugins = sortGroup(checksFile.Plugins)

	// Marshal the Config struct back into YAML
	err = storeChecksFile(checksFile)

	if err != nil {
		return errors.Wrap(err, "Failed to write checks file")
	}
	// open the csv file

	// parse the csv into a struct
	return nil
}
