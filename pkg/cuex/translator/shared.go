package translator

import "cuelang.org/go/cue"

var shared = Translators{}

func Register(m Translator) {
	shared[m.Name()] = m
}

func MarshalCueValue(v cue.Value) (b []byte, err error) {
	return shared.Marshal(v)
}
