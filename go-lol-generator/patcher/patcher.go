package patcher

import (
	"fmt"
	"go/token"
	"go/types"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/pkg/errors"
)

type ErrPatchRequired struct {
	ResID string
	For   string
}

func patchRequired(resID, format string, args ...interface{}) *ErrPatchRequired {
	return &ErrPatchRequired{
		ResID: resID,
		For:   fmt.Sprintf(format, args...),
	}
}
func (e ErrPatchRequired) Error() string {
	if e.For == "" {
		return fmt.Sprintf("patcher: override required for resource %q\n", e.ResID)
	}
	return fmt.Sprintf("patcher: override required for %s in resource %q", e.For, e.ResID)
}

func forResource(id string) (*ResPatch, error) {
	rp, ok := overrides.Resources[id]
	if !ok {
		return nil, patchRequired(id, "")
	}
	return &rp, nil
}

func ForClass(resID, clsName string) (*ClassPatch, error) {
	rp, err := forResource(resID)
	if err != nil {
		return nil, err
	}
	cp, err := rp.Class(clsName)
	if err != nil {
		return nil, err
	}
	return cp, nil
}

func StructName(resID, clsName string) (string, error) {
	cp, err := ForClass(resID, clsName)
	if err != nil {
		return "", err
	}
	return cp.Name, nil
}

func ForOperation(resID, opPath string) (*OpPatch, error) {
	rp, err := forResource(resID)
	if err != nil {
		return nil, err
	}
	for pathSuffix, pp := range rp.Operations {
		if strings.HasSuffix(opPath, pathSuffix) {
			return &pp, nil
		}
	}
	return nil, patchRequired(resID, "operation %q", opPath)
}

func OperationName(resID, opPath string) (string, error) {
	op, err := ForOperation(resID, opPath)
	if err != nil {
		return "", err
	}
	return op.Name, nil
}

var javaPrimitives = map[string]types.BasicKind{
	"boolean": types.Bool,
	"int":     types.Int32,
	"long":    types.Int64,
	"string":  types.String,
	"double":  types.Float64,
	"float":   types.Float32,
}

// Type converts types used in riot document to golang type.
func Type(resID, s string) (types.Type, error) {
	s = strings.TrimSpace(s)

	if bk, ok := javaPrimitives[s]; ok {
		return types.Typ[bk], nil
	}

	if s == "" {
		return nil, errors.Errorf("cannot parse empty string as types.Type\n")
	}

	if strings.HasPrefix(s, "List[") && strings.HasSuffix(s, "]") {
		elem, err := Type(resID, s[5:len(s)-1])
		if err != nil {
			return nil, errors.Wrapf(err, "invalid list type: %q\n", s)
		}

		return types.NewSlice(elem), nil
	}

	if strings.HasPrefix(s, "Set[") && strings.HasSuffix(s, "]") {
		elem, err := Type(resID, s[4:len(s)-1])
		if err != nil {
			return nil, errors.Wrapf(err, "invalid set type: %q\n", s)
		}

		return types.NewSlice(elem), nil
	}

	if strings.HasPrefix(s, "Map[") && strings.HasSuffix(s, "]") {
		bs := strings.Split(s[4:len(s)-1], ", ")
		if len(bs) != 2 {
			return nil, errors.Errorf("invalid map type: %q\n", s)
		}

		key, err := Type(resID, bs[0])
		if err != nil {
			return nil, errors.Wrap(err, "invalid map key type\n")
		}

		elem, err := Type(resID, bs[1])
		if err != nil {
			return nil, errors.Wrap(err, "invalid map elem type\n")
		}
		return types.NewMap(key, elem), nil
	}

	return ClassType(resID, s)
}

