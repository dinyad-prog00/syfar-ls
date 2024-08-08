package completion

import (
	"syfar-ls/helpers"

	protocol "github.com/tliron/glsp/protocol_3_16"
)

func CreateSnippet(label string, snippet string, detail string) protocol.CompletionItem {
	return protocol.CompletionItem{
		Label:            label,
		Detail:           helpers.StrPtr(detail),
		Kind:             helpers.KindPtr(protocol.CompletionItemKindFunction),
		InsertTextFormat: helpers.FormatPtr(protocol.InsertTextFormatSnippet),
		InsertText:       helpers.StrPtr(snippet),
	}
}

var ActionSnippet = CreateSnippet("action", "action \"${1:type}\" ${2:name} {\n\t${3}\n}", "Action block")
var StepperSnippet = CreateSnippet("steps", "steps \"${2:description}\" {\n\t${3}\n}", "Stepper block")
var TestSetSnippet = CreateSnippet("tests", "tests \"${2:description}\" {\n\t${3}\n}", "Tests set block")
var TestSnippet = CreateSnippet("test", "test \"${2:description}\"{\n\texpect {\n\t\t${3}\n\t}\n}", "Test block")
var ExpectSnippet = CreateSnippet("expect", "expect {\n\t${3}\n}", "Expect block")
var OutSnippet = CreateSnippet("out", "out {\n\t${3}\n}", "Output block")
var PrintSnippet = CreateSnippet("print", "print {\n\t${1}\n}", "Print")
var VarSnippet = CreateSnippet("var", "var ${2:name} = ${3:value}", "Variable")
var MultiVarSnippet = CreateSnippet("varm", "var (\n\t${2:var1} = ${3:val1}\n\t${4:var2} = ${5:val2}\n)", "MultiVar")
var VarSetSnippet = CreateSnippet("vars", "vars ${1:name} {\n\t${2:var1} = ${3:val1}\n\t${4:var2} = ${5:val2}\n}", "Variables set")
var SecretSetSnippet = CreateSnippet("secrets", "secrets ${1:name} {\n\t${2:set1} = ${3:val1}\n\t${4:sct2} = ${5:val2}\n}", "Secrets set")
