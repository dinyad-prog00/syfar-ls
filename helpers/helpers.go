package helpers

import protocol "github.com/tliron/glsp/protocol_3_16"

func KindPtr(kind protocol.CompletionItemKind) *protocol.CompletionItemKind {
	k := kind
	return &k
}
func FormatPtr(format protocol.InsertTextFormat) *protocol.InsertTextFormat {
	f := format
	return &f
}
func StrPtr(str string) *string {
	s := str
	return &s
}

func BoolPtr(str bool) *bool {
	s := str
	return &s
}

func IsInStringList(target string, list []string) bool {
	for _, item := range list {
		if item == target {
			return true
		}
	}
	return false
}

func GetMapAttr[T any](m map[string]interface{}, key string) (s *T, ok bool) {
	v := m[key]
	val, ok := v.(T)
	if ok {
		return &val, true
	}
	return nil, false
}
