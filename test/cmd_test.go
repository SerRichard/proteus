package testing

import (
    "testing"
    
	"github.com/SerRichard/proteus/cli/cmd"
)

func TestProteusCommand(t *testing.T) {

	name := "proteus"

	var command = cmd.ProteusCommand()

	if name != command.Use {
		t.Fatalf(`Typo in main command, found:  %q, expected:  %v`,command.Use, name )
	}

}