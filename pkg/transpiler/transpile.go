package transpiler

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/SerRichard/proteus/pkg/cwl"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/pretty"
	"gopkg.in/yaml.v3"
)

func extractFileName(filename string, ext string) (string, error) {
	if len(filename) <= len(ext) {
		return "", fmt.Errorf("filename %s is not greater than only the extension %v", filename, ext)
	}
	name := filename[0 : len(filename)-len(ext)]
	return name, nil
}

func TranspileCommandlineTool(cl cwl.CommandLineTool, inputs map[string]cwl.CWLInputEntry, locations cwl.FileLocations, outputFile string) error {

	log.Infof("TypeCheckCommandlineTool")
	err := cwl.TypeCheckCommandlineTool(&cl, inputs)
	if err != nil {
		return err
	}

	log.Infof("EmitCommandlineTool")
	wf, err := EmitCommandlineTool(&cl, inputs, locations)
	if err != nil {
		return err
	}
	// HACK: yaml Marshalling doesn't marshal correctly
	// therefore we turn the Workflow to map[string]interface and marshal that
	data, err := json.Marshal(wf)
	if err != nil {
		return err
	}

	data = pretty.Pretty(data)

	m := make(map[string]interface{})
	err = json.Unmarshal(data, &m)
	if err != nil {
		return err
	}
	data, err = yaml.Marshal(m)
	if err != nil {
		return err
	}
	return os.WriteFile(outputFile, data, 0644)
}

func TranspileCWLWorkflow(workflow cwl.Workflow, inputs map[string]cwl.CWLInputEntry, locations cwl.FileLocations, outputFile string) error {

	// Check the workflow provided
	err := cwl.TypeCheckWorkflow(&workflow, inputs)
	if err != nil {
		return err
	}

	// Convert the CWL Workflow to Argo Workflows
	wf, err := EmitWorkflow(&workflow, inputs, locations)
	if err != nil {
		return err
	}

	data, err := json.Marshal(wf)
	if err != nil {
		return err
	}

	data = pretty.Pretty(data)

	m := make(map[string]interface{})
	err = json.Unmarshal(data, &m)
	if err != nil {
		return err
	}
	data, err = yaml.Marshal(m)
	if err != nil {
		return err
	}

	return os.WriteFile(outputFile, data, 0644)
}

func ProcessFile(inputFile string, inputsFile string, locationsFile string) error {

	log.Infof("Processing on CWL Version: %s ", cwl.CWLVersion)

	ext := filepath.Ext(inputFile)

	var cwlData map[string]interface{}
	var inputs map[string]cwl.CWLInputEntry
	var fileLocations cwl.FileLocations

	if ext != ".cwl" {
		return fmt.Errorf("invalid file extension %s, only common workflow language (.cwl) files are allowed", ext)
	}

	name, err := extractFileName(inputFile, ext)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	outputFile := fmt.Sprintf("%s_argo_output.yaml", name)

	def, err := os.ReadFile(inputFile)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(def, &cwlData)
	if err != nil {
		return err
	}

	if inputsFile != "" {
		data, err := os.ReadFile(inputsFile)
		if err != nil {
			return err
		}

		err = yaml.Unmarshal(data, &inputs)
		if err != nil {
			return err
		}
	}

	log.Infof("Found CommandLineTool")
	if locationsFile != "" {
		data, err := os.ReadFile(locationsFile)
		if err != nil {
			return err
		}
		err = json.Unmarshal(data, &fileLocations)
		if err != nil {
			return err
		}
	}

	class, ok := cwlData["class"]
	if !ok {
		return errors.New("<class> expected")
	}

	if class == "CommandLineTool" {

		log.Infof("Found CommandLineTool")

		var cliTool cwl.CommandLineTool

		err := yaml.Unmarshal(def, &cliTool)
		if err != nil {
			return err
		}

		log.Infof("About to TypeCheckCommandlineTool")

		return TranspileCommandlineTool(cliTool, inputs, fileLocations, outputFile)
	} else if class == "Workflow" {

		log.Infof("Found Workflow")
		var workflow cwl.Workflow

		err := yaml.Unmarshal(def, &workflow)
		if err != nil {
			return err
		}

		return TranspileCWLWorkflow(workflow, inputs, fileLocations, outputFile)
	}

	return nil
}
