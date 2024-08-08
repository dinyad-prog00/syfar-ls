package tmp

import (
	"reflect"

	"github.com/alecthomas/participle/v2"
)

func ParseFile(content string, filename string) (*SyfarFile, error) {
	var ps = participle.MustBuild[SyfarFile](participle.Unquote())

	ast, err := ps.ParseString(filename, string(content))

	if err != nil {
		return nil, err
	}
	return ast, nil
}

func GetVariableList(ast *SyfarFile) []Variable {
	list := []Variable{}
	for _, entry := range ast.Entries {
		if entry.Variable != nil {
			list = append(list, *entry.Variable)
		}
	}
	return list
}

func GetMultiVarList(ast *SyfarFile) []MultiVariable {
	list := []MultiVariable{}
	for _, entry := range ast.Entries {
		if entry.MultiVariable != nil {
			list = append(list, *entry.MultiVariable)
		}
	}
	return list
}

func GetVarSetList(ast *SyfarFile) []VarSet {
	list := []VarSet{}
	for _, entry := range ast.Entries {
		if entry.VarSet != nil {
			list = append(list, *entry.VarSet)
		}
	}
	return list
}

func GetActionList(ast *SyfarFile) []Action {
	list := []Action{}
	for _, entry := range ast.Entries {
		if entry.Action != nil {
			list = append(list, *entry.Action)
		}
	}
	return list
}

func GetVarSet(ast *SyfarFile, id string) (*VarSet, bool) {
	for _, entry := range ast.Entries {
		if entry.VarSet != nil && entry.VarSet.Id == id {
			return entry.VarSet, true
		}
	}
	return nil, false
}

func GetSecretSetList(ast *SyfarFile) []SecretSet {
	list := []SecretSet{}
	for _, entry := range ast.Entries {
		if entry.SecretSet != nil {
			list = append(list, *entry.SecretSet)
		}
	}
	return list
}

func GetSecretSet(ast *SyfarFile, id string) (*SecretSet, bool) {
	for _, entry := range ast.Entries {
		if entry.SecretSet != nil && entry.SecretSet.Id == id {
			return entry.SecretSet, true
		}
	}
	return nil, false
}

func GetActionsParametersName(act Action) []string {
	list := []string{}
	for _, at := range act.Attributes {
		if at.Parameter != nil {
			list = append(list, at.Parameter.Name)
		}
	}
	return list
}

type Input struct {
	Name     string
	Type     reflect.Kind
	Required bool
}

type Ouput struct {
	Name string
	Type reflect.Kind
}

var RequestInput = []Input{
	{Name: "url", Type: reflect.String, Required: true},
	{Name: "method", Type: reflect.String, Required: false},
	{Name: "params", Type: reflect.Map, Required: false},
	{Name: "query", Type: reflect.Map, Required: false},
	{Name: "body", Type: reflect.Map, Required: false},
	{Name: "headers", Type: reflect.Map, Required: false},
}

var RequestOutput = []Ouput{
	{Name: "status", Type: reflect.String},
	{Name: "statusCode", Type: reflect.Int},
	{Name: "body", Type: reflect.String},
	{Name: "json", Type: reflect.Map},
}

var ReadFileInput = []Input{
	{Name: "path", Type: reflect.String, Required: true},
}

var ReadFileOuput = []Ouput{
	{Name: "content", Type: reflect.String},
	{Name: "ok", Type: reflect.Bool},
}

var Parameters = map[string][]Input{"file_read": ReadFileInput, "http_request": RequestInput}

var Outputs = map[string][]Ouput{"file_read": ReadFileOuput, "http_request": RequestOutput}
