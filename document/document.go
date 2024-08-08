package document

import (
	"strings"

	"syfar-ls/tmp"

	protocol "github.com/tliron/glsp/protocol_3_16"
)

type Document struct {
	URI                     protocol.DocumentUri
	Filename                string
	Path                    string
	NeedsRefreshDiagnostics bool
	Content                 string
	Ast                     *tmp.SyfarFile
	Error                   bool
	lines                   []string
}

// ApplyChanges updates the content of the document from LSP textDocument/didChange events.
func (d *Document) ApplyChanges(changes []interface{}) {
	for _, change := range changes {
		switch c := change.(type) {
		case protocol.TextDocumentContentChangeEvent:
			startIndex, endIndex := c.Range.IndexesIn(d.Content)
			d.Content = d.Content[:startIndex] + c.Text + d.Content[endIndex:]
		case protocol.TextDocumentContentChangeEventWhole:
			d.Content = c.Text
		}
	}

	ast, err := tmp.ParseFile(d.Content, d.Path)
	if err == nil {
		d.Ast = ast
	}

	d.lines = nil
}

func (d *Document) GetLine(index int) (string, bool) {
	lines := d.GetLines()
	if index < 0 || index > len(lines) {
		return "", false
	}
	return lines[index], true
}

func (d *Document) GetLines() []string {
	if d.lines == nil {

		d.lines = strings.Split(d.Content, "\n")
	}
	return d.lines
}

func (d *Document) WordAt(pos protocol.Position) string {
	line, ok := d.GetLine(int(pos.Line))
	if !ok {
		return ""
	}
	return WordAt(line, int(pos.Character))
}
