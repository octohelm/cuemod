package translator

import (
	"encoding/json"
	"fmt"

	"cuelang.org/go/cue"
)

func UnmarshalCueValue(cueValue cue.Value, v interface{}) error {
	data, err := cueValue.MarshalJSON()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

func ValueFromCueValue(v cue.Value) (interface{}, error) {
	switch v.Kind() {
	case cue.NullKind:
		return nil, nil
	case cue.BoolKind:
		return v.Bool()
	case cue.IntKind:
		return v.Int64()
	case cue.FloatKind:
		return v.Float64()
	case cue.BytesKind:
		return v.Bytes()
	case cue.StringKind:
		return v.String()
	case cue.StructKind:
		m := map[string]interface{}{}
		l, _ := v.Fields()

		if l.Next() {
			for i := 0; ; i++ {
				k := l.Label()
				v := l.Value()
				fieldValue, err := ValueFromCueValue(v)
				if err != nil {
					return nil, err
				}
				m[k] = fieldValue
			}
		}

		return m, nil
	case cue.ListKind:
		list := make([]interface{}, 0)
		l, _ := v.List()
		if l.Next() {
			x, err := ValueFromCueValue(l.Value())
			if err != nil {
				return nil, err
			}
			list = append(list, x)
		}
		return list, nil
	}
	return nil, fmt.Errorf("unsupported value %s: %#v", v.Kind(), v)
}
