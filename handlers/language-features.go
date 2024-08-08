package handlers

import (
	"syfar-ls/document"
	"syfar-ls/features/completion"
	"syfar-ls/fs"
	"syfar-ls/helpers"
	"syfar-ls/mappers"

	"github.com/tliron/commonlog"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"

	_ "github.com/tliron/commonlog/simple"
)

func Initialized(context *glsp.Context, params *protocol.InitializedParams) error {
	return nil

}

func Shutdown(context *glsp.Context) error {
	protocol.SetTraceValue(protocol.TraceValueOff)
	return nil
}

func SetTrace(context *glsp.Context, params *protocol.SetTraceParams) error {
	protocol.SetTraceValue(params.Value)
	return nil
}

func DidOpen(context *glsp.Context, params *protocol.DidOpenTextDocumentParams) (*document.Document, error) {
	langID := params.TextDocument.LanguageID
	if langID != "syfar" {
		return nil, nil
	}
	uri := params.TextDocument.URI
	return fs.NewDocument(uri, params.TextDocument.Text)
}

var Glogger commonlog.Logger

func TextDocumentCompletion(doc *document.Document, context *glsp.Context, params *protocol.CompletionParams, logger commonlog.Logger) (interface{}, error) {
	Glogger = logger
	var completionItems []protocol.CompletionItem
	if true {
		switch params.Context.TriggerKind {
		case protocol.CompletionTriggerKindTriggerCharacter:
			return completion.HandleTriggerKindTriggerCharacter(doc, context, params, logger)
		case protocol.CompletionTriggerKindInvoked:
			return completion.HandleTriggerKindTriggerInvoked(doc, context, params, logger)
		case protocol.CompletionTriggerKindTriggerForIncompleteCompletions:
			for _, word := range mappers.KeyWords {
				completionItems = append(completionItems, protocol.CompletionItem{
					Label: word,
					Kind:  helpers.KindPtr(protocol.CompletionItemKindKeyword),
				})

			}
			return completionItems, nil
		default:

			for _, word := range mappers.KeyWords {
				completionItems = append(completionItems, protocol.CompletionItem{
					Label: word,
					Kind:  helpers.KindPtr(protocol.CompletionItemKindKeyword),
				})

			}
			return completionItems, nil
		}
	}
	return completionItems, nil
}
