package completion

import (
	"fmt"
	"strings"
	"syfar-ls/document"
	"syfar-ls/tmp"

	"syfar-ls/helpers"

	"github.com/tliron/commonlog"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func IsVariableReference(word string) bool {
	parts := strings.Split(word, ".")
	return len(parts) == 2 && parts[0] == "var"
}

func IsVariableSetReference(word string) bool {
	parts := strings.Split(word, ".")
	return len(parts) == 2 && parts[0] == "vars"
}

func IsVariableSetAttrReference(word string) bool {
	parts := strings.Split(word, ".")
	return len(parts) == 3 && parts[0] == "vars"

}

func IsSecretSetReference(word string) bool {
	parts := strings.Split(word, ".")
	return len(parts) == 2 && parts[0] == "secrets"
}

func IsSecretSetAttrReference(word string) bool {
	parts := strings.Split(word, ".")
	return len(parts) == 3 && parts[0] == "secrets"
}

func GetSetID(word string) string {
	parts := strings.Split(word, ".")
	if len(parts) == 3 {
		return parts[1]
	}
	return ""
}

func IsActionResultReference(word string) bool {
	parts := strings.Split(word, ".")
	return len(parts) == 2 && parts[0] == "r"
}

func HandleTriggerKindTriggerCharacter(doc *document.Document, context *glsp.Context, params *protocol.CompletionParams, logger commonlog.Logger) (interface{}, error) {
	if params.Context.TriggerCharacter != nil {
		switch *params.Context.TriggerCharacter {
		case ".":
			return HandlePointTriggerCharacter(doc, context, params, logger)
		case "\n":
			ctxZone := GetContextZone(doc, params.Position)
			switch ctxZone.Type {
			case ContextZoneTypeTestSet:
				return []protocol.CompletionItem{{
					Label: "test",
					Kind:  helpers.KindPtr(protocol.CompletionItemKindKeyword),
				}}, nil
			case ContextZoneTypeTest:
				return []protocol.CompletionItem{{
					Label: "expect",
					Kind:  helpers.KindPtr(protocol.CompletionItemKindKeyword),
				}}, nil
			case ContextZoneTypePrint:
				return BuildCompletion(BuildVarAndSecret(doc), protocol.CompletionItemKindVariable), nil
			case ContextZoneTypeAction:
				return BuildActionCompletion(ctxZone), nil
			}

		case "=":
			return BuildCompletion(BuildVarAndSecret(doc), protocol.CompletionItemKindVariable), nil
		}
	}

	return nil, nil
}

func HandleTriggerKindTriggerInvoked(doc *document.Document, context *glsp.Context, params *protocol.CompletionParams, logger commonlog.Logger) (interface{}, error) {
	ctxZone := GetContextZone(doc, params.Position)
	switch ctxZone.Type {
	case ContextZoneTypePrint:
		return BuildCompletion(BuildVarAndSecret(doc), protocol.CompletionItemKindVariable), nil
	case ContextZoneTypeAction:
		return BuildActionCompletion(ctxZone), nil
	case ContextZoneTypeTestSet:
		return []protocol.CompletionItem{TestSetSnippet}, nil
	case ContextZoneTypeTest:
		return []protocol.CompletionItem{ExpectSnippet}, nil
	case ContextZoneTypeUnkown:
		return []protocol.CompletionItem{
			ActionSnippet,
			TestSetSnippet,
			TestSnippet,
			VarSetSnippet,
			VarSnippet,
			MultiVarSnippet,
			PrintSnippet,
			StepperSnippet,
			SecretSetSnippet,
			ExpectSnippet,
		}, nil
	}
	return nil, nil
}

func HandlePointTriggerCharacter(doc *document.Document, context *glsp.Context, params *protocol.CompletionParams, logger commonlog.Logger) (interface{}, error) {
	var completionItems []protocol.CompletionItem
	word := doc.WordAt(params.Position)
	switch {
	case IsVariableReference(word):
		logger.Info(word)
		if doc.Ast != nil {
			vars := tmp.GetVariableList(doc.Ast)
			for _, v := range vars {
				completionItems = append(completionItems, protocol.CompletionItem{
					Label: v.Name,
					Kind:  helpers.KindPtr(protocol.CompletionItemKindProperty),
				})
			}
			//multi
			mvars := tmp.GetMultiVarList(doc.Ast)
			for _, mv := range mvars {
				for _, v := range mv.Variables {
					completionItems = append(completionItems, protocol.CompletionItem{
						Label: v.Name,
						Kind:  helpers.KindPtr(protocol.CompletionItemKindProperty),
					})
				}
			}

		}
		return completionItems, nil

	case IsVariableSetReference(word):
		if doc.Ast != nil {
			vars := tmp.GetVarSetList(doc.Ast)
			for _, v := range vars {
				completionItems = append(completionItems, protocol.CompletionItem{
					Label: v.Id,
					Kind:  helpers.KindPtr(protocol.CompletionItemKindProperty),
				})
			}

		}
		return completionItems, nil
	case IsVariableSetAttrReference(word):
		id := GetSetID(word)
		varSet, ok := tmp.GetVarSet(doc.Ast, id)
		if ok {
			for _, v := range varSet.Variables {
				completionItems = append(completionItems, protocol.CompletionItem{
					Label: v.Name,
					Kind:  helpers.KindPtr(protocol.CompletionItemKindProperty),
				})
			}
		}
		return completionItems, nil
	case IsSecretSetReference(word):
		if doc.Ast != nil {
			vars := tmp.GetSecretSetList(doc.Ast)
			for _, v := range vars {
				completionItems = append(completionItems, protocol.CompletionItem{
					Label:      v.Id,
					InsertText: helpers.StrPtr(v.Id),
					Kind:       helpers.KindPtr(protocol.CompletionItemKindProperty),
				})
			}

		}
		return completionItems, nil
	case IsSecretSetAttrReference(word):
		id := GetSetID(word)
		sSet, ok := tmp.GetSecretSet(doc.Ast, id)
		if ok {
			for _, v := range sSet.Variables {
				completionItems = append(completionItems, protocol.CompletionItem{
					Label: v.Name,
					Kind:  helpers.KindPtr(protocol.CompletionItemKindProperty),
				})
			}
		}
		return completionItems, nil

	case IsActionResultReference(word):
		ctxZone := GetContextZone(doc, params.Position)
		switch ctxZone.Type {
		case ContextZoneTypeExpect, ContextZoneTypeOut:
			return BuildOutCompletion(ctxZone), nil
		}

		return completionItems, nil

	default:
		return completionItems, nil
	}
}

func BuildCompletion(labels []string, kink protocol.CompletionItemKind) []protocol.CompletionItem {
	var completionItems []protocol.CompletionItem

	for _, l := range labels {
		completionItems = append(completionItems, protocol.CompletionItem{Label: l, Kind: &kink})
	}

	return completionItems
}

func BuildVarAndSecret(doc *document.Document) []string {

	list := MapListTo(tmp.GetVariableList(doc.Ast), func(e tmp.Variable) string {
		return fmt.Sprintf("var.%s", e.Name)
	})
	mvars := tmp.GetMultiVarList(doc.Ast)
	for _, mv := range mvars {
		for _, v := range mv.Variables {
			list = append(list, fmt.Sprintf("var.%s", v.Name))
		}
	}
	setvars := tmp.GetVarSetList(doc.Ast)
	for _, sv := range setvars {
		list = append(list, fmt.Sprintf("vars.%s", sv.Id))
		l := MapListTo(sv.Variables, func(e *tmp.Assignment) string {
			return fmt.Sprintf("vars.%s.%s", sv.Id, e.Name)
		})
		list = append(list, l...)
	}
	svars := tmp.GetSecretSetList(doc.Ast)
	for _, sv := range svars {
		list = append(list, fmt.Sprintf("secrets.%s", sv.Id))
		l := MapListTo(sv.Variables, func(e *tmp.Assignment) string {
			return fmt.Sprintf("secrets.%s.%s", sv.Id, e.Name)
		})
		list = append(list, l...)
	}

	acts := tmp.GetActionList(doc.Ast)
	for _, act := range acts {
		for _, at := range act.Attributes {
			if at.Out != nil {
				l := MapListTo(at.Out.Variables, func(e *tmp.OutAssignment) string {
					return fmt.Sprintf("%s.%s", act.Id, e.Name)
				})
				list = append(list, l...)
			}
		}

	}

	return list
}

func MapListTo[I any, R any](list []I, mapper func(elem I) R) []R {
	result := []R{}
	for _, elem := range list {
		result = append(result, mapper(elem))
	}
	return result
}

func BuildActionCompletion(ctxZn ContextZone) []protocol.CompletionItem {
	var completionItems []protocol.CompletionItem

	setupParams, ok := helpers.GetMapAttr[[]string](ctxZn.Props, "parameters")

	if !ok {
		return completionItems
	}
	ptype, ok := helpers.GetMapAttr[string](ctxZn.Props, "type")
	if !ok {
		return completionItems
	}

	for _, att := range tmp.Parameters[*ptype] {
		if !helpers.IsInStringList(att.Name, *setupParams) {
			completionItems = append(completionItems, protocol.CompletionItem{
				Label:      att.Name,
				Kind:       helpers.KindPtr(protocol.CompletionItemKindProperty),
				InsertText: helpers.StrPtr(fmt.Sprintf("%s = ", att.Name)),
				Detail:     helpers.StrPtr(att.Type.String()),
				Preselect:  helpers.BoolPtr(true),
				SortText:   helpers.StrPtr("______________________"),
				//Documentation: helpers.StrPtr("Docuements azertyu"),
			})
		}
	}

	completionItems = append(completionItems, OutSnippet)
	completionItems = append(completionItems, TestSetSnippet)
	completionItems = append(completionItems, TestSnippet)

	return completionItems
}

func BuildOutCompletion(ctxZn ContextZone) []protocol.CompletionItem {
	var completionItems []protocol.CompletionItem
	ptype, ok := helpers.GetMapAttr[string](ctxZn.Props, "type")
	if !ok {
		return completionItems
	}

	for _, att := range tmp.Outputs[*ptype] {

		completionItems = append(completionItems, protocol.CompletionItem{
			Label:      att.Name,
			Kind:       helpers.KindPtr(protocol.CompletionItemKindProperty),
			InsertText: helpers.StrPtr(att.Name),
			Detail:     helpers.StrPtr(att.Type.String()),
			Preselect:  helpers.BoolPtr(true),
		})

	}
	return completionItems
}
