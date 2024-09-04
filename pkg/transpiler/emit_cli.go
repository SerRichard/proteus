package transpiler

import (
	"errors"
	"fmt"
	"math"
	"sort"

	log "github.com/sirupsen/logrus"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/SerRichard/proteus/pkg/cwl"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
)

const (
	ArgoType             = "Workflow"
	ArgoVersion          = "argoproj.io/v1alpha1"
	volumeClaimName      = "argovolume"
	volumeClaimMountPath = "/mnt/pvol"
)

// de-sum typed "CommandlineInputParameter
// contains all necessary data to output argo yaml
type flatCommandlineInputParameter struct {
	Type             cwl.Type
	Label            *string
	StringValue      *string               // string value
	BoolValue        *bool                 // boolean value
	IntValue         *int                  // int value
	ArrayValue       *[]any                // Array value
	Emit             bool                  // boolean value
	File             *cwl.CWLFile          // file value
	FileLocationData *cwl.FileLocationData // file location data
	SecondaryFiles   cwl.SecondaryFiles
	Streamable       *bool
	Doc              cwl.Strings
	Id               *string
	Format           *cwl.CWLFormat
	LoadContents     *bool
	LoadListing      *cwl.LoadListingEnum
	InputBinding     *cwl.CommandlineBinding
}

type flatCommandlineOutputParameter struct {
	Type             cwl.Type
	Label            *string
	FileLocationData *cwl.FileLocationData
	File             *cwl.CWLFile
	SecondaryFiles   cwl.SecondaryFiles
	Streamable       *bool
	Doc              cwl.Strings
	Id               *string
	Format           *cwl.CWLFormat
	OutputBinding    *cwl.CommandlineOutputBinding
}

func emitDockerRequirement(container *apiv1.Container, d *cwl.DockerRequirement) error {
	tmpContainer := container.DeepCopy()

	if d.DockerPull == nil {
		return errors.New("dockerPull is a required field")
	}

	tmpContainer.Image = *d.DockerPull

	if d.DockerOutputDirectory != nil {
		tmpContainer.WorkingDir = *d.DockerOutputDirectory
		log.Warn("I am assuming the DockerOutputDirectory and WorkingDir are equivalent")
		log.Infof("Changing docker WorkingDir to %s", tmpContainer.WorkingDir)
	}

	if d.DockerFile != nil {
		return errors.New("dockerfile is not currently supported")
	}

	if d.DockerImageID != nil {
		return errors.New("dockerimageid is not currently supported")
	}

	if d.DockerImport != nil {
		return errors.New("docker import is not currently supported")
	}

	*container = *tmpContainer
	return nil
}

func emitInputParam(input flatCommandlineInputParameter) *v1alpha1.Parameter {
	name := *input.Id
	param := v1alpha1.Parameter{Name: name}
	return &param
}

func dockerNotPresent() error              { return errors.New("DockerRequirement was not found") }
func resourceRequirementNotPresent() error { return errors.New("ResourceRequirement was not found") }

func findDockerRequirement(requirements cwl.Requirements) (*cwl.DockerRequirement, error) {
	var docker *cwl.DockerRequirement
	docker = nil
	for _, req := range requirements {
		d, ok := req.(cwl.DockerRequirement)
		if ok {
			docker = &d
		}
	}

	if docker != nil {
		return docker, nil
	} else {
		return nil, dockerNotPresent()
	}
}

func findResourceRequirement(requirements cwl.Requirements) (*cwl.ResourceRequirement, error) {
	var resource *cwl.ResourceRequirement

	resource = nil
	for _, req := range requirements {
		r, ok := req.(cwl.ResourceRequirement)
		if ok {
			resource = &r
		}
	}
	if resource != nil {
		return resource, nil
	} else {
		return nil, resourceRequirementNotPresent()
	}
}

