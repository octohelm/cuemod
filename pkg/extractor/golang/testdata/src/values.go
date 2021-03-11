package src

import (
	"time"

	yamlv3 "gopkg.in/yaml.v3"
)

const (
	// test
	StringValue = "some str" // test2
	BoolValue   = true
)

const (
	IntValue int = iota + 1
	IntValue2
)

const Rename = time.ANSIC

const ScalarNode = yamlv3.ScalarNode