func ClassType(resID, clsName string) (types.Type, error) {
	switch resID {
	case "lol-static-data":
		if clsName == "SpellRange" {
			goto ret
		}
	}

	{
		cp, err := ForClass(resID, clsName)
		if err != nil {
			return nil, errors.Wrapf(err, "unknown type %q\n", clsName)
		}
		clsName = cp.Name
		goto ret
	}

ret:
	return types.NewPointer(types.NewNamed(types.NewTypeName(token.NoPos, nil, clsName, nil), nil, nil)), nil
}

func PathParamType(paramName, typeStr string) string {
	switch paramName {
	case "summonerIds":
		return "List[long]"
	case "summonerNames":
		return "List[string]"
	}
	return typeStr
}

func FieldTypeString(resID, clsName, rawFieldName, typeStr string) string {
	if resID == "lol-static-data" {
		switch clsName {
		case "SummonerSpellDto", "ChampionSpellDto":
			switch rawFieldName {
			case "effect": // typeStr == List[object]
				return "List[List[double]]"
			case "range":
				return "SpellRange"
			}
		}
	}

	return typeStr
}

func FieldName(orig string) string {
	if orig == "" {
		return ""
	}

	name := upperFirst(orig)

	return lintName(name)
}

func upperFirst(s string) string {
	if s == "" {
		return ""
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[n:]
}

// lintName returns a different name if it should be different.
func lintName(name string) (should string) {
	// Fast path for simple cases: "_" and all lowercase.
	if name == "_" {
		return name
	}
	allLower := true
	for _, r := range name {
		if !unicode.IsLower(r) {
			allLower = false
			break
		}
	}
	if allLower {
		return name
	}

	// Split camelCase at any lower->upper transition, and split on underscores.
	// Check each word for common initialisms.
	runes := []rune(name)
	w, i := 0, 0 // index of start of word, scan
	for i+1 <= len(runes) {
		eow := false // whether we hit the end of a word
		if i+1 == len(runes) {
			eow = true
		} else if runes[i+1] == '_' {
			// underscore; shift the remainder forward over any run of underscores
			eow = true
			n := 1
			for i+n+1 < len(runes) && runes[i+n+1] == '_' {
				n++
			}

			// Leave at most one underscore if the underscore is between two digits
			if i+n+1 < len(runes) && unicode.IsDigit(runes[i]) && unicode.IsDigit(runes[i+n+1]) {
				n--
			}

			copy(runes[i+1:], runes[i+n+1:])
			runes = runes[:len(runes)-n]
		} else if unicode.IsLower(runes[i]) && !unicode.IsLower(runes[i+1]) {
			// lower->non-lower
			eow = true
		}
		i++
		if !eow {
			continue
		}

		// [w,i) is a word.
		word := string(runes[w:i])
		if u := strings.ToUpper(word); commonInitialisms[u] {
			// Keep consistent case, which is lowercase only at the start.
			if w == 0 && unicode.IsLower(runes[w]) {
				u = strings.ToLower(u)
			}
			// All the common initialisms are ASCII,
			// so we can replace the bytes exactly.
			copy(runes[w:], []rune(u))
		} else if w > 0 && strings.ToLower(word) == word {
			// already all lowercase, and not the first word, so uppercase the first character.
			runes[w] = unicode.ToUpper(runes[w])
		}
		w = i
	}
	return string(runes)
}

// commonInitialisms is a set of common initialisms.
// Only add entries that are highly unlikely to be non-initialisms.
// For instance, "ID" is fine (Freudian code is rare), but "AND" is not.
var commonInitialisms = map[string]bool{
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CSS":   true,
	"DNS":   true,
	"EOF":   true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IP":    true,
	"JSON":  true,
	"LHS":   true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SQL":   true,
	"SSH":   true,
	"TCP":   true,
	"TLS":   true,
	"TTL":   true,
	"UDP":   true,
	"UI":    true,
	"UID":   true,
	"UUID":  true,
	"URI":   true,
	"URL":   true,
	"UTF8":  true,
	"VM":    true,
	"XML":   true,
	"XSRF":  true,
	"XSS":   true,
}
