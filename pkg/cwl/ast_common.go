package cwl

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"gopkg.in/yaml.v3"
)

const (
	CWLVersion     = "v1.2"
	localPath      = "/usr/share/commonwl/"
	localsharePath = "/usr/local/share/commonwl/"
	homesharePath  = "%+v/.local/share/commonwl/"
)

type (
	String         string
	Bool           bool
	Int            int
	Float          float32
	Strings        []string
	SecondaryFiles []CWLSecondaryFileSchema
)

type CWLFormatKind int32

const (
	FormatStringKind CWLFormatKind = iota
	FormatStringsKind
	FormatExpressionKind
)

type CWLExpressionKind int32

type CWLExpression struct {
	Kind       CWLExpressionKind
	Raw        string
	Expression string
	Bool       bool
	Int        int
	Float      float64
}

const (
	RawKind CWLExpressionKind = iota
	ExpressionKind
	BoolKind
	IntKind
	FloatKind
)

type CWLSecondaryFileSchema struct {
	Pattern  CWLExpression `yaml:"pattern"`
	Required CWLExpression `yaml:"required"`
}

type CWLFormat struct {
	Kind       CWLFormatKind
	String     String
	Strings    Strings
	Expression CWLExpression
}

type CWLStdin struct{}

type CWLInputProvider interface {
	GetKind() string
}

type CWLFile struct {
	Class          string     `yaml:"class"` // constant value File
	Location       *string    `yaml:"location"`
	Path           *string    `yaml:"path"`
	Dirname        *string    `yaml:"dirname"`
	Nameroot       *string    `yaml:"nameroot"`
	Nameext        *string    `yaml:"nameext"`
	Checksum       *string    `yaml:"checksum"`
	Size           *int64     `yaml:"size"`
	SecondaryFiles []CWLFile  `yaml:"secondaryFiles"`
	Format         *CWLFormat `yaml:"format"`
	Contents       *string    `yaml:"contents"`
}

type CWLInputEntry struct {
	Kind       Type
	FileData   *CWLFile
	BoolData   *bool
	StringData *string
	IntData    *int
	Array      *[]any
}

type LoadListingEnum string

const (
	LoadListingNone    LoadListingEnum = "no_listing"
	LoadListingDeep    LoadListingEnum = "deep_listing"
	LoadListingShallow LoadListingEnum = "shallow_listing"
)

type FileLocationKind string

const (
	HTTPKind FileLocationKind = "http"
	S3Kind   FileLocationKind = "s3"
	GITKind  FileLocationKind = "git"
)

type FileLocationData struct {
	Name string                 `json:"name"`
	Type FileLocationKind       `json:"type"`
	HTTP *v1alpha1.HTTPArtifact `json:"http"`
	S3   *v1alpha1.S3Artifact   `json:"s3"`
	HDFS *v1alpha1.HDFSArtifact `json:"hdfs"`
}

type FileLocations struct {
	Inputs  map[string]FileLocationData `json:"inputs"`
	Outputs map[string]FileLocationData `json:"outputs"`
}

func (cwlInputEntry *CWLInputEntry) UnmarshalYAML(value *yaml.Node) error {

	var b bool
	err := value.Decode(&b)
	if err == nil {
		cwlInputEntry.Kind = CWLBoolKind
		cwlInputEntry.BoolData = &b
		return nil
	}

	var i int
	err = value.Decode(&i)
	if err == nil {
		cwlInputEntry.Kind = CWLIntKind
		cwlInputEntry.IntData = &i
		return nil
	}

	var s string
	err = value.Decode(&s)
	if err == nil {
		cwlInputEntry.Kind = CWLStringKind
		cwlInputEntry.StringData = &s
		return nil
	}

	var arr []any
	err = value.Decode(&arr)
	if err == nil {
		cwlInputEntry.Kind = CWLArrayKind
		cwlInputEntry.Array = &arr
		return nil
	}

	var file CWLFile
	err = value.Decode(&file)
	if err == nil {
		if file.Class != "File" {
			return fmt.Errorf("%s was received instead of ", file.Class)
		}
		cwlInputEntry.Kind = CWLFileKind
		cwlInputEntry.FileData = &file
		return nil
	}

	return errors.New("unable to convert into CWLInputEntry")
}

func (f *FileLocationData) UnmarshalJSON(data []byte) error {
	type rawFileLocationData FileLocationData

	tmp := rawFileLocationData{}

	err := json.Unmarshal(data, &tmp)

	if err != nil {
		return err
	}

	switch tmp.Type {
	case HTTPKind:
		if tmp.HTTP == nil {
			return errors.New("http data not provided")
		}
	case S3Kind:
		if tmp.S3 == nil {
			return errors.New("s3 data not provided")
		}
	default:
		return fmt.Errorf("%s is not a valid type", tmp.Type)
	}
	f.Name = tmp.Name
	f.Type = tmp.Type
	f.HTTP = tmp.HTTP
	f.S3 = tmp.S3
	return nil
}
