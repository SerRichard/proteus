package cmd

import (
  "github.com/spf13/cobra"
)


func ProteusCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "proteus",
		Short: "Proteus cli tool for transpiling to Argo Workflows",
		Long: `
Proteus currently supports transpiling the following languages to Argo Workflows`,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
		},
	}
		
	command.AddCommand(TranspileCommand())

	return command
}
  