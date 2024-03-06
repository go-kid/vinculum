package vinculum

import (
	"github.com/go-kid/ioc/configure"
	"github.com/go-kid/ioc/scanner/meta"
	"github.com/go-kid/ioc/syslog"
	"github.com/google/go-cmp/cmp"
	"github.com/samber/lo"
	"reflect"
)

type distributionCenter struct {
	binder     configure.Binder
	rsInjector RefreshScopeInjector
	Spy        Spy `wire:""`
	scopes     []string
}

func NewDistributionCenter(binder configure.Binder, injector RefreshScopeInjector) any {
	return &distributionCenter{
		binder:     binder,
		rsInjector: injector,
	}
}

func (w *distributionCenter) Order() int {
	return 0
}

func (w *distributionCenter) Run() error {
	watchedScopes := w.rsInjector.WatchedScopes()
	w.scopes = lo.Keys(watchedScopes)
	go func(ch <-chan []byte) {
		for changeBytes := range ch {
			syslog.Infof("config spy receive updated config:\n%s", string(changeBytes))
			originValues := cloneMap(lo.SliceToMap(w.scopes, func(scope string) (string, any) {
				return scope, w.binder.Get(scope)
			}))
			err := w.binder.SetConfig(changeBytes)
			if err != nil {
				syslog.Panicf("refresh config error: %v", err)
			}
			for scope, originVal := range originValues {
				if newVal := w.binder.Get(scope); !reflect.DeepEqual(originVal, newVal) {
					diff := cmp.Diff(originVal, newVal)
					syslog.Infof("distribution identified changes on scope '%s'\nchanges:\n%s", scope, diff)
					nodes := watchedScopes[scope]
					for _, node := range nodes {
						sourceType := node.Source.Type
						err = w.binder.PropInject([]*meta.Node{node})
						if err != nil {
							syslog.Panicf("refresh component %s config scope '%s' error: %v", sourceType.String(), scope, err)
						}

						if rsc, ok := node.Source.Value.Interface().(RefreshScopeComponent); ok {
							err := rsc.OnChange(scope)
							if err != nil {
								syslog.Panicf("refresh component %s trigger OnChange scope '%s' error: %v", sourceType.String(), scope, err)
							}
						}
					}
				}
			}
		}
	}(w.Spy.Change())
	return nil
}

func cloneMap(m map[string]any) map[string]any {
	cloneM := make(map[string]any)
	for k, v := range m {
		if mv, ok := v.(map[string]any); ok {
			cloneM[k] = cloneMap(mv)
		} else {
			cloneM[k] = v
		}
	}
	return cloneM
}
