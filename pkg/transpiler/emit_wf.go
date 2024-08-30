package transpiler

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	apiv1 "k8s.io/api/core/v1"

	"github.com/SerRichard/proteus/pkg/cwl"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
)

func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz"

	// Seed the random number generator
	rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func EmitWorkflowArguments(inputs *cwl.WorkflowInputs) (*v1alpha1.Arguments, error) {

	var args v1alpha1.Arguments

	for key, input := range *inputs {
		var tmpParam v1alpha1.Parameter
		tmpParam.Name = key

		for _, _type := range input.Type {
			switch _type.Kind {
			case cwl.CWLStringKind:
				tmpParam.Value = (*v1alpha1.AnyString)(input.Default)
			case cwl.CWLFileKind:
				tmpParam.Value = (*v1alpha1.AnyString)(input.Default)
			default:
				return nil, fmt.Errorf("%T currently unsupported type", _type.Kind)
			}
		}
		args.Parameters = append(args.Parameters, tmpParam)
	}

	return &args, nil
}

func EmitStepInput(input *cwl.WorkflowStepInput, default_name string) (*v1alpha1.Parameter, error) {

	var paramName string
	if input.Id == nil {
		paramName = default_name
	} else {
		paramName = *input.Id
	}

	var stringList []string = strings.Split(*input.Source, "/")
	var inputReference string

	var scope string = stringList[0]
	var key string = stringList[1]

	// If a global input, we will use workflow.parameters... else steps.{{stepId}}.outputs.parameters.{{param.Name}}
	if scope == "global" {
		inputReference = fmt.Sprintf("{{workflow.parameters.%s}}", key)
	} else {
		var updatedScope string = strings.ReplaceAll(scope, "_", "-")

		inputReference = fmt.Sprintf("{{steps.%s.outputs.parameters.%s}}", updatedScope, key)
	}

	var returnParam v1alpha1.Parameter
	returnParam.Name = paramName
	returnParam.Value = (*v1alpha1.AnyString)(&inputReference)

	return &returnParam, nil
}

func cleanArgs(s string) (string, error) {

	startIndex := strings.Index(s, "$(")
	endIndex := strings.Index(s, ")")

	if startIndex == -1 {
		return s, nil
	}

	subString := strings.TrimSpace(s[startIndex : endIndex+1])
	splitStrings := strings.Split(subString, ".")

	if len(splitStrings) > 2 {
		return "", fmt.Errorf("argument substring %+v contains multiple '.' can only contain one", subString)
	}
	addedParam := splitStrings[0] + ".parameters." + splitStrings[1]

	addedParam = strings.Replace(addedParam, "$(", "{{", 1)
	addedParam = strings.Replace(addedParam, ")", "}}", 1)

	updatedString := strings.Replace(s, subString, addedParam, 1)

	return cleanArgs(updatedString)
}

// Need to change the input ref syntax from cwl to argo
func sanitizeArgs(arguments *cwl.Arguments) (*[]string, error) {

	var returnArgs []string
	for _, arg := range *arguments {

		cleanArg, err := cleanArgs(arg)
		if err != nil {
			return nil, err
		}
		returnArgs = append(returnArgs, cleanArg)
	}

	return &returnArgs, nil
}

func EmitCommandArgs(container *apiv1.Container, run *cwl.WorkflowCommandLineTool) error {
	tmpContainer := container.DeepCopy()

	tmpContainer.Command = run.BaseCommand

	cleanArgs, err := sanitizeArgs(&run.Arguments)

	if err != nil {
		return err
	}
	tmpContainer.Args = *cleanArgs

	*container = *tmpContainer

	return nil
}