func emitInputParams(template *v1alpha1.Template, inputs []flatCommandlineInputParameter) {
	params := make([]v1alpha1.Parameter, 0)
	for _, input := range inputs {
		newInput := emitInputParam(input)
		params = append(params, *newInput)
	}
	template.Inputs.Parameters = params
}

func canFindType(input cwl.CWLInputEntry, tys cwl.CWLTypes) error {
	for _, currTy := range tys {
		if currTy.Kind == input.Kind {
			return nil
		}
	}
	return errors.New("unable to find type")
}

type CommandlineInputParameter struct {
	cwl.CommandlineInputParameter
}

func (inputParameter CommandlineInputParameter) getInputBindings(inputs map[string]cwl.CWLInputEntry) (*flatCommandlineInputParameter, error) {
	if inputParameter.CommandlineInputParameter.ID == nil {
		return nil, errors.New("input parameter is nil")
	}

	binding := flatCommandlineInputParameter{
		Type:           inputParameter.CommandlineInputParameter.Type[0].Kind,
		SecondaryFiles: inputParameter.CommandlineInputParameter.SecondaryFiles,
		Streamable:     inputParameter.CommandlineInputParameter.Streamable,
		Doc:            inputParameter.CommandlineInputParameter.Doc,
		Id:             inputParameter.CommandlineInputParameter.ID,
		Format:         inputParameter.CommandlineInputParameter.Format,
		InputBinding:   inputParameter.CommandlineInputParameter.InputBinding,
		Emit:           true,
	}

	if len(inputs) == 0 {
		switch v := inputParameter.CommandlineInputParameter.Default.(type) {
		case string:
			binding.StringValue = &v
		}
		return &binding, nil
	}

	input, ok := inputs[*inputParameter.CommandlineInputParameter.ID]

	if !ok {
		return nil, fmt.Errorf("%s was not present in input", *inputParameter.CommandlineInputParameter.ID)
	}

	err := canFindType(input, inputParameter.CommandlineInputParameter.Type)
	if err != nil {
		return nil, err
	}

	binding.Type = input.Kind
	switch input.Kind {
	case cwl.CWLStringKind:
		binding.StringValue = input.StringData
	case cwl.CWLIntKind:
		binding.IntValue = input.IntData
	case cwl.CWLFileKind:
		binding.File = input.FileData
	case cwl.CWLBoolKind:
		binding.BoolValue = input.BoolData
	case cwl.CWLArrayKind:
		binding.ArrayValue = input.Array
	default:
		return nil, fmt.Errorf("%T unknown type", input.Kind)
	}

	return &binding, nil
}

type CommandlineOutputParameter struct {
	cwl.CommandlineOutputParameter
}

func (outputParameter CommandlineOutputParameter) getOutputBindings() (*flatCommandlineOutputParameter, error) {
	binding := flatCommandlineOutputParameter{
		SecondaryFiles: outputParameter.CommandlineOutputParameter.SecondaryFiles,
		Streamable:     outputParameter.CommandlineOutputParameter.Streamable,
		Doc:            outputParameter.CommandlineOutputParameter.Doc,
		Id:             outputParameter.CommandlineOutputParameter.ID,
		Format:         outputParameter.CommandlineOutputParameter.Format,
		OutputBinding:  outputParameter.CommandlineOutputParameter.OutputBinding,
	}

	if len(outputParameter.CommandlineOutputParameter.Type) != 1 {
		return nil, fmt.Errorf("only single output types expected: expected len(Type)==1 got len(Type)==%d in array %v", len(outputParameter.CommandlineOutputParameter.Type), outputParameter.Type)
	}
	ty := outputParameter.CommandlineOutputParameter.Type[0].Kind
	switch ty {
	case cwl.CWLStringKind:
		break
	case cwl.CWLIntKind:
		break
	case cwl.CWLFileKind:
		break
	default:
		return nil, fmt.Errorf("%T unknown type", ty)
	}
	binding.Type = ty
	return &binding, nil
}

