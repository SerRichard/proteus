package cwl

import (
	"errors"
	"fmt"
)

func errorNilRequirements(id *string) error {
	if id != nil {
		return fmt.Errorf("Requirements cannot be nil in %s", *id)
	}
	return errors.New("Requirements cannot be nil")
}

func errorDockerRequirement(id *string) error {
	if id != nil {
		return fmt.Errorf("DockerRequirement must be present in all Argo CWL definitions, %s does not satisfy this", *id)
	}
	return errors.New("DockerRequirement must be present in all Argo CWL definitions")
}

func isAllFiles(tys []CWLType) bool {
	for _, ty := range tys {
		if ty.Kind != CWLFileKind {
			return false
		}
	}
	return true
}

func isAllDirectories(tys []CWLType) bool {
	for _, ty := range tys {
		if ty.Kind != CWLDirectoryKind {
			return false
		}
	}
	return true
}

// TypeCheckCommandlineInputs checks the validity of command-line inputs.
func TypeCheckCommandlineInputs(clins []CommandlineInputParameter) error {
	for _, clin := range clins {

		allFiles := isAllFiles(clin.Type)
		allDirectories := isAllDirectories(clin.Type)
		// type check secondary files
		if clin.SecondaryFiles != nil {
			if !allFiles {
				return errors.New("File|[]File expected when secondaryFiles is set")
			}
		}
		if clin.Streamable != nil && !allFiles {
			return errors.New("Streamable only valid when types are of File|[]File")
		}
		if clin.Format != nil && !allFiles {
			return errors.New("Format only valid when types are of File|[]File")
		}
		if clin.LoadContents != nil {
			return errors.New("LoadContents only valid when types of File|[]File")
		}
		if clin.LoadListing != nil && !allDirectories {
			return errors.New("LoadListing only valid when types of Directory|[]Directory")
		}
	}
	return nil
}

// TypeCheckCommandlineOutputs checks the validity of command-line outputs.
func TypeCheckCommandlineOutputs(clouts []CommandlineOutputParameter) error {
	for _, clout := range clouts {

		allFiles := isAllFiles(clout.Type)
		// type check secondary files
		if clout.SecondaryFiles != nil {
			if !allFiles {
				return errors.New("File|[]File expected when secondaryFiles is set")
			}
		}
		if clout.Streamable != nil && !allFiles {
			return errors.New("streamable only valid when types are of File|[]File")
		}
		if clout.Format != nil && !allFiles {
			return errors.New("Format only valid when types are of File|[]File")
		}
	}
	return nil
}

// TypeCheckCommandlineClass validates the class of command-line tools.
func TypeCheckCommandlineClass(id *string, class string) error {
	if class == "CommandLineTool" {
		return nil
	}
	if id != nil {
		return fmt.Errorf("\"CommandLineTool\" required but %s was provided in %s", class, *id)
	}
	return fmt.Errorf("\"CommandLineTool\" required but %s provided", class)
}

// TypeCheckCommandlineID validates the ID of command-line tools.
func TypeCheckCommandlineID(id *string) error {
	if id == nil {
		return errors.New("\"id\" cannot be nil")
	}
	return nil
}

func typeCheckDockerRequirement(d *DockerRequirement) error {
	if d == nil {
		return errors.New("docker requirement required, nil received")
	}
	if d.DockerPull == nil {
		return errors.New("dockerPull is required")
	}
	return nil
}

// TypeCheckCommandlineRequirements checks the requirements for command-line tools.
func TypeCheckCommandlineRequirements(id *string, clrs []CWLRequirements) error {
	if clrs == nil {
		return errorNilRequirements(id)
	}

	foundDocker := false

	for _, requirement := range clrs {
		if docker, ok := requirement.(DockerRequirement); ok {
			if err := typeCheckDockerRequirement(&docker); err != nil {
				return err
			}
			foundDocker = true
		}
	}

	if !foundDocker {
		return errorDockerRequirement(id)
	}
	return nil
}

// TypeCheckCLICWLVersion checks the CWL version of the CLI.
func TypeCheckCLICWLVersion(id *string, cwlVersion *string) error {
	// allowed to be nil
	if cwlVersion == nil {
		return nil
	}
	if *cwlVersion == CWLVersion {
		return nil
	}
	if id != nil {
		return fmt.Errorf("In %s cwlVerion provided was %s but %s was expected", *id, *cwlVersion, CWLVersion)
	}
	return fmt.Errorf("cwlVersion provided was %s but %s was expected", *cwlVersion, CWLVersion)
}

// TypeCheckBaseCommand validates the base command in a CLI.
func TypeCheckBaseCommand(id *string, baseCommand []string, arguments []string) error {

	if len(baseCommand) > 0 || len(arguments) > 0 {
		return nil
	}
	if id != nil {
		return fmt.Errorf("In %s len(baseCommand) == 0 and len(arguments) was not > 0", *id)
	}
	return errors.New("If len(baseCommand) == 0 then len(arguments) must be > 0")
}

// TypeCheckCommandlineTool checks the overall validity of a command-line tool.
func TypeCheckCommandlineTool(cl *CommandLineTool, inputs map[string]CWLInputEntry) error {
	var err error

	err = TypeCheckCommandlineInputs(cl.Inputs)
	if err != nil {
		return err
	}

	err = TypeCheckCommandlineOutputs(cl.Outputs)
	if err != nil {
		return err
	}

	err = TypeCheckCommandlineClass(cl.ID, cl.Class)
	if err != nil {
		return err
	}
	err = TypeCheckCommandlineID(cl.ID)
	if err != nil {
		return err
	}

	err = TypeCheckCommandlineRequirements(cl.ID, cl.Requirements)
	if err != nil {
		return err
	}

	// err = TypeCheckCommandlineHints(cl.Id, cl.Hints)
	// if err != nil {
	// 	return err
	// }

	err = TypeCheckCLICWLVersion(cl.ID, cl.CWLVersion)
	if err != nil {
		return nil
	}

	err = TypeCheckBaseCommand(cl.ID, cl.BaseCommand, cl.Arguments)
	if err != nil {
		return err
	}

	return nil
}
