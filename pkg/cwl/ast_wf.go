package cwl

type InputEnumSchema struct {
	Symbols []string
	Type    string //constant enum
	Label   *string
	Doc     *[]string
	Name    *string
}

type InputArraySchema struct {
	Items []SharedWorkflowInputSumTypes
	Type  string //constant array
	Label *string
	Doc   *[]string
	Name  string
}

type InputRecordField struct {
	Name           string
	Type           CWLTypes
	Doc            *[]string
	Label          *string
	SecondaryFiles *CWLSecondaryFileSchema
	Streamable     *bool
	Format         *CWLFormat
	LoadContents   *bool
	LoadListing    *LoadListingEnum
}

type InputRecordSchema struct {
	Type   string // constant record
	Fields *[]InputRecordField
	Label  *string
	Doc    *[]string
	Name   *string
}

type SharedWorkflowInputSumTypes interface {
	isSharedWorkflowInputSumType()
}

type InputBinding struct {
	LoadContents *bool
}

type WorkflowInputParameter struct {
	Type           CWLTypes
	Label          *string
	SecondaryFiles SecondaryFiles
	Streamable     *bool
	Doc            *[]string
	Id             *string
	Format         *CWLFormat
	LoadContents   *bool
	LoadListing    *LoadListingEnum
	Default        *string
	InputBinding   InputBinding
}

type WorkflowOutputParameterType interface{}

type LinkMergeMethod interface {
	isLinkMergeMethod()
}
type MergeNested struct{}
type MergeFlattened struct{}

func (_ MergeNested) isLinkMergeMethod()    {}
func (_ MergeFlattened) isLinkMergeMethod() {}

type PickValueMethod interface {
	isPickValueMethod()
}
type FirstNonNull struct{}
type TheOnlyNonNull struct{}
type AllNonNull struct{}

func (_ FirstNonNull) isPickValueMethod()   {}
func (_ TheOnlyNonNull) isPickValueMethod() {}
func (_ AllNonNull) isPickValueMethod()     {}

type SharedWorkflowOutputSumTypes interface {
	isSharedWorkflowOutputSumType()
}

type OutputRecordFields struct {
	Name           string
	Type           SharedWorkflowOutputSumTypes
	Doc            []string
	Label          *string
	SecondaryFiles []CWLSecondaryFileSchema
	Streamable     *bool
	Format         *CWLFormat
}

type OutputRecordSchema struct {
	Type   string //constant record
	Fields []OutputRecordFields
	Label  *string
	Doc    []string
	Name   *string
}

type WorkflowOutputParameter struct {
	Type           CWLTypes
	Label          *string
	SecondaryFiles []CWLSecondaryFileSchema
	Streamable     *bool
	Doc            []string
	Id             *string
	Format         *CWLFormat
	OutputSource   []string
	LinkMerge      *LinkMergeMethod
	PickValue      *PickValueMethod
}

type WorkflowStepInput struct {
	Id           *string          `yaml:"id"`
	Source       *string          `yaml:"source"`
	LinkMerge    *string          `yaml:"linkMerge"`
	LoadContents *bool            `yaml:"loadContents"`
	LoadListing  *LoadListingEnum `yaml:"loadListing"`
	Label        *string          `yaml:"label"`
	Default      *string          `yaml:"default"`
	ValueFrom    *CWLExpression   `yaml:"valueFrom"`
}

type WorkflowStepOutput struct {
	Id *string
}

type WorkflowRunnable interface {
	isWorkflowRunnable()
}

func (_ CommandLineTool) isWorkflowRunnable() {}
func (_ String) isWorkflowRunnable()          {}
func (_ Workflow) isWorkflowRunnable()        {}

type SubworkflowFeatureRequirement struct {
	Class string // constant SubworkflowFeatureRequirement
}
type ScatterFeatureRequirement struct {
	Class string // constant ScatterFeatureRequirement
}
type MultipleInputFeatureRequirement struct {
	Class string // constant MultipleInputFeatureRequirement
}

type StepInputExpressionRequirement struct {
	Class string // constant StepInputExpressionRequirement
}

type ScatterMethod string

const (
	DotProduct         ScatterMethod = "dotproduct"
	NestedCrossProduct ScatterMethod = "nested_crossproduct"
	FlatCrossProduct   ScatterMethod = "flat_crossproduct"
)

type Scatter struct {
	String string
	Array  []string
}

type WorkflowCommandLineTool struct {
	CommandLineTool
}

type WorkflowStep struct {
	In            WorkflowStepInputs      `yaml:"in"`
	Out           WorkflowStepOutputs     `yaml:"out"`
	Run           WorkflowCommandLineTool `yaml:"run"`
	Id            string                  `yaml:"id"`
	Label         *string                 `yaml:"label"`
	Doc           Strings                 `yaml:"doc"`
	Requirements  Requirements            `yaml:"requirements"`
	Hints         Hints                   `yaml:"hints"`
	Scatter       Scatter                 `yaml:"scatter"`
	ScatterMethod ScatterMethod           `yaml:"scatterMethod"`
}

type WorkflowStepInputs struct {
	Array []WorkflowStepInput
	Map   map[string]WorkflowStepInput
}

type WorkflowStepOutputs []WorkflowStepOutput

type WorkflowInputs map[string]WorkflowInputParameter
type WorkflowOutputs map[string]WorkflowOutputParameter
type WorkflowSteps []WorkflowStep

type Workflow struct {
	Inputs       WorkflowInputs  `yaml:"inputs"`
	Outputs      WorkflowOutputs `yaml:"outputs"`
	Class        string          `yaml:"class"` // Only Workflow
	Steps        WorkflowSteps   `yaml:"steps"`
	Id           *string         `yaml:"id"`
	Label        *string         `yaml:"label"`
	Doc          []string        `yaml:"doc"`
	Requirements Requirements    `yaml:"requirements"`
	Hints        Hints           `yaml:"hints"`
	CWLVersion   *string         `yaml:"cwlVersion"`
	Intent       []string        `yaml:"intent"`
}
