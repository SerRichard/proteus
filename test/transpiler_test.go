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

func TestTranspileCommandLineToolInputs(t *testing.T) {

	// TODO Need to unit test the arg input function independently. They output is incorrect here.

	var input = "data/composite-cli/inputs/inp.cwl"
	var inputs_file = "data/composite-cli/inputs/inp-job.yml"
	var output = "data/composite-cli/inputs/inp_argo_output.yaml"

	err := transpiler.ProcessFile(input, inputs_file, "")
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

func TestTranspileCommandLineToolArrayInputs(t *testing.T) {

	var input = "data/composite-cli/array-inputs/array-inputs.cwl"
	var inputs_file = "data/composite-cli/array-inputs/array-inputs-job.yml"
	var output = "data/composite-cli/inputs/inp_argo_output.yaml"

	err := transpiler.ProcessFile(input, inputs_file, "")
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

func TestTranspileCommandLineToolParamRef(t *testing.T) {

	var input = "data/composite-cli/param-ref/tar_param.cwl"
	var inputs_file = "data/composite-cli/param-ref/tar_param_job.yml"
	var output = "data/composite-cli/param-ref/tar_param_argo_output.yaml"

	err := transpiler.ProcessFile(input, inputs_file, "")
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

func TestTranspileCommandLineToolLocations(t *testing.T) {

	var input = "data/composite-cli/param-ref/tar_param.cwl"
	var inputs_file = "data/composite-cli/param-ref/tar_param_job.yml"
	var output = "data/composite-cli/param-ref/tar_param_argo_output.yaml"

	err := transpiler.ProcessFile(input, inputs_file, "")
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
