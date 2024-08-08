package server

import (
	lsfs "syfar-ls/fs"

	"syfar-ls/handlers"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	glspserv "github.com/tliron/glsp/server"
)

type Server struct {
	server      *glspserv.Server
	fileStorage *lsfs.FileStorage
}

type ServerOpts struct {
	Name    string
	Version string
	IsDebug bool
}

func NewServer(opts ServerOpts) *Server {

	handler := protocol.Handler{}

	glspServer := glspserv.NewServer(&handler, opts.Name, opts.IsDebug)
	server := &Server{
		server:      glspServer,
		fileStorage: lsfs.NewFileStorage(),
	}

	//var clientCapabilities protocol.ClientCapabilities

	handler.Initialize = func(context *glsp.Context, params *protocol.InitializeParams) (any, error) {
		//commonlog.NewInfoMessage(0, "Initializing server...")
		server.server.Log.Info("Initializing server...")
		//clientCapabilities = params.Capabilities

		capabilities := handler.CreateServerCapabilities()

		// capabilities.HoverProvider = true
		// capabilities.DefinitionProvider = true
		// capabilities.CodeActionProvider = true

		triggerChars := []string{".", "\n", "="}
		capabilities.CompletionProvider = &protocol.CompletionOptions{
			TriggerCharacters: triggerChars,
			//	ResolveProvider: boolPtr(true)
		}
		//capabilities.ReferencesProvider = &protocol.ReferenceOptions{}

		return protocol.InitializeResult{
			Capabilities: capabilities,
			ServerInfo: &protocol.InitializeResultServerInfo{
				Name:    opts.Name,
				Version: &opts.Version,
			},
		}, nil
	}

	handler.Initialized = handlers.Initialized
	handler.Shutdown = handlers.Shutdown
	handler.SetTrace = handlers.SetTrace
	handler.TextDocumentCompletion = func(context *glsp.Context, params *protocol.CompletionParams) (interface{}, error) {
		doc, err := server.fileStorage.GetDocument(params.TextDocument.URI)
		if err != nil {
			return nil, err
		}
		return handlers.TextDocumentCompletion(doc, context, params, server.server.Log)
	}

	handler.TextDocumentDidOpen = func(context *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
		doc, err := handlers.DidOpen(context, params)
		server.fileStorage.AddDocument(params.TextDocument.URI, doc)
		return err
	}

	handler.TextDocumentDidChange = func(context *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {

		doc, err := server.fileStorage.GetDocument(params.TextDocument.URI)
		if err != nil {
			return err
		}
		doc.ApplyChanges(params.ContentChanges)
		return nil
	}

	handler.TextDocumentDidClose = func(context *glsp.Context, params *protocol.DidCloseTextDocumentParams) error {
		server.fileStorage.ColseDocument(params.TextDocument.URI)
		return nil
	}

	return server
}

func (s *Server) Run() error {
	return s.server.RunStdio()
}
