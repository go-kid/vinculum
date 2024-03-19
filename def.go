package vinculum

import (
	"github.com/go-kid/ioc/configure"
	"github.com/go-kid/ioc/injector"
	"github.com/go-kid/ioc/scanner/meta"
)

const Tag = "refreshScope"

type RefreshScopeConfiguration interface {
	RefreshScope() string
}

type RefreshScopeInjector interface {
	injector.InjectProcessor
	WatchedScopes() map[string][]*meta.Node
}

type RefreshScopeComponent interface {
	OnScopeChange(path string) error
}

type Spy interface {
	Change() <-chan UpdateHandler
}

type UpdateHandler func(binder configure.Binder) error