func sortBindingsByPosition(bindings []flatCommandlineInputParameter) {
	sort.Slice(bindings[:], func(i, j int) bool {
		leftPost := 0
		rightPost := 0
		if bindings[i].InputBinding.Position != nil {
			leftPost = *bindings[i].InputBinding.Position
		}
		if bindings[i].InputBinding.Position != nil {
			rightPost = *bindings[j].InputBinding.Position
		}
		return leftPost < rightPost
	})
}

func emitArgumentParams(container *apiv1.Container,
	baseCommand cwl.Strings,
	arguments cwl.Arguments,
	bindings []flatCommandlineInputParameter) error {
	cmds := make([]string, 0)
	skip := false

	if len(baseCommand) == 0 {
		if len(arguments) == 0 {
			return errors.New("len(baseCommand)==0 && len(arguments)==0")
		}

		cmds = append(cmds, arguments[0])
		skip = false
	}

	for _, cmd := range baseCommand {
		cmds = append(cmds, cmd)
	}

	for i, arg := range arguments {
		if i == 0 && skip {
			continue
		}

		cmds = append(cmds, arg)
	}

	sortBindingsByPosition(bindings)

	args := make([]string, 0)
	for _, binding := range bindings {

		prefix := ""
		if binding.InputBinding != nil && binding.InputBinding.Prefix != nil {
			sep := true
			if binding.InputBinding.Separate != nil {
				sep = *binding.InputBinding.Separate
			}

			if sep {
				sepArg := *binding.InputBinding.Prefix
				args = append(args, sepArg)
			} else {
				prefix = *binding.InputBinding.Prefix
			}
		}
		var arg string
		arg = fmt.Sprintf("%s{{inputs.parameters.%s}}", prefix, *binding.Id)

		if binding.Type == cwl.CWLFileKind {
			if binding.InputBinding == nil {
				continue
			}
			if binding.File == nil || binding.File.Path == nil {
				log.Info("binding, ", binding.File)
				return errors.New("file information was not available")
			}
			arg = *binding.File.Path
		}
		args = append(args, arg)
	}

	container.Command = cmds
	container.Args = args

	return nil
}

func emitArguments(spec *v1alpha1.WorkflowSpec, bindings []flatCommandlineInputParameter) error {
	params := make([]v1alpha1.Parameter, 0)
	arts := make([]v1alpha1.Artifact, 0)

	for _, binding := range bindings {

		switch binding.Type {
		case cwl.CWLStringKind:
			params = append(params, v1alpha1.Parameter{Name: *binding.Id, Value: (*v1alpha1.AnyString)(binding.StringValue)})
		case cwl.CWLIntKind:
			intString := fmt.Sprintf("%d", *binding.IntValue)
			params = append(params, v1alpha1.Parameter{Name: *binding.Id, Value: (*v1alpha1.AnyString)(&intString)})
		case cwl.CWLBoolKind:
			boolString := fmt.Sprintf("%t", *binding.BoolValue)
			params = append(params, v1alpha1.Parameter{Name: *binding.Id, Value: (*v1alpha1.AnyString)(&boolString)})
		default:
			log.Info("HERE ", binding.Type)
			return fmt.Errorf("%T is not supported", binding.Type)
		}
	}
	args := v1alpha1.Arguments{Parameters: params, Artifacts: arts}
	spec.Arguments = args
	return nil
}

func flattenInput(inputs *cwl.Inputs, input map[string]cwl.CWLInputEntry) ([]flatCommandlineInputParameter, error) {
	flatInputs := make([]flatCommandlineInputParameter, 0)
	for _, inputBinding := range *inputs {

		var tmpInputParam = CommandlineInputParameter{inputBinding}

		newBindings, err := tmpInputParam.getInputBindings(input)
		if err != nil {
			return nil, err
		}
		flatInputs = append(flatInputs, *newBindings)
	}
	return flatInputs, nil
}

