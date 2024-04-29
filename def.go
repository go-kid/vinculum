package vinculum

import (
	"github.com/go-kid/properties"
)

const Tag = "refreshScope"

type RefreshScopeConfiguration interface {
	RefreshScope()
}

type RefreshScopeComponent interface {
	OnScopeChange(path string) error
}

type Spy interface {
	RegisterChannel(ch chan<- properties.Properties)
}
