package cwl

import (
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"
)

// intermediate representation used to
// parse into interfaces. The class string is used
// to decode Node into a structure.
type intermediateRepr struct {
	Class *string
	Node  *yaml.Node
}

func getCWLExpressionInner(str string) *string {
	if len(str) < 3 {
		return nil
	}
	if str[0] != '$' {
		return nil
	}
	if str[1] == '[' && str[len(str)-1] == ']' {
		return &str
	}
	if str[1] == '{' && str[len(str)-1] == '}' {
		return &str
	}
	return nil
}

// UnmarshalYAML decodes YAML data into a Strings object.
func (s *Strings) UnmarshalYAML(value *yaml.Node) error {
	strings := make([]string, 0)
	switch value.Kind {
	case yaml.ScalarNode:
		var s string
		if err := value.Decode(&s); err != nil {
			return err
		}
		strings = append(strings, s)
	case yaml.SequenceNode:
		if err := value.Decode(&strings); err != nil {
			return err
		}
	default:
		return errors.New("string | []string expected")
	}
	*s = strings
	return nil
}

// UnmarshalYAML decodes YAML data into a CWLTypes object.
func (tys *CWLTypes) UnmarshalYAML(value *yaml.Node) error {
	newTys := make([]CWLType, 0)
	switch value.Kind {
	case yaml.ScalarNode:
		var s string
		if err := value.Decode(&s); err != nil {
			return err
		}
		var ty CWLType
		switch s {
		case "string":
			ty.Kind = CWLStringKind
		case "null":
			ty.Kind = CWLNullKind
		case "boolean":
			ty.Kind = CWLBoolKind
		case "int":
			ty.Kind = CWLIntKind
		case "long":
			ty.Kind = CWLLongKind
		case "float":
			ty.Kind = CWLFloatKind
		case "double":
			ty.Kind = CWLDoubleKind
		case "File":
			ty.Kind = CWLFileKind
		case "Directory":
			ty.Kind = CWLDirectoryKind
		default:
			return fmt.Errorf("%s is not a supported type", s)
		}
		newTys = append(newTys, ty)
	case yaml.MappingNode:
		return errors.New("complex types not supported yet")
	case yaml.SequenceNode:
		return errors.New("array types not supported yet")
	default:
		return errors.New("type not supported")
	}
	*tys = newTys
	return nil
}

// UnmarshalYAML decodes YAML data into a CWLFormat object.
func (format *CWLFormat) UnmarshalYAML(value *yaml.Node) error {
	switch value.Kind {
	case yaml.ScalarNode:
		var s string
		if err := value.Decode(&s); err != nil {
			return err
		}
		format.Kind = FormatStringKind
		format.String = String(s)
		return nil
	case yaml.SequenceNode:
		s := make([]string, 0)
		if err := value.Decode(&s); err != nil {
			return err
		}
		format.Kind = FormatStringsKind
		format.Strings = s
		return nil
	default:
		return errors.New("string | []string expected")
	}
}

// UnmarshalYAML decodes YAML data into a CommandlineInputParameter object.
func (input *CommandlineInputParameter) UnmarshalYAML(value *yaml.Node) error {
	type rawParamType CommandlineInputParameter

	err := value.Decode((*rawParamType)(input))
	if err != nil {
		return err
	}
	return nil
}

// UnmarshalYAML decodes YAML data into an Inputs object.
func (inp *Inputs) UnmarshalYAML(value *yaml.Node) error {
	inputs := make([]CommandlineInputParameter, 0)
	switch value.Kind {
	case yaml.MappingNode:
		m := make(map[string]CommandlineInputParameter)
		err := value.Decode(&m)
		if err != nil {
			return err
		}
		for key, input := range m {
			newKey := key
			input.ID = &newKey
			inputs = append(inputs, input)
		}
	case yaml.SequenceNode:
		err := value.Decode(&inputs)
		if err != nil {
			return err
		}
	default:
		return errors.New("sequence or mapping type expected")
	}
	*inp = inputs

	return nil
}

// UnmarshalYAML decodes YAML data into an Outputs object.
func (out *Outputs) UnmarshalYAML(value *yaml.Node) error {
	outputs := make([]CommandlineOutputParameter, 0)

	switch value.Kind {
	case yaml.MappingNode:
		m := make(map[string]CommandlineOutputParameter)
		err := value.Decode(&m)
		if err != nil {
			return err
		}
		for key, output := range m {
			newKey := key
			output.ID = &newKey
			outputs = append(outputs, output)
		}
	case yaml.SequenceNode:
		err := value.Decode(&outputs)
		if err != nil {
			return err
		}
	default:
		return errors.New("Sequence or mapping type expected")
	}

	*out = outputs

	return nil
}

func (ir *intermediateRepr) UnmarshalYAML(value *yaml.Node) error {
	m := make(map[string]interface{})
	err := value.Decode(&m)
	if err != nil {
		return err
	}
	classi, ok := m["class"]
	if ok {
		class, ok := classi.(string)
		if !ok {
			return errors.New("string expected")
		}
		ir.Class = &class
	}
	ir.Node = value
	return nil
}