func flattenOutput(outputs *cwl.Outputs) ([]flatCommandlineOutputParameter, error) {
	flatOutputs := make([]flatCommandlineOutputParameter, 0)
	for _, outputBinding := range *outputs {

		var tmpOutputParam = CommandlineOutputParameter{outputBinding}

		newBindings, err := tmpOutputParam.getOutputBindings()
		if err != nil {
			return nil, err
		}
		flatOutputs = append(flatOutputs, *newBindings)
	}
	return flatOutputs, nil
}

func filterParams(inputs []flatCommandlineInputParameter) []flatCommandlineInputParameter {
	newInputs := make([]flatCommandlineInputParameter, 0)
	for _, input := range inputs {
		switch input.Type {
		case cwl.CWLFileKind:
			continue
		case cwl.CWLRecordFieldKind:
			continue
		case cwl.CWLArrayKind:
			continue
		case cwl.CWLEnumKind:
			continue
		default:
			newInputs = append(newInputs, input)
		}
	}
	return newInputs
}

func needPVC(outputs []flatCommandlineOutputParameter) bool {
	for _, binding := range outputs {
		if binding.Type == cwl.CWLFileKind {
			return true
		}
	}
	return false
}

func expressionToQuantity(expr *cwl.CWLExpression) (*resource.Quantity, error) {
	var qstr string

	switch expr.Kind {
	case cwl.RawKind:
		qstr = expr.Raw
	case cwl.IntKind:
		qstr = fmt.Sprintf("%dMi", expr.Int)
	case cwl.FloatKind:
		round := int(math.Ceil(expr.Float))
		qstr = fmt.Sprintf("%dMi", round)
	default:
		return nil, fmt.Errorf("%T is not a supported type for quantity conversion", expr.Kind)
	}

	quantity, err := resource.ParseQuantity(qstr)
	if err != nil {
		return nil, err
	}
	return &quantity, nil
}

func emitPVC(spec *v1alpha1.WorkflowSpec, resourceReq *cwl.ResourceRequirement) error {
	pSpec := apiv1.PersistentVolumeClaimSpec{}
	resources := apiv1.ResourceRequirements{}
	resourceMap := make(map[apiv1.ResourceName]resource.Quantity)

	quantity, err := expressionToQuantity(resourceReq.OutdirMin)
	if err != nil {
		return err
	}
	resourceMap[apiv1.ResourceStorage] = *quantity
	resources.Requests = resourceMap
	pSpec.Resources = resources
	pSpec.AccessModes = []apiv1.PersistentVolumeAccessMode{apiv1.ReadWriteMany}
	pVolClaim := apiv1.PersistentVolumeClaim{}
	pVolClaim.Spec = pSpec

	pVolClaim.Name = volumeClaimName

	spec.VolumeClaimTemplates = []apiv1.PersistentVolumeClaim{pVolClaim}
	return nil
}

func emitInputArtifacts(template *v1alpha1.Template, inputs map[string]cwl.CWLInputEntry, locations cwl.FileLocations) error {
	arts := make([]v1alpha1.Artifact, 0)

	// If there are no locations, do not try to infer an artifact should exist.
	if len(locations.Inputs) == 0 {
		return nil
	}

	for key, inputEntry := range inputs {
		if inputEntry.Kind != cwl.CWLFileKind {
			continue
		}
		location, ok := locations.Inputs[key]
		if !ok {
			return fmt.Errorf("location data not present for %s", key)
		}

		art := v1alpha1.Artifact{}
		art.Name = location.Name
		art.Path = *inputEntry.FileData.Path
		art.HTTP = location.HTTP
		art.S3 = location.S3
		arts = append(arts, art)
	}

	template.Inputs.Artifacts = arts
	return nil
}

func evalCommandlineBindingOutputGlob(bglob *cwl.CommandlineOutputBindingGlob) (string, error) {
	if bglob == nil {
		return "", errors.New("output binding invalid")
	}
	switch bglob.Kind {
	case cwl.GlobStringKind:
		return *bglob.String, nil
	default:
		return "", errors.New("only string is supported at the moment")
	}
}

