package cwl

// CWLRequirements defines the requirements for a CWL tool.
type CWLRequirements interface {
	isCWLRequirement()
	getClass() string
}

// DockerRequirement specifies Docker requirements for a CWL tool.
type DockerRequirement struct {
	Class                 string  `yaml:"class"`
	DockerPull            *string `yaml:"dockerPull"`
	DockerLoad            *string `yaml:"dockerLoad"`
	DockerFile            *string `yaml:"dockerFile"`
	DockerImport          *string `yaml:"dockerImport"`
	DockerImageID         *string `yaml:"dockerImageID"`
	DockerOutputDirectory *string `yaml:"dockerOutputDirectory"`
}

// InlineJavascriptRequirement specifies inline JavaScript requirements for a CWL tool.
type InlineJavascriptRequirement struct {
	Class         string  `yaml:"class"`
	ExpressionLib Strings `yaml:"expressionLib"`
}

// SoftwarePackage represents a software package requirement for a CWL tool.
type SoftwarePackage struct {
	Package string  `yaml:"package"`
	Version Strings `yaml:"version"`
	Specs   Strings `yaml:"specs"`
}

// SoftwareRequirement specifies software requirements for a CWL tool.
type SoftwareRequirement struct {
	Class    string            `yaml:"class"`
	Packages []SoftwarePackage `yaml:"packages"`
}

// LoadListingRequirement specifies the load listing requirements for a CWL tool.
type LoadListingRequirement struct {
	Class       string           `yaml:"class"`
	LoadListing *LoadListingEnum `yaml:"loadListing"`
}

// InitialWorkDirRequirementListing defines initial working directory requirements for a CWL tool.
type InitialWorkDirRequirementListing interface {
	isInitialWorkDirRequirementListing()
}

// InitialWorkDirRequirement specifies initial working directory requirements for a CWL tool.
type InitialWorkDirRequirement struct {
	Class   string                           `yaml:"class"` // constant InitialWorkDirRequirement
	Listing InitialWorkDirRequirementListing `yaml:"listing"`
}

// SchemaDefRequirementType defines the type of schema definition requirement for a CWL tool.
type SchemaDefRequirementType interface {
	isSchemaDefRequirementType()
}

// SchemaDefRequirement specifies schema definition requirements for a CWL tool.
type SchemaDefRequirement struct {
	Class string                     `yaml:"class"` // constant SchemaDefRequirement
	Types []SchemaDefRequirementType `yaml:"types"`
}

// EnvironmentDef defines environment requirements for a CWL tool.
type EnvironmentDef struct {
	EnvName  string        `yaml:"envName"`
	EnvValue CWLExpression `yaml:"envValue"`
}

// EnvVarRequirement specifies environment variable requirements for a CWL tool.
type EnvVarRequirement struct {
	Class  string           `yaml:"class"` // constant EnvVarRequirement
	EnvDef []EnvironmentDef `yaml:"envDef"`
}

// ShellCommandRequirement specifies shell command requirements for a CWL tool.
type ShellCommandRequirement struct {
	Class string `yaml:"class"` // constant ShellCommandRequirement
}

// WorkReuse specifies work reuse requirements for a CWL tool.
type WorkReuse struct {
	Class       string        `yaml:"class"` // constant WorkReuse
	EnableReuse CWLExpression `yaml:"enableReuse"`
}

// NetworkAccess specifies network access requirements for a CWL tool.
type NetworkAccess struct {
	Class         string // constant NetworkAccess
	NetworkAccess CWLExpression
}

// InplaceUpdateRequirement defines inplace update requirements for a CWL tool.
type InplaceUpdateRequirement struct {
	Class         string `yaml:"class"` // constant InplaceUpdateRequirement
	InplaceUpdate Bool   `yaml:"inplaceUpdate"`
}

// ToolTimeLimit specifies time limits for running a CWL tool.
type ToolTimeLimit struct {
	Class     string        `yaml:"class"` // constant ToolTimeLimit
	TimeLimit CWLExpression `yaml:"timeLimit"`
}

// ResourceRequirement defines resource requirements for a CWL tool.
type ResourceRequirement struct {
	Class     string         `yaml:"class"` // constand ResourceRequirement
	CoresMin  *CWLExpression `yaml:"coresMin"`
	CoresMax  *CWLExpression `yaml:"coresMax"`
	RamMin    *CWLExpression `yaml:"ramMin"`
	RamMax    *CWLExpression `yaml:"ramMax"`
	TmpdirMin *CWLExpression `yaml:"tmpdirMin"`
	TmpdirMax *CWLExpression `yaml:"tmpdirMax"`
	OutdirMin *CWLExpression `yaml:"outdirMin"`
	OutdirMax *CWLExpression `yaml:"outdirMax"`
}

