package toml

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"cuelang.org/go/cue"
	"github.com/octohelm/cuemod/pkg/cuemod/translator"

	"github.com/google/go-jsonnet"
	tankanative "github.com/grafana/tanka/pkg/jsonnet/native"
)

func init() {
	translator.Register(t{})
}

type t struct{}

func (t) Name() string {
	return "jsonnet"
}

func (t) MarshalCueValue(value cue.Value) ([]byte, error) {
	data, err := value.MarshalJSON()
	if err != nil {
		return nil, err
	}

	jsonnetInput := JsonnetInput{}

	if err := json.Unmarshal(data, &jsonnetInput); err != nil {
		return nil, err
	}

	vm := jsonnet.MakeVM()
	funcs := tankanative.Funcs()
	for i := range funcs {
		vm.NativeFunction(funcs[i])
	}

	filename := "./main.jsonnet"

	vm.Importer(newMemoryImporter(filename, jsonnetInput.Imports))

	jsonraw, err := vm.EvaluateAnonymousSnippet(filename, string(jsonnetInput.Data))
	if err != nil {
		return nil, err
	}

	return []byte(jsonraw), nil
}

type JsonnetInput struct {
	Data    []byte                  `json:"data,omitempty"`
	Code    string                  `json:"code,omitempty"`
	Imports map[string]JsonnetInput `json:"imports"`
}

func (input *JsonnetInput) Contents() []byte {
	if len(input.Code) > 0 {
		return []byte(input.Code)
	}
	return input.Data
}

func newMemoryImporter(filename string, imports map[string]JsonnetInput) memoryImporter {
	mi := memoryImporter{}

	var makeContent func(importedAt string, imports map[string]JsonnetInput)

	makeContent = func(importedAt string, imports map[string]JsonnetInput) {
		for importPath := range imports {
			input := imports[importPath]

			nextImportedAt := mi.nextImportedAt(importedAt, importPath)

			mi[nextImportedAt] = jsonnet.MakeContents(string(input.Contents()))

			makeContent(nextImportedAt, input.Imports)
		}
	}

	makeContent(filename, imports)

	return mi
}

type memoryImporter map[string]jsonnet.Contents

func (memoryImporter) nextImportedAt(importedAt string, importPath string) string {
	return filepath.Join(filepath.Dir(importedAt), importPath)
}

func (data memoryImporter) Import(importedAt string, importPath string) (contents jsonnet.Contents, foundAt string, err error) {
	nextImportedAt := data.nextImportedAt(importedAt, importPath)
	if content, ok := data[nextImportedAt]; ok {
		return content, nextImportedAt, nil
	}
	return jsonnet.Contents{}, "", fmt.Errorf("import `%s` not available in %s", importPath, importedAt)
}
