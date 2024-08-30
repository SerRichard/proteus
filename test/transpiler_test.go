package testing

import (
	"log"
	"os"
	"testing"

	"github.com/SerRichard/proteus/pkg/transpiler"
)

func TestTranspileCommandLineTool(t *testing.T) {

	var input = "data/hello-cli.cwl"
	var output = "data/hello-cli_argo_output.yaml"

	err := transpiler.ProcessFile(input, "", "")
	if err != nil {
		t.Logf("Error caught %d", err)
		t.Fail()
	}

	if _, err := os.Stat(output); err == nil {
		e := os.Remove(output)
		if e != nil {
			log.Fatal(e)
		}
	}

}

func TestTranspileWorkflow(t *testing.T) {

	var input = "data/hello-workflow.cwl"
	var output = "data/hello-workflow_argo_output.yaml"

	err := transpiler.ProcessFile(input, "", "")
	if err != nil {
		t.Logf("Error caught %d", err)
		t.Fail()
	}

	if _, err := os.Stat(output); err == nil {
		e := os.Remove(output)
		if e != nil {
			log.Fatal(e)
		}
	}

}
