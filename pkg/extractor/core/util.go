package core

import (
	encodingjson "encoding/json"
	"reflect"
	"sort"
	"strings"
	"unicode"

	cueast "cuelang.org/go/cue/ast"
	cuetoken "cuelang.org/go/cue/token"
	"cuelang.org/go/encoding/json"
)

func SafeIdentifierFromImportPath(s string) string {
	parts := strings.Split(s, "/")

	lastIdx := len(parts)

	//
	for {
		lastIdx = lastIdx - 1

		if lastIdx < 0 {
			continue
		}

		last := parts[lastIdx]

		// use parent when v1 v2
		if len(last) > 2 && last[0] == 'v' && unicode.IsNumber(rune(last[1])) {
			continue
		}

		// use parent when number only
		if len(last) > 0 && unicode.IsNumber(rune(last[0])) {
			continue
		}

		runes := []rune(last)

		for i, r := range runes {
			if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
				runes[i] = '_'
			}
		}

		return string(runes)
	}
}

func Extract(v interface{}) (cueast.Expr, error) {
	data, err := encodingjson.Marshal(v)
	if err != nil {
		return nil, err
	}
	return json.Extract("", data)
}

func ExtractWithType(v interface{}) (cueast.Expr, error) {
	switch x := v.(type) {
	case map[string]interface{}:
		if len(x) == 0 {
			return cueast.NewStruct(&cueast.Field{
				Label: cueast.NewList(cueast.NewIdent("string")),
				Value: cueast.NewIdent("_"),
			}), nil
		}

		keys := make([]string, 0)
		for k := range x {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		fields := make([]interface{}, len(keys))

		for i, k := range keys {
			valueExpr, err := ExtractWithType(x[k])
			if err != nil {
				return nil, err
			}
			fields[i] = &cueast.Field{
				Label:    cueast.NewString(k),
				Token:    cuetoken.COLON,
				Optional: cuetoken.Blank.Pos(),
				Value:    valueExpr,
			}
		}

		return cueast.NewStruct(fields...), nil
	case []interface{}:
		typ := cueast.Expr(cueast.NewIdent("_"))
		if len(x) > 0 {
			t, err := ExtractWithType(x[0])
			if err != nil {
				return nil, err
			}
			typ = t
		}
		return cueast.NewList(&cueast.Ellipsis{Type: typ}), nil
	case nil:
		return cueast.NewIdent("_"), nil
	default:
		d, _ := encodingjson.Marshal(v)
		expr, err := json.Extract("", d)
		if err != nil {
			return nil, err
		}
		return defaultValueAndType(
			expr,
			cueast.NewIdent(reflect.TypeOf(v).String()),
		), nil
	}
}

func defaultValueAndType(defaultValue cueast.Expr, t cueast.Expr) cueast.Expr {
	return cueast.NewBinExpr(cuetoken.OR, &cueast.UnaryExpr{Op: cuetoken.MUL, X: defaultValue}, t)
}

func NewBytes(data []byte) *cueast.BasicLit {
	return &cueast.BasicLit{
		Kind:     cuetoken.STRING,
		ValuePos: cuetoken.NoPos,
		Value:    "'''\n" + strings.Replace(string(data), "\\", "\\\\", -1) + "'''",
	}
}
