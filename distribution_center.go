package vinculum

import (
	"github.com/go-kid/ioc/configure"
	"github.com/go-kid/ioc/scanner/meta"
	"github.com/go-kid/ioc/syslog"
	"github.com/google/go-cmp/cmp"
	"github.com/samber/lo"
	"reflect"
	"sync"
)

type distributionCenter struct {
	binder        configure.Binder
	rsInjector    RefreshScopeInjector
	Spies         []Spy `wire:""`
	watchedScopes map[string][]*meta.Node
	scopes        []string
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
	w.watchedScopes = w.rsInjector.WatchedScopes()
	w.scopes = lo.Keys(w.watchedScopes)
	for _, spy := range w.Spies {
		go w.refresh(spy)
	}
	return nil
}

func (w *distributionCenter) refresh(spy Spy) {
	for handler := range spy.Change() {
		originValues := w.currentCopy()
		err := handler(w.binder)
		if err != nil {
			syslog.Panicf("refresh config error: %v", err)
		}
		wg := sync.WaitGroup{}
		wg.Add(len(originValues))
		for scope, originVal := range originValues {
			go func(scope string, originVal any) {
				defer wg.Done()
				if newVal := w.binder.Get(scope); !reflect.DeepEqual(originVal, newVal) {
					diff := cmp.Diff(originVal, newVal)
					syslog.Infof("distribution identified changes on scope '%s'\nchanges:\n%s", scope, diff)
					nodes := w.watchedScopes[scope]
					for _, node := range nodes {
						err = w.binder.PropInject([]*meta.Node{node})
						if err != nil {
							syslog.Panicf("refresh component %s config scope '%s' error: %v", node.Holder.Meta.ID(), scope, err)
						}

						if rsc, ok := node.Holder.Meta.Raw.(RefreshScopeComponent); ok {
							err := rsc.OnScopeChange(scope)
							if err != nil {
								syslog.Panicf("refresh component %s trigger OnScopeChange scope '%s' error: %v", node.Holder.Meta.ID(), scope, err)
							}
						}
					}
				}
			}(scope, originVal)
		}
		wg.Wait()
	}
}

func (w *distributionCenter) currentCopy() map[string]any {
	return cloneMap(lo.SliceToMap(w.scopes, func(scope string) (string, any) {
		return scope, w.binder.Get(scope)
	}))
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
