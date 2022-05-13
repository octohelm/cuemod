package src

import (
	"sync"
	"time"
)

type Alias = Protocol

type String string

type Map map[Protocol]string

type Bytes []byte

type Slice []string

type Array [2]String

type Any interface{}

type private string

// struct comment
type Struct struct {
	Sub
	// name
	// name
	Name string `json:"name"`
}

type Sub struct {
	time.Time

	// should ignore
	sync.Mutex

	Optional *uint64 `json:"optional,omitempty"`
	Nullable *int    `json:"nullable"`
}

type Strfmt struct{}

func (Strfmt) MarshalText() (text []byte, err error) {
	return []byte("="), nil
}

func (*Strfmt) UnmarshalText(text []byte) (err error) {
	return nil
}

type YAML struct{}

func (YAML) MarshalYAML() (data interface{}, err error) {
	return []byte("="), nil
}
