package translator

import (
	"encoding/json"

	"cuelang.org/go/cue"
)

type Translator interface {
	// Name renderer name, be will used in attr @render()
	Name() string
	// MarshalCue marshal cue as json raw
	MarshalCueValue(v cue.Value) ([]byte, error)
}
type Translators map[string]Translator

func (m Translators) Marshal(v cue.Value) (b []byte, err error) {
	attr := v.Attribute("translate")
	if attr.Err() == nil {
		name, err := attr.String(0)
		if err != nil {
			return nil, err
		}

		if renderer, ok := m[name]; ok {
			data, err := renderer.MarshalCueValue(v)
			if err != nil {
				return nil, err
			}
			switch v.Kind() {
			// string as string
			case cue.StringKind:
				return json.Marshal(string(data))
			// bytes as bytes
			case cue.BytesKind:
				return json.Marshal(data)
			default:
				return data, nil
			}
		}
	}

	switch v.Kind() {
	case cue.StructKind:
		b = append(b, '{')
		l, _ := v.Fields()

		if l.Next() {
			for i := 0; ; i++ {
				k := l.Label()
				s, err := json.Marshal(k)
				if err != nil {
					return nil, err
				}
				b = append(b, s...)
				b = append(b, ':')

				v := l.Value()
				bb, err := m.Marshal(v)
				if err != nil {
					return nil, err
				}
				b = append(b, bb...)

				if !l.Next() {
					break
				}
				b = append(b, ',')
			}
		}

		b = append(b, '}')
		return b, nil
	case cue.ListKind:
		b = append(b, '[')
		l, _ := v.List()
		if l.Next() {
			for i := 0; ; i++ {
				x, err := m.Marshal(l.Value())
				if err != nil {
					return nil, err
				}
				b = append(b, x...)
				if !l.Next() {
					break
				}
				b = append(b, ',')
			}
		}
		b = append(b, ']')
		return b, nil
	}
	return v.MarshalJSON()
}
