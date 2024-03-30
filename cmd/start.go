package commands

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

func StartCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "",
		Run: func(cmd *cobra.Command, args []string) {
			stdout, err := handlesStartCmd(cmd, args)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Fprint(cmd.OutOrStdout(), stdout)
		},
	}
}

func handlesStartCmd(cobraCommand *cobra.Command, args []string) (string, error) {
	engine := Engine{}
	engine.Setup()

	durationInMinutes, err := cobraCommand.Flags().GetInt("duration")
	description, err := cobraCommand.Flags().GetString("description")

	if err != nil {
		Log.ErrorLog.Printf("%s", err)
		return err.Error(), err
	}

	msg, err := engine.StartSession(time.Duration(durationInMinutes)*time.Minute, description)
	return msg, err
}

func init() {
	startCmd := StartCmd()
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().IntP(
		"duration",
		"t",
		5,
		"duration in minutes",
	)

	startCmd.Flags().StringP(
		"description",
		"d",
		"no description",
		"description of this session",
	)

}
