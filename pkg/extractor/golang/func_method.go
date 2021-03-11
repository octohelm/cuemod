package golang

import "go/types"

var stringInterfaces = map[string]method{
	"MarshalText": {
		Results: []string{"[]byte", "error"},
	},
	"UnmarshalText": {
		PtrRecv: true,
		Params:  []string{"[]byte"},
		Results: []string{"error"},
	},
}

var topInterfaces = map[string]method{
	"MarshalJSON": {
		Results: []string{"[]byte", "error"},
	},
	"UnmarshalJSON": {
		PtrRecv: true,
		Params:  []string{"[]byte"},
		Results: []string{"error"},
	},
	"MarshalYAML": {
		Results: []string{"interface{}", "error"},
	},
	"UnmarshalYAML": {
		PtrRecv: true,
		Params:  []string{"interface{}"},
		Results: []string{"error"},
	},
}

type method struct {
	PtrRecv bool
	Params  []string
	Results []string
}

func (m *method) Equal(fn *types.Func) bool {
	s := fn.Type().(*types.Signature)

	if m.PtrRecv {
		if _, ok := s.Recv().Type().(*types.Pointer); !ok {
			return false
		}
	}

	params := s.Params()
	results := s.Results()

	paramsLen := params.Len()
	resultsLen := results.Len()

	if paramsLen != len(m.Params) || resultsLen != len(m.Results) {
		return false
	}

	for i := 0; i < paramsLen; i++ {
		if params.At(i).Type().String() != m.Params[i] {
			return false
		}
	}

	for i := 0; i < resultsLen; i++ {
		if results.At(i).Type().String() != m.Results[i] {
			return false
		}
	}

	return true
}
