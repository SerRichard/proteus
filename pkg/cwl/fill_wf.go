package cwl

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func (inp *WorkflowInputs) UnmarshalYAML(value *yaml.Node) error {
	var inputs map[string]WorkflowInputParameter

	// Try to unmarshal as an Map of WorkflowInputParameter
	if err := value.Decode(&inputs); err != nil {
		return err
	}

	*inp = inputs

	return nil
}

func (inp *WorkflowStepInput) UnmarshalYAML(value *yaml.Node) error {
	var formatInput WorkflowStepInput

	// In this case we want to try and split the single value into Id and Source
	if value.Value != "" && len(value.Content) == 0 {

		char := "/"
		found := false
		for _, c := range value.Value {
			if string(c) == char {
				found = true
				break
			}
		}

		// Using global as a prefix for step inputs that should reference workflow inputs
		var globalInput string = "global/" + value.Value
		if !found {
			formatInput.Source = &globalInput
		} else {
			formatInput.Source = &value.Value
		}

		*inp = formatInput
		return nil
	} else {
		type alias WorkflowStepInput

		var temp alias
		if err := value.Decode(&temp); err != nil {
			return err
		}

		*inp = WorkflowStepInput(temp)
	}

	return nil
}

func (inp *WorkflowStepInputs) UnmarshalYAML(value *yaml.Node) error {
	var inputs WorkflowStepInputs

	// Try to unmarshal as an WorkflowStepInput
	var oneInput WorkflowStepInput
	if err := value.Decode(&oneInput); err == nil {
		var asArray []WorkflowStepInput
		asArray = append(asArray, oneInput)
		inputs.Array = asArray
		inputs.Map = nil
	}

	// Try to unmarshal as an Array of WorkflowStepInput
	var arrayInputs []WorkflowStepInput
	if err := value.Decode(&arrayInputs); err == nil {
		inputs.Array = arrayInputs
		inputs.Map = nil
	}

	// Try to unmarshal as an Map of WorkflowStepInput
	var mapInputs map[string]WorkflowStepInput
	if err := value.Decode(&mapInputs); err == nil {
		inputs.Array = nil
		inputs.Map = mapInputs
	}

	*inp = inputs

	return nil
}

func (steps *WorkflowSteps) UnmarshalYAML(value *yaml.Node) error {
	var outSteps WorkflowSteps

	// Preserve the order of the nodes
	var keys []string
	for _, node := range value.Content {
		if node.Tag == "!!str" {
			keys = append(keys, node.Value)
		}
	}

	var temp map[string]WorkflowStep
	if err := value.Decode(&temp); err != nil {
		return err
	}

	for _, element := range keys {
		var tmpStep WorkflowStep = temp[string(element)]

		if tmpStep.Id == "" {
			tmpStep.Id = element
		}

		outSteps = append(outSteps, tmpStep)
	}

	*steps = outSteps
	return nil
}

func (out *WorkflowStepOutputs) UnmarshalYAML(value *yaml.Node) error {
	var outputs WorkflowStepOutputs

	// Try to unmarshal as an WorkflowStepInput
	var oneOutput []WorkflowStepOutput
	if err := value.Decode(&oneOutput); err == nil {
		*out = outputs
	}

	// Try to unmarshal as a list of strings
	var arrayOutputs []string
	if err := value.Decode(&arrayOutputs); err == nil {
		for _, thing := range arrayOutputs {
			var tmpOutput WorkflowStepOutput
			tmpOutput.Id = &thing

			outputs = append(outputs, tmpOutput)
		}
	}

	*out = outputs

	return nil
}

func (out *WorkflowOutputs) UnmarshalYAML(value *yaml.Node) error {

	var tmpOutputs map[string]WorkflowOutputParameter
	if err := value.Decode(&tmpOutputs); err != nil {
		return err
	}

	for key, output := range tmpOutputs {
		var tmpKey = &key
		output.Id = tmpKey

		tmpOutputs[key] = output
	}

	*out = tmpOutputs

	return nil

}

func (run *WorkflowCommandLineTool) UnmarshalYAML(value *yaml.Node) error {

	var tmpRun WorkflowCommandLineTool
	var tmpCLT CommandLineTool

	var runString string
	if err := value.Decode(&runString); err == nil {
		if strings.Contains(runString, ".cwl") {
			tmpContent, err := getrunContents(&runString)
			if err != nil {
				return err
			}

			tmpRun = *tmpContent
		} else {
			return fmt.Errorf("expected .cwl file reference, got %s", runString)
		}
	} else {

		if err := value.Decode(&tmpCLT); err != nil {
			return err
		}

		tmpRun.CommandLineTool = tmpCLT
	}

	*run = tmpRun

	return nil
}

func locateFile(run *string) (*string, error) {

	var localRun = localPath + *run
	var localshareRun = localsharePath + *run
	var homeshareRun = fmt.Sprintf(homesharePath, os.Getenv("HOME")) + *run

	if _, err := os.Stat(*run); err == nil {
		return run, nil
	} else if _, err := os.Stat(localRun); err == nil {
		return &localRun, nil
	} else if _, err := os.Stat(localshareRun); err == nil {
		return &localshareRun, nil
	} else if _, err := os.Stat(homeshareRun); err == nil {
		return &homeshareRun, nil
	}

	return nil, fmt.Errorf("could not find the file: %+v", *run)
}

func getrunContents(runFilePath *string) (*WorkflowCommandLineTool, error) {
	// Expecting a file path

	existingPath, err := locateFile(runFilePath)
	if err != nil {
		return nil, err
	}

	def, err := os.ReadFile(*existingPath)
	if err != nil {
		return nil, err
	}

	var tmpData map[string]interface{}
	err = yaml.Unmarshal(def, &tmpData)
	if err != nil {
		return nil, err
	}

	if tmpData["class"] == "CommandLineTool" {
		var cliTool WorkflowCommandLineTool

		err := yaml.Unmarshal(def, &cliTool)
		if err != nil {
			return nil, err
		}

		return &cliTool, nil
	} else {
		return nil, fmt.Errorf("only CommandLineTools are currently supported as step input")
	}
}
