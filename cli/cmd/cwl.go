package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/SerRichard/proteus/pkg/transpiler"
	"github.com/spf13/cobra"
)

// TranspileCommand converts a command-line tool into a different format.
func TranspileCommand() *cobra.Command {

	var inputsFile string
	var locationsFile string

	command := &cobra.Command{
		Use:   "transpile",
		Short: "transpile the provided CWL file into an Argo Workflows resource",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			if len(args) == 0 {
				cmd.HelpFunc()(cmd, args)
				os.Exit(1)
			}

			fmt.Println("You can transpile this file, ", inputsFile, "locations", locationsFile)

			var mainFile = args[0]
			err := transpiler.ProcessFile(mainFile, inputsFile, locationsFile)
			if err != nil {
				log.Fatal(err)
			}

		},
	}

	command.Flags().StringVar(&inputsFile, "inputs", "", "Additional file defining any inputs for the main CWL file.")
	command.Flags().StringVar(&locationsFile, "locations", "", "Additional file defining any loctions for the main CWL file.")

	return command
}
