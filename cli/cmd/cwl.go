package cmd

import (
  "fmt"
  "os" 

  "github.com/spf13/cobra"

  "github.com/SerRichard/proteus/pkg/transpiler"
)


func TranspileCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "transpile",
		Short: "transpile the provided CWL file into an Argo Workflows resource",
		Run: func(cmd *cobra.Command, args []string) {
			
			if len(args) == 0 {
				cmd.HelpFunc()(cmd, args)
				os.Exit(1)
			}

			fmt.Println("You can transpile this file")


			transpiler.TranspileFile()

		},
	}
 
	return command
}
  