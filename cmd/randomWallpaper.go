package commands

import (
	"fmt"

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
	return "start pomodoro session", nil
}

func init() {
	startCmd := StartCmd()
	rootCmd.AddCommand(startCmd)
}