func EmitStep(step *cwl.WorkflowStep, locations cwl.FileLocations, outputs v1alpha1.Outputs) (*v1alpha1.WorkflowStep, error) {
	outStep := v1alpha1.WorkflowStep{}
	template := v1alpha1.Template{}
	container := apiv1.Container{}

	outStep.Name = strings.Replace(step.Id, "_", "-", -1)

	dockerRequirement, err := findDockerRequirement(step.Requirements)
	if err != nil {
		return nil, err
	}

	err = emitDockerRequirement(&container, dockerRequirement)
	if err != nil {
		return nil, err
	}

	err = EmitCommandArgs(&container, &step.Run)
	if err != nil {
		return nil, err
	}

	template.Container = &container
	var templateInputs []v1alpha1.Parameter

	// A step input can either be a workflow argument, or the output of another step!
	if step.In.Array != nil {
		for idx, input := range step.In.Array {
			var stepName string = "step-" + fmt.Sprint(idx)

			newInput, err := EmitStepInput(&input, stepName)

			if err != nil {
				return nil, err
			}
			templateInputs = append(templateInputs, *newInput)
		}
	} else if step.In.Map != nil {
		for key, input := range step.In.Map {

			newInput, err := EmitStepInput(&input, key)
			if err != nil {
				return nil, err
			}
			templateInputs = append(templateInputs, *newInput)
		}
	}

	// Add the parsed argo outputs into the template if they are relevant!
	for _, output := range step.Out {
		for _, argoOutput := range outputs.Parameters {
			if argoOutput.Name == *output.Id {
				template.Outputs.Parameters = append(template.Outputs.Parameters, argoOutput)
			}
		}
	}

	template.Inputs.Parameters = templateInputs

	outStep.Inline = &template

	return &outStep, nil
}

func emitWorkflowStepOutputs(workflow *cwl.Workflow) (*v1alpha1.Outputs, error) {

	var stepOutputs v1alpha1.Outputs

	// Get output resources from the steps
	for _, step := range workflow.Steps {
		// Get the step run resource

		// For every output, add a tmp parameter and check the step run for the output value from
		for _, out := range step.Out {
			var tmpParameter v1alpha1.Parameter

			tmpParameter.Name = *out.Id

			for _, output := range step.Run.Outputs {

				if *output.ID == tmpParameter.Name {
					var tmpValueFrom v1alpha1.ValueFrom
					tmpValueFrom.Path = *output.OutputBinding.Glob.String
					tmpParameter.ValueFrom = &tmpValueFrom

				}
			}

			stepOutputs.Parameters = append(stepOutputs.Parameters, tmpParameter)

		}

	}

	return &stepOutputs, nil
}

func EmitWorkflow(workflow *cwl.Workflow, inputs map[string]cwl.CWLInputEntry, locations cwl.FileLocations) (*v1alpha1.Workflow, error) {
	var wf v1alpha1.Workflow

	var workflowTemplate v1alpha1.Template

	if workflow.Id != nil {
		wf.Name = *workflow.Id
	} else {
		randomStr := RandomString(10)
		wf.Name = "generated-workflow-" + randomStr
	}

	wf.APIVersion = ArgoVersion
	wf.Kind = ArgoType

	spec := v1alpha1.WorkflowSpec{}

	args, err := EmitWorkflowArguments(&workflow.Inputs)
	if err != nil {
		return nil, err
	}
	spec.Arguments = *args

	// Get the workflow outputs.
	workflowOutputs, err := emitWorkflowStepOutputs(workflow)
	if err != nil {
		return nil, err
	}

	// For every step in the workflow, we create a ParrallelStep
	outSteps := make([]v1alpha1.ParallelSteps, 0)
	for _, step := range workflow.Steps {
		var tmpParralel v1alpha1.ParallelSteps

		tmp, err := EmitStep(&step, locations, *workflowOutputs)

		if err == nil {
			tmpParralel.Steps = append(tmpParralel.Steps, *tmp)

		} else {
			return nil, fmt.Errorf("ran into %+v on step %+v", err, step.Id)
		}

		outSteps = append(outSteps, tmpParralel)
	}

	workflowTemplate.Name = "global-template"
	workflowTemplate.Steps = outSteps

	spec.Entrypoint = workflowTemplate.Name
	spec.Templates = append(spec.Templates, workflowTemplate)

	wf.Spec = spec

	return &wf, nil
}