func (InlineJavascriptRequirement) isCWLRequirement()  {}
func (d InlineJavascriptRequirement) getClass() string { return d.Class }

func (SchemaDefRequirement) isCWLRequirement()  {}
func (d SchemaDefRequirement) getClass() string { return d.Class }

func (LoadListingRequirement) isCWLRequirement()  {}
func (d LoadListingRequirement) getClass() string { return d.Class }

func (DockerRequirement) isCWLRequirement()  {}
func (d DockerRequirement) getClass() string { return d.Class }

func (SoftwareRequirement) isCWLRequirement()  {}
func (d SoftwareRequirement) getClass() string { return d.Class }

func (InitialWorkDirRequirement) isCWLRequirement()  {}
func (d InitialWorkDirRequirement) getClass() string { return d.Class }

func (EnvVarRequirement) isCWLRequirement()  {}
func (d EnvVarRequirement) getClass() string { return d.Class }

func (ShellCommandRequirement) isCWLRequirement()  {}
func (d ShellCommandRequirement) getClass() string { return d.Class }

func (WorkReuse) isCWLRequirement()  {}
func (d WorkReuse) getClass() string { return d.Class }

func (NetworkAccess) isCWLRequirement()  {}
func (d NetworkAccess) getClass() string { return d.Class }

func (InplaceUpdateRequirement) isCWLRequirement()  {}
func (d InplaceUpdateRequirement) getClass() string { return d.Class }

func (ToolTimeLimit) isCWLRequirement()  {}
func (d ToolTimeLimit) getClass() string { return d.Class }

func (ResourceRequirement) isCWLRequirement()  {}
func (d ResourceRequirement) getClass() string { return d.Class }

func (CommandlineInputRecordSchema) isSchemaDefRequirementType() {}
func (CommandlineInputEnumSchema) isSchemaDefRequirementType()   {}
func (CommandlineInputArraySchema) isSchemaDefRequirementType()  {}
func (DockerRequirement) isSchemaDefRequirementType()            {}
func (SoftwareRequirement) isSchemaDefRequirementType()          {}
func (InitialWorkDirRequirement) isSchemaDefRequirementType()    {}

// CommandlineInputRecordField represents a field in a command-line input record.
type CommandlineInputRecordField struct {
	Name           string              `yaml:"name"`
	Type           CWLTypes            `yaml:"type"` // len(1) represents scalar len > 1 represents array
	Doc            Strings             `yaml:"doc"`
	Label          *string             `yaml:"label"`
	SecondaryFiles SecondaryFiles      `yaml:"secondaryFiles"`
	Streamable     *bool               `yaml:"streamable"`
	Format         CWLFormat           `yaml:"format"`
	LoadContents   *bool               `yaml:"loadContents"`
	LoadListing    LoadListingEnum     `yaml:"loadListing"`
	InputBinding   *CommandlineBinding `yaml:"inputBinding"`
}

// CommandlineInputArraySchema defines the schema for a command-line input array.
type CommandlineInputArraySchema struct {
	Items        CWLTypes            `yaml:"items"`
	Type         string              `yaml:"type"` // MUST be array
	Label        *string             `yaml:"label"`
	Doc          Strings             `yaml:"doc"`
	Name         *string             `yaml:"name"`
	InputBinding *CommandlineBinding `yaml:"inputBinding"`
}

// CommandlineInputEnumSchema specifies the schema for a command-line input enumeration.
type CommandlineInputEnumSchema struct {
	Symbols      Strings `yaml:"symbols"`
	Type         string  `yaml:"type"` // MUST BE enum, only a placeholder for type verification purposes
	Label        *string `yaml:"label"`
	Doc          Strings `yaml:"doc"`
	Name         *string `yaml:"name"`
	InputBinding *CommandlineBinding
}

// CommandlineInputRecordSchema defines the schema for a command-line input record.
type CommandlineInputRecordSchema struct {
	Type   string                         `yaml:"type"` // MUST BE "record"
	Fields *[]CommandlineInputRecordField `yaml:"fields"`
	Label  *string                        `yaml:"label"`
	Doc    *Strings                       `yaml:"doc"`
	Name   *string                        `yaml:"name"`
	// will be used for processing later on hence we disable the linter
	inputBinding *CommandlineBinding `yaml:"inputBinding"` //nolint:unused,structcheck
}

// Type represents a type used in command-line tools.
type Type int32

const (
	// CWLNullKind represents the kind of null value in CWL.
	CWLNullKind Type = iota
	CWLBoolKind
	CWLIntKind
	CWLLongKind
	CWLFloatKind
	CWLDoubleKind
	CWLFileKind
	CWLDirectoryKind
	CWLStdinKind
	CWLStringKind
	CWLRecordKind
	CWLRecordFieldKind
	CWLEnumKind
	CWLArrayKind
)

