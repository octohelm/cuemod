package a

import (
	j "encoding/json"
)

services: test: {
	_hidden: true

	metadata: annotations: configAsJson: j.Marshal({a: int: 1})
	metadata: annotations: {
		configAsToml: j.Marshal({a: int: 1}) @translate("toml")
	}
}
