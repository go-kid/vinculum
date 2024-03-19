package vinculum

import (
	"fmt"
	"github.com/go-kid/ioc/defination"
	"github.com/go-kid/ioc/registry"
	"github.com/go-kid/ioc/scanner/meta"
)

type refreshScopeInjector struct {
	watchScopes map[string][]*meta.Node
}

func newInjector() RefreshScopeInjector {
	return &refreshScopeInjector{
		watchScopes: make(map[string][]*meta.Node),
	}
}

func (i *refreshScopeInjector) Priority() int {
	return 0
}

func (i *refreshScopeInjector) RuleName() string {
	return "Refresh_Scope_Injector"
}

func (i *refreshScopeInjector) Filter(d *meta.Node) bool {
	return d.Tag == Tag
}

func (i *refreshScopeInjector) Inject(_ registry.Registry, d *meta.Node) error {
	if d.TagVal == "" {
		configPrefix, err := extractTag(d)
		if err != nil {
			return err
		}
		d.TagVal = configPrefix
	}
	i.watchScopes[d.TagVal] = append(i.watchScopes[d.TagVal], d)
	return nil
}

func (i *refreshScopeInjector) WatchedScopes() map[string][]*meta.Node {
	return i.watchScopes
}

func extractTag(d *meta.Node) (string, error) {
	if configuration, ok := d.Value.Interface().(defination.Configuration); ok {
		return configuration.Prefix(), nil
	}
	if value, ok := d.Field.Tag.Lookup(defination.PropTag); ok {
		return value, nil
	}
	return "", fmt.Errorf("field %s is not configuration", d.ID())
}