func emitOutputArtifact(tmpl *v1alpha1.Template, output flatCommandlineOutputParameter, locations cwl.FileLocations) error {

	// If there are no locations, do not try to infer an artifact should exist.
	if len(locations.Outputs) == 0 {
		return nil
	}

	if output.Type != cwl.CWLFileKind {
		return errors.New("emitOutputArtifact only accepts CWLFileKind")
	}

	path, err := evalCommandlineBindingOutputGlob(&output.OutputBinding.Glob)
	if err != nil {
		return err
	}

	location, ok := locations.Outputs[*output.Id]
	if !ok {
		return fmt.Errorf("unable to find output for %s", *output.Id)
	}
	art := v1alpha1.Artifact{Name: *output.Id, Path: path}
	art.HTTP = location.HTTP
	art.S3 = location.S3

	tmpl.Outputs.Artifacts = append(tmpl.Outputs.Artifacts, art)
	return nil
}

func emitOutputs(tmpl *v1alpha1.Template, outputs []flatCommandlineOutputParameter, locations cwl.FileLocations) error {
	for _, output := range outputs {
		switch output.Type {
		case cwl.CWLFileKind:
			err := emitOutputArtifact(tmpl, output, locations)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("%T is not supported", output.Type)
		}
	}
	return nil
}

func attachVolume(container *apiv1.Container, volumeName string, mountpath string) {
	if container.WorkingDir != "" {
		mountpath = container.WorkingDir
	}

	mnt := apiv1.VolumeMount{}
	mnt.Name = volumeName

	mnt.MountPath = mountpath
	container.VolumeMounts = []apiv1.VolumeMount{mnt}
}

func EmitCommandlineTool(clTool *cwl.CommandLineTool, inputs map[string]cwl.CWLInputEntry, locations cwl.FileLocations) (*v1alpha1.Workflow, error) {
	var wf v1alpha1.Workflow
	var err error

	wf.Name = *clTool.ID
	spec := v1alpha1.WorkflowSpec{}
	wf.APIVersion = ArgoVersion
	wf.Kind = ArgoType

	container := apiv1.Container{}

	dockerRequirement, err := findDockerRequirement(clTool.Requirements)
	if err != nil {
		return nil, err
	}

	err = emitDockerRequirement(&container, dockerRequirement)
	if err != nil {
		return nil, err
	}

	template := v1alpha1.Template{}
	template.Container = &container
	template.Name = *clTool.ID

	bindings, err := flattenInput(&clTool.Inputs, inputs)
	if err != nil {
		return nil, err
	}

	paramBindings := filterParams(bindings)

	emitInputParams(&template, paramBindings)

	outputBindings, err := flattenOutput(&clTool.Outputs)
	if err != nil {
		return nil, err
	}

	log.Info("Need pvc? ", outputBindings)

	if needPVC(outputBindings) {

		resourceRequirement, err := findResourceRequirement(clTool.Requirements)
		if err != nil {
			return nil, err
		}

		err = emitPVC(&spec, resourceRequirement)
		if err != nil {
			return nil, err
		}
		attachVolume(&container, volumeClaimName, volumeClaimMountPath)
	}

	err = emitArgumentParams(&container, clTool.BaseCommand, clTool.Arguments, bindings)
	if err != nil {
		return nil, err
	}

	err = emitArguments(&spec, paramBindings)
	if err != nil {
		return nil, err
	}

	err = emitInputArtifacts(&template, inputs, locations)
	if err != nil {
		return nil, err
	}

	err = emitOutputs(&template, outputBindings, locations)
	if err != nil {
		return nil, err
	}

	spec.Templates = []v1alpha1.Template{template}
	spec.Entrypoint = template.Name

	wf.Spec = spec
	return &wf, nil
}