// CWLType defines a CWL type used in command-line tools.
type CWLType struct {
	Kind   Type
	Record *CommandlineInputRecordSchema
	Enum   *CommandlineInputEnumSchema
	Array  *CommandlineInputArraySchema
	File   *CWLFile
}

// CWLTypes defines multiple CWL types.
type CWLTypes []CWLType

// CommandlineBinding specifies bindings for command-line arguments.
type CommandlineBinding struct {
	LoadContents  *bool         `yaml:"loadContents"`
	Position      *int          `yaml:"position"`
	Prefix        *string       `yaml:"prefix"`
	Separate      *bool         `yaml:"separate"`
	ItemSeperator *string       `yaml:"itemSeperator"`
	ValueFrom     CWLExpression `yaml:"valueFrom"`
	ShellQuote    *bool         `yaml:"bool"`
}

// CommandlineInputParameter represents a command-line input parameter.
type CommandlineInputParameter struct {
	Type           CWLTypes            `yaml:"type"` // len(1) == scalar while len > 1 == array
	Label          *string             `yaml:"label"`
	SecondaryFiles SecondaryFiles      `yaml:"secondaryFiles"` // len(1) == scalar while len > 1 == array
	Streamable     *bool               `yaml:"streamable"`
	Doc            Strings             `yaml:"doc"`
	ID             *string             `yaml:"ID"`
	Format         *CWLFormat          `yaml:"format"`
	LoadContents   *bool               `yaml:"loadContents"`
	LoadListing    *LoadListingEnum    `yaml:"loadListing"`
	Default        interface{}         `yaml:"default"`
	InputBinding   *CommandlineBinding `yaml:"inputBinding"`
}

// OutputBindingGlobKind defines the kind of glob used in output bindings.
type OutputBindingGlobKind int32

const (
	// GlobStringKind represents different types of glob strings.
	GlobStringKind OutputBindingGlobKind = iota
	GlobStringsKind
	GlobExpressionKind
)

// CommandlineOutputBindingGlob defines output binding globs for command-line tools.
type CommandlineOutputBindingGlob struct {
	Kind       OutputBindingGlobKind
	String     *string
	Strings    []string
	Expression CWLExpression
}

// CommandlineOutputBinding specifies output bindings for command-line tools.
type CommandlineOutputBinding struct {
	LoadContents *bool                        `yaml:"loadContents"`
	LoadListing  LoadListingEnum              `yaml:"loadListing"`
	Glob         CommandlineOutputBindingGlob `yaml:"glob"`
	OutputEval   CWLExpression                `yaml:"outputEval"`
}

// CommandlineOutputParameter represents a command-line output parameter.
type CommandlineOutputParameter struct {
	Type           CWLTypes                  `yaml:"type"`
	Label          *string                   `yaml:"label"`
	SecondaryFiles SecondaryFiles            `yaml:"secondaryFiles"`
	Streamable     *bool                     `yaml:"streamable"`
	Doc            Strings                   `yaml:"doc"`
	ID             *string                   `yaml:"ID"`
	Format         *CWLFormat                `yaml:"format"`
	OutputBinding  *CommandlineOutputBinding `yaml:"outputBinding"`
}

// CommandlineArgumentKind defines the kind of command-line arguments.
type CommandlineArgumentKind int32

// ArgumentStringKind represents types of argument strings.
const (
	ArgumentStringKind CommandlineArgumentKind = iota
	ArgumentExpressionKind
	ArgumentCLIBindingKind
)

type CommandlineArgument struct {
	Kind               CommandlineArgumentKind
	String             String
	Expression         CWLExpression
	CommandlineBinding CommandlineBinding
}

type Inputs []CommandlineInputParameter
type Outputs []CommandlineOutputParameter
type Requirements []CWLRequirements
type Hints struct {
	Array []interface{}
	Map   map[string]interface{}
}
type Arguments []string // CommandlineArgument

type CommandLineTool struct {
	Inputs       Inputs         `yaml:"inputs"`
	Outputs      Outputs        `yaml:"outputs"`
	Class        string         `yaml:"class"` // Must be "CommandLineTool"
	ID           *string        `yaml:"id"`
	Label        *string        `yaml:"label"`
	Doc          Strings        `yaml:"doc"`
	Requirements Requirements   `yaml:"requirements"`
	Hints        Hints          `yaml:"hints"`
	CWLVersion   *string        `yaml:"cwlVersion"`
	Intent       Strings        `yaml:"intent"`
	BaseCommand  Strings        `yaml:"baseCommand"`
	Arguments    Arguments      `yaml:"arguments"`
	Stdin        *CWLExpression `yaml:"stdin"`
	Stderr       *CWLExpression `yaml:"stderr"`
	Stdout       *CWLExpression `yaml:"stdout"`
}
