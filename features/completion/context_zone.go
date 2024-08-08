package completion

import (
	"syfar-ls/document"
	"syfar-ls/tmp"

	"github.com/alecthomas/participle/v2/lexer"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

type ContextZoneType int

const (
	ContextZoneTypeFile ContextZoneType = iota
	ContextZoneTypeStepper
	ContextZoneTypeAction
	ContextZoneTypeTestSet
	ContextZoneTypeTest
	ContextZoneTypeExpect
	ContextZoneTypeOut
	ContextZoneTypePrint
	ContextZoneTypeSecretSet
	ContextZoneTypeMultiVariable
	ContextZoneTypeVarSet
	ContextZoneTypeUnkown
)

type ContextZone struct {
	Type  ContextZoneType
	Props map[string]interface{}
}

func GetContextZone(doc *document.Document, pos protocol.Position) ContextZone {
	if doc.Ast != nil {
		for _, entry := range doc.Ast.Entries {
			if PosIsBetween(entry.Pos, entry.EndPos, pos) {
				switch {
				case entry.Test != nil:
					ctxZn, ok := GetContextZoneFromTest(entry.Test, pos)
					if ok {
						return ctxZn
					}
				case entry.TestSet != nil:
					ctxZn, ok := GetContextZoneFromTestSet(entry.TestSet, pos)
					if ok {
						return ctxZn
					}
				case entry.Action != nil:
					ctxZn, ok := GetContextZoneFromAction(entry.Action, pos)
					if ok {
						return ctxZn
					}
				case entry.Print != nil:
					if PosIsBetween(entry.Print.StartBlock.Pos, entry.Print.EndBlock.Pos, pos) {
						return ContextZone{Type: ContextZoneTypePrint}
					}

				}
			}
		}
	}
	return ContextZone{Type: ContextZoneTypeUnkown}
}

func GetContextZoneFromTestSet(testSet *tmp.TestSet, pos protocol.Position) (ctxZn ContextZone, ok bool) {
	for _, tt := range testSet.Tests {
		ctxZn, ok := GetContextZoneFromTest(tt, pos)
		if ok {
			return ctxZn, true
		}

	}
	if PosIsBetween(testSet.StartBlock.Pos, testSet.EndBlock.Pos, pos) {
		return ContextZone{Type: ContextZoneTypeTestSet}, ok
	}

	return ContextZone{Type: ContextZoneTypeUnkown}, false
}

func GetContextZoneFromTest(test *tmp.Test, pos protocol.Position) (ctxZn ContextZone, ok bool) {
	for _, exp := range test.Expectations {
		if PosIsBetween(exp.StartBlock.Pos, exp.EndBlock.Pos, pos) {
			return ContextZone{Type: ContextZoneTypeExpect}, true
		}
	}
	if PosIsBetween(test.StartBlock.Pos, test.EndBlock.Pos, pos) {
		return ContextZone{Type: ContextZoneTypeTest}, true
	}
	return ContextZone{Type: ContextZoneTypeUnkown}, false
}

func GetContextZoneFromAction(act *tmp.Action, pos protocol.Position) (ctxZn ContextZone, ok bool) {
	for _, at := range act.Attributes {
		switch {
		case at.TestSet != nil:
			ctxZn, ok := GetContextZoneFromTestSet(at.TestSet, pos)
			if ok {
				ctxZn.Props = map[string]interface{}{"type": act.Type}
				return ctxZn, true
			}
		case at.Test != nil:
			ctxZn, ok := GetContextZoneFromTest(at.Test, pos)
			if ok {
				ctxZn.Props = map[string]interface{}{"type": act.Type}
				return ctxZn, true
			}
		case at.Out != nil:
			if PosIsBetween(at.Out.StartBlock.Pos, at.Out.EndBlock.Pos, pos) {
				return ContextZone{Type: ContextZoneTypeOut, Props: map[string]interface{}{"type": act.Type}}, true
			}
		}

	}
	if PosIsBetween(act.StartBlock.Pos, act.EndBlock.Pos, pos) {
		attrs := tmp.GetActionsParametersName(*act)
		return ContextZone{Type: ContextZoneTypeAction, Props: map[string]interface{}{"parameters": attrs, "type": act.Type}}, true
	}

	return ContextZone{Type: ContextZoneTypeUnkown}, false
}

func PosIsBetween(start lexer.Position, end lexer.Position, pos protocol.Position) bool {
	//Glogger.Info(fmt.Sprintf("%d:%d to %d:%d --- %d:%d", start.Line, start.Column, end.Line, end.Column, pos.Line+1, pos.Character+1))
	// Si la position de départ est après la position de fin, échanger les positions
	if start.Line > end.Line || (start.Line == end.Line && start.Column > end.Column) {
		start, end = end, start
	}

	// Vérifier si la position se situe entre les positions de début et de fin
	return (int(pos.Line)+1 > start.Line || (int(pos.Line)+1 == start.Line && int(pos.Character)+1 >= start.Column)) &&
		(int(pos.Line)+1 < end.Line || (int(pos.Line)+1 == end.Line && int(pos.Character)+1 <= end.Column))
}