// UnmarshalYAML decodes YAML data into a Requirements object.
func (reqs *Requirements) UnmarshalYAML(value *yaml.Node) error {
	rs := make(map[string]intermediateRepr, 0)
	err := value.Decode(&rs)
	if err != nil {
		rsArray := make([]intermediateRepr, 0)
		err = value.Decode(&rsArray)
		if err != nil {
			return errors.New("[]requirement or map[class]requirement was expected")
		}
		for _, req := range rsArray {
			if req.Class == nil {
				return errors.New("class expected")
			}
			rs[*req.Class] = req
		}
	}

	newRequests := make([]CWLRequirements, 0)
	for class, req := range rs {
		switch class {
		case "DockerRequirement":
			var d DockerRequirement
			err := req.Node.Decode(&d)
			if err != nil {
				return err
			}
			newRequests = append(newRequests, d)
		case "ResourceRequirement":
			var r ResourceRequirement
			err := req.Node.Decode(&r)
			if err != nil {
				return err
			}
			newRequests = append(newRequests, r)
		default:
			return fmt.Errorf("%s is not implemented", class)
		}
	}
	*reqs = newRequests
	return nil
}

// UnmarshalYAML decodes YAML data into a Hints object.
func (h *Hints) UnmarshalYAML(value *yaml.Node) error {

	// Attempt to unmarshal as an array of Any
	var arrayHints []interface{}
	if err := value.Decode(&arrayHints); err == nil {
		h.Array = arrayHints
		return nil
	}

	// Attempt to unmarshal as a map with class keys
	var mapHints map[string]interface{}
	if err := value.Decode(&mapHints); err == nil {
		h.Map = mapHints
		return nil
	}

	return fmt.Errorf("hints must be an array or a map with class keys")
}

// UnmarshalYAML decodes YAML data into a Scatter object.
func (s *Scatter) UnmarshalYAML(value *yaml.Node) error {

	// Attempt to unmarshal as an string
	var stringScatter string
	err := value.Decode(&stringScatter)
	if err == nil {
		s.String = stringScatter
		return nil
	}

	// Attempt to unmarshal as an array of strings
	var arrayScatter []string
	err = value.Decode(&arrayScatter)
	if err == nil {
		s.Array = arrayScatter
		return nil
	}

	return fmt.Errorf("scatter must be a string or an array of strings")
}

// UnmarshalYAML method for ScatterMethod
func (s *ScatterMethod) UnmarshalYAML(value *yaml.Node) error {
	var methodName string
	if err := value.Decode(&methodName); err != nil {
		return err
	}

	switch methodName {
	case string(DotProduct), string(NestedCrossProduct), string(FlatCrossProduct):
		*s = ScatterMethod(methodName)
	default:
		return fmt.Errorf("unknown scatter method: %s", methodName)
	}

	return nil
}

// UnmarshalYAML decodes YAML data into a CWLExpression object.
func (expr *CWLExpression) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind != yaml.ScalarNode {
		return errors.New("can only be string | bool | int | float")
	}

	if value.Tag == "!!int" {

		var i int
		err := value.Decode(&i)
		if err == nil {
			expr.Kind = IntKind
			expr.Int = i
			return nil
		}
	} else if value.Tag == "!!float" {

		var f float64
		err := value.Decode(&f)
		if err == nil {
			expr.Kind = FloatKind
			expr.Float = f
			return nil
		}
	}

	var b bool
	err := value.Decode(&b)
	if err == nil {
		expr.Kind = BoolKind
		expr.Bool = b
		return nil
	}

	var s string
	err = value.Decode(&s)
	if err == nil {
		exprS := getCWLExpressionInner(s)
		if exprS != nil {
			expr.Kind = ExpressionKind
			expr.Expression = *exprS
			return nil
		}
		expr.Kind = RawKind
		expr.Raw = s
		return nil
	}
	return errors.New("can only be string | bool | int | float")
}

// UnmarshalYAML decodes YAML data into a CommandLineTool object.
func (cl *CommandLineTool) UnmarshalYAML(value *yaml.Node) error {
	type rawCLITool CommandLineTool
	if err := value.Decode((*rawCLITool)(cl)); err != nil {
		return err
	}
	return nil
}

// UnmarshalYAML decodes YAML data into a CommandlineOutputBindingGlob object.
func (clOutputBindingGlob *CommandlineOutputBindingGlob) UnmarshalYAML(value *yaml.Node) error {
	var s string
	err := value.Decode(&s)
	if err == nil {
		exprS := getCWLExpressionInner(s)
		if exprS != nil {
			clOutputBindingGlob.Kind = GlobExpressionKind
			clOutputBindingGlob.Expression = CWLExpression{Kind: ExpressionKind, Expression: *exprS}
			return nil
		}
		clOutputBindingGlob.Kind = GlobStringKind
		clOutputBindingGlob.String = &s
		return nil
	}

	ss := make([]string, 0)
	err = value.Decode(&ss)
	if err == nil {
		return nil
	}
	return nil
}
