package document

import "regexp"

var wordRegex = regexp.MustCompile(`[^ \t\n\f\r,;\[\]\"\']+`)

func WordAt(str string, index int) string {
	wordIdxs := wordRegex.FindAllStringIndex(str, -1)
	for _, wordIdx := range wordIdxs {
		if wordIdx[0] <= index && index <= wordIdx[1] {
			return str[wordIdx[0]:wordIdx[1]]
		}
	}

	return ""
}
