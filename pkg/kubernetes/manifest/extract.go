package manifest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/octohelm/cuemod/pkg/cuex"

	releasev1alpha1 "github.com/octohelm/cuemod/pkg/apis/release/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/pkg/errors"
	"github.com/stretchr/objx"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func Extract(m interface{}) (map[string]Object, error) {
	extracted := map[string]Object{}

	if err := walk(m, extracted, nil); err != nil {
		return nil, err
	}

	return extracted, nil
}

func walk(v interface{}, extracted map[string]Object, path path) error {
	switch v := v.(type) {
	case map[string]interface{}:
		return walkObj(v, extracted, path)
	case []interface{}:
		return walkList(v, extracted, path)
	}
	return errors.Errorf("unsupported %T %s", v, path)
}

func walkObj(obj objx.Map, extracted map[string]Object, p path) error {
	if isKubernetesManifest(obj) {
		gv, _ := schema.ParseGroupVersion(obj.Get("apiVersion").Str())

		if gv.Group == releasev1alpha1.SchemeGroupVersion.Group {
			switch obj.Get("kind").Str() {
			case "ReleaseTemplate":
				if obj.Get("template").IsStr() {
					var overwrites []byte

					if obj.Get("overwrites").IsStr() {
						overwrites = []byte(obj.Get("overwrites").Str())
					}

					i, err := cuex.InstanceFromTemplateAndOverwrites([]byte(obj.Get("template").Str()), overwrites)
					if err != nil {
						return err
					}
					jsonraw, err := cuex.Eval(i, cuex.JSON)
					if err != nil {
						return err
					}
					next := map[string]interface{}{}
					if err := json.NewDecoder(bytes.NewBuffer(jsonraw)).Decode(&next); err != nil {
						return err
					}
					if err := walkObj(next, extracted, p); err != nil {
						return err
					}
				}
			case "Release":
				if err := walk(obj.Get("spec").Data(), extracted, p); err != nil {
					return err
				}
			}
			return nil
		}
		co, err := ObjectFromRuntimeObject(&unstructured.Unstructured{Object: obj})
		if err != nil {
			return err
		}
		extracted[p.Full()] = co
		return nil
	}

	keys := make([]string, 0, len(obj))
	for k := range obj {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		nextP := append(p, key)
		if obj[key] == nil { // result from false if condition in Jsonnet
			continue
		}
		if err := walk(obj[key], extracted, nextP); err != nil {
			return err
		}
	}

	return nil
}

func walkList(list []interface{}, extracted map[string]Object, p path) error {
	for idx, value := range list {
		err := walk(value, extracted, append(p, idx))
		if err != nil {
			return err
		}
	}
	return nil
}

func isKubernetesManifest(obj objx.Map) bool {
	return obj.Get("apiVersion").IsStr() &&
		obj.Get("apiVersion").Str() != "" &&
		obj.Get("kind").IsStr() &&
		obj.Get("kind").Str() != ""
}

type path []interface{}

func (p path) Full() string {
	b := bytes.NewBuffer(nil)

	for _, v := range p {
		switch value := v.(type) {
		case string:
			_, _ = fmt.Fprintf(b, ".%s", value)
		case int:
			_, _ = fmt.Fprintf(b, "[%d]", value)
		}
	}

	return b.String()
}

func (p path) Base() string {
	if len(p) > 0 {
		return p[:len(p)-1].Full()
	}
	return "."
}
