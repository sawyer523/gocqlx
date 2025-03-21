package qb

import (
	"strings"
	"unicode"
)

// camelToSnakeASCII converts camel case strings to snake case. For performance
// reasons it only works with ASCII strings.
func camelToSnakeASCII(s string) string {
	buf := []byte(s)
	out := make([]byte, 0, len(buf)+3)

	l := len(buf)
	for i := 0; i < l; i++ {
		if !(allowedBindRune(buf[i]) || buf[i] == '_') {
			continue
		}

		b := rune(buf[i])

		if unicode.IsUpper(b) {
			if i > 0 &&
				buf[i-1] != '_' &&
				buf[i-1] != '"' &&
				(unicode.IsLower(rune(buf[i-1])) || (i+1 < l && unicode.IsLower(rune(buf[i+1])))) {
				out = append(out, '_')
			}
			b = unicode.ToLower(b)
		}

		out = append(out, byte(b))
	}

	return string(out)
}

func allowedBindRune(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9')
}

func hasUpper(s string) bool {
	if strings.ContainsAny(s, " ") {
		return false
	}
	for i := 0; i < len(s); i++ {
		if unicode.IsUpper(rune(s[i])) {
			return true
		}
	}
	return false
}

func quote(s string) string {
	if len(s) == 0 {
		return s
	}
	if !hasUpper(s) {
		if _, ok := keywords[s]; ok {
			return "\"" + s + "\""
		}
		return s
	}
	if s[0] == '"' && s[len(s)-1] == '"' {
		return s
	}
	return "\"" + s + "\""
}

var keywords = map[string]struct{}{
	"add":          {},
	"allow":        {},
	"alter":        {},
	"and":          {},
	"apply":        {},
	"asc":          {},
	"authorize":    {},
	"batch":        {},
	"begin":        {},
	"by":           {},
	"columnfamily": {},
	"create":       {},
	"delete":       {},
	"desc":         {},
	"describe":     {},
	"drop":         {},
	"entries":      {},
	"execute":      {},
	"from":         {},
	"full":         {},
	"grant":        {},
	"if":           {},
	"in":           {},
	"index":        {},
	"infinity":     {},
	"insert":       {},
	"into":         {},
	"keyspace":     {},
	"limit":        {},
	"modify":       {},
	"nan":          {},
	"norecursive":  {},
	"not":          {},
	"null":         {},
	"of":           {},
	"on":           {},
	"or":           {},
	"order":        {},
	"primary":      {},
	"rename":       {},
	"replace":      {},
	"revoke":       {},
	"schema":       {},
	"select":       {},
	"set":          {},
	"table":        {},
	"to":           {},
	"token":        {},
	"truncate":     {},
	"unlogged":     {},
	"update":       {},
	"use":          {},
	"using":        {},
	"view":         {},
	"where":        {},
	"with":         {},
}
