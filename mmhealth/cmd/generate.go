package cmd

import (
	"archive/zip"
	"os"
	"strings"
	"time"

	mmhealth "github.com/coltoneshaw/mmhealth/mmhealth"
	healthchecks "github.com/coltoneshaw/mmhealth/mmhealth/healthchecks"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var ProcessCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate the entire health report from the support packet.",
	Long:  "Generates the entire health report from the support packet, and outputting a pdf file.",
	RunE:  generateCmdF,
}

func init() {
	ProcessCmd.Flags().StringP("packet", "p", "", "the support packet file to process")
	ProcessCmd.Flags().StringP("outputName", "o", "healthcheck-report", "the output file name for the PDF.")
	ProcessCmd.Flags().StringP("company", "c", "", "The company name for the final report")

	ProcessCmd.Flags().Bool("raw", false, "Skips the generation of a pdf file. ")

	ProcessCmd.Flags().Bool("debug", true, "Whether to show debug logs or not")

	if err := ProcessCmd.MarkFlagRequired("packet"); err != nil {
		panic(err)
	}

	RootCmd.AddCommand(
		ProcessCmd,
	)
}

func generateCmdF(cmd *cobra.Command, args []string) error {
	supportPacketFile, err := cmd.Flags().GetString("packet")
	if err != nil {
		return errors.Wrap(err, "failed to get packet flag")
	}

	outputFileName, err := cmd.Flags().GetString("outputName")
	if err != nil {
		return errors.Wrap(err, "failed to get output file name")
	}

	rawReport, _ := cmd.Flags().GetBool("raw")

	packetReader, err := os.Open(supportPacketFile)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := packetReader.Close(); closeErr != nil {
			mmhealth.HandleError("Error closing packet file:", closeErr)
		}
	}()

	zipFileInfo, err := packetReader.Stat()
	if err != nil {
		return err
	}

	zipReader, err := zip.NewReader(packetReader, zipFileInfo.Size())
	if err != nil || zipReader.File == nil {
		return err
	}

	packetContents, err := mmhealth.UnzipToMemory(zipReader)

	if err != nil {
		return err
	}

	hc := healthchecks.ProcessPacket{}

	report, err := hc.ProcessPacket(*packetContents)
	if err != nil {
		return err
	}

	report.Metadata.CompanyName = packetContents.Diagnostics.License.Company

	// override the company name if it was provided
	if companyName, err := cmd.Flags().GetString("company"); err == nil && companyName != "" {
		report.Metadata.CompanyName = companyName
	}
	report.Metadata.Date = time.Now().Format("Jan 2, 2006")

	err = saveMarkdownReportToFile(outputFileName, report)
	if err != nil {
		return err
	}

	if !rawReport {
		err = mmhealth.ReportToPDF(outputFileName+".md", outputFileName+".pdf")
		if err != nil {
			return err
		}

		err = deleteMarkdownReportFile(outputFileName + ".md")
		if err != nil {
			return err
		}
	}

	return nil
}

func saveMarkdownReportToFile(outputFileName string, results healthchecks.CheckResults) error {
	// Convert the results to YAML
	data, err := yaml.Marshal(results)
	if err != nil {
		return errors.Wrap(err, "failed to marshal results to yaml")
	}

	file, err := os.Create(outputFileName + ".md")
	if err != nil {
		return errors.Wrap(err, "failed to create report markdown")
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			mmhealth.HandleError("Error closing report file:", closeErr)
		}
	}()

	markdown := "---\n" + string(data) + "---\n"

	// special character that causes strange formatting in latex
	markdown = strings.ReplaceAll(markdown, ">", "\\textgreater{} ")
	markdown = strings.ReplaceAll(markdown, "<", "\\textless{} ")

	err = os.WriteFile(outputFileName+".md", []byte(markdown), 0644)
	if err != nil {
		return errors.Wrap(err, "failed to write to report markdown")
	}
	return nil
}

func deleteMarkdownReportFile(reportFilePath string) error {
	err := os.Remove(reportFilePath)
	if err != nil {
		return errors.Wrap(err, "failed to delete report markdown")
	}
	return nil
}
