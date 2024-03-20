package vinculum

import (
	"fmt"
	"github.com/go-kid/ioc/registry"
	"github.com/go-kid/ioc/scanner/meta"
	"github.com/go-kid/ioc/syslog"
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
	err := bindTagVal(d)
	if err != nil {
		return err
	}
	i.watchScopes[d.TagVal] = append(i.watchScopes[d.TagVal], d)
	syslog.Tracef("vinculum watched config: %s", d.TagVal)
	return nil
}

func (i *refreshScopeInjector) WatchedScopes() map[string][]*meta.Node {
	return i.watchScopes
}

func bindTagVal(d *meta.Node) error {
	if d.TagVal != "" {
		return nil
	}
	nodes := d.Holder.Meta.GetConfigurationNodes()
	if len(nodes) == 0 {
		return fmt.Errorf("field %s is not configuration", d.ID())
	}
	d.TagVal = nodes[0].TagVal
	return nil
}
