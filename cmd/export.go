package commands

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

func ExportCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "export",
		Short: "",
		Run: func(cmd *cobra.Command, args []string) {
			stdout, err := handlesExportCmd(cmd, args)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Fprint(cmd.OutOrStdout(), stdout)
		},
	}
}

func handlesExportCmd(cobraCommand *cobra.Command, args []string) (string, error) {
	engine := Engine{}
	engine.Setup("./sessions.json")

	path, err := cobraCommand.Flags().GetString("path")
	if err != nil {
		return "", err
	}
	err = engine.ExportToMDTable(path, time.Now())
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("exported to %s", path), nil
}

func init() {
	exportCmd := ExportCmd()
	rootCmd.AddCommand(exportCmd)

	exportCmd.Flags().StringP(
		"path",
		"p",
		"sessions.md",
		"path to export sessions to",
	)

}
