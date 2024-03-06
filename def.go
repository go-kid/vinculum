package vinculum

import (
	"github.com/go-kid/ioc/injector"
	"github.com/go-kid/ioc/scanner/meta"
)

const Tag = "refreshScope"

type Spy interface {
	Change() <-chan []byte
}

type RefreshScopeInjector interface {
	injector.InjectProcessor
	WatchedScopes() map[string][]*meta.Node
}

type RefreshScopeComponent interface {
	OnChange(path string) error
}
