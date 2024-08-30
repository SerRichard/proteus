package cwl

import (
	"errors"
	"fmt"
)

func TypeCheckWorkflowInputParameters(inputs WorkflowInputs) error {
	for _, wfin := range inputs {

		allFiles := isAllFiles(wfin.Type)
		allDirectories := isAllDirectories(wfin.Type)

		if wfin.SecondaryFiles != nil {
			if !allFiles {
				return errors.New("File|[]File expected when secondaryFiles is set")
			}
		}
		if wfin.Streamable != nil && !allFiles {
			return errors.New("Streamable only valid when types are of File|[]File")
		}
		if wfin.Format != nil && !allFiles {
			return errors.New("Format only valid when types are of File|[]File")
		}
		if wfin.LoadContents != nil {
			return errors.New("LoadContents only valid when types of File|[]File")
		}
		if wfin.LoadListing != nil && !allDirectories {
			return errors.New("LoadListing only valid when types of Directory|[]Directory")
		}
	}
	return nil
}

func TypeCheckOutputs(outputs WorkflowOutputs) error {
	for _, wfout := range outputs {

		allFiles := isAllFiles(wfout.Type)
		// type check secondary files
		if wfout.SecondaryFiles != nil {
			if !allFiles {
				return errors.New("File|[]File expected when secondaryFiles is set")
			}
		}
		if wfout.Streamable != nil && !allFiles {
			return errors.New("streamable only valid when types are of File|[]File")
		}
		if wfout.Format != nil && !allFiles {
			return errors.New("Format only valid when types are of File|[]File")
		}
	}
	return nil
}

func TypeCheckSteps(steps WorkflowSteps) error {
	for _, step := range steps {

		// DockerRequirements must be set for steps
		var dockerReq bool
		for _, req := range step.Requirements {
			if req.getClass() == "DockerRequirement" {
				dockerReq = true
			}
		}
		if !dockerReq {
			return fmt.Errorf("no DockerRequirement found in step %+v", step.Id)
		}

		// We want to raise an error on scatter arrays greater than 1 where ScatterMethod is not set
		if (len(step.Scatter.Array) > 1) && (step.ScatterMethod == "") {
			return errors.New("ScatterMethod must be set when scatter arrays are greater than 1")
		}
	}
	return nil
}

func TypeCheckRequirements(reqs Requirements) error {
	if len(reqs) != 0 {
		fmt.Printf("You have tried to set global requirements! These are currently ignored!\n")
	}
	return nil // errors.New("Requirements currently must be specified per step")
}

func TypeCheckHints(hints Hints) error {

	// Need to tidy this up a bit more tomorrow
	if len(hints.Array) != 0 {
		fmt.Printf("You have tried to set global hints! These are currently ignored!\n")
	}

	if len(hints.Map) != 0 {
		fmt.Printf("You have tried to set global hints! These are currently ignored!\n")
	}

	return nil
}

func TypeCheckWorkflow(wf *Workflow, inputs map[string]CWLInputEntry) error {

	err := TypeCheckWorkflowInputParameters(wf.Inputs)
	if err != nil {
		return err
	}

	err = TypeCheckOutputs(wf.Outputs)
	if err != nil {
		return err
	}

	// Currently assumes docker requirements are attached directly to the steps
	err = TypeCheckSteps(wf.Steps)
	if err != nil {
		return err
	}

	// If there are not requirements on the previous step, but they do exist here, then do not error!
	err = TypeCheckRequirements(wf.Requirements)
	if err != nil {
		return err
	}

	err = TypeCheckHints(wf.Hints)
	if err != nil {
		return err
	}

	return nil
}
