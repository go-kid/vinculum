package vinculum

import (
	"github.com/go-kid/ioc/component_definition"
	"github.com/go-kid/ioc/configure"
	"github.com/go-kid/ioc/container"
	"github.com/go-kid/ioc/container/processors"
	"github.com/go-kid/ioc/syslog"
	"github.com/go-kid/ioc/util/framework_helper"
	"github.com/go-kid/properties"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
)

type distributionCenter struct {
	processors.DefaultInstantiationAwareComponentPostProcessor
	Logger                  syslog.Logger `logger:"Vinculum"`
	tspp                    *processors.DefaultTagScanDefinitionRegistryPostProcessor
	binder                  configure.Binder
	Spies                   []Spy `wire:""`
	watchedScopes           map[string][]func() error
	scopes                  []string
	ch                      chan properties.Properties
	configurationProcessors []container.InstantiationAwareComponentPostProcessor
}

func New() any {
	return &distributionCenter{
		tspp: &processors.DefaultTagScanDefinitionRegistryPostProcessor{
			NodeType: component_definition.PropertyTypeConfiguration,
			Tag:      Tag,
			ExtractHandler: func(meta *component_definition.Meta, field *component_definition.Field) (tag, tagVal string, ok bool) {
				if _, infer := meta.Raw.(RefreshScopeConfiguration); infer {
					tag = Tag
					ok = true
				}
				return
			},
			Required: false,
		},
		watchedScopes: make(map[string][]func() error),
		ch:            make(chan properties.Properties, 1),
		configurationProcessors: framework_helper.SortOrderedComponents([]container.InstantiationAwareComponentPostProcessor{
			processors.NewConfigQuoteAwarePostProcessors(),
			processors.NewExpressionTagAwarePostProcessors(),
			processors.NewPropertiesAwarePostProcessors(),
			processors.NewValueAwarePostProcessors(),
			processors.NewValidateAwarePostProcessors(),
		}),
	}
}

func (w *distributionCenter) PostProcessComponentFactory(f container.Factory) error {
	w.binder = f.GetConfigure()
	for _, processor := range w.configurationProcessors {
		if cfp, ok := processor.(container.ComponentFactoryPostProcessor); ok {
			err := cfp.PostProcessComponentFactory(f)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (w *distributionCenter) PostProcessDefinitionRegistry(registry container.DefinitionRegistry, component any, componentName string) error {
	return w.tspp.PostProcessDefinitionRegistry(registry, component, componentName)
}

func (w *distributionCenter) PostProcessAfterInstantiation(component any, componentName string) (bool, error) {
	return true, nil
}

func (w *distributionCenter) PostProcessProperties(properties []*component_definition.Property, component any, componentName string) ([]*component_definition.Property, error) {
	rsc, ok := component.(RefreshScopeComponent)
	for _, property := range properties {
		property := property
		for path, _ := range property.Configurations {
			path := path
			w.watchedScopes[path] = append(w.watchedScopes[path], func() error {
				for _, processor := range w.configurationProcessors {
					_, err := processor.PostProcessProperties([]*component_definition.Property{property}, component, componentName)
					if err != nil {
						return err
					}
				}
				w.Logger.Debugf("refreshed configuration '%s'", property)
				if ok {
					err := rsc.OnScopeChange(path)
					if err != nil {
						return errors.Wrapf(err, "invoke %s.OnScopeChange(%s)", componentName, path)
					}
					w.Logger.Debugf("invoke '%s'.OnScopeChange(%s)", componentName, path)
				}
				return nil
			})
		}
	}
	return properties, nil
}

func (w *distributionCenter) Init() error {
	for _, spy := range w.Spies {
		spy.RegisterChannel(w.ch)
	}
	go w.refresh()
	return nil
}

func (w *distributionCenter) refresh() {
	for prop := range w.ch {
		for scope, processHandlers := range w.watchedScopes {
			newVal, ok := prop.Get(scope)
			if !ok {
				continue
			}
			originVal := w.binder.Get(scope)
			diff := cmp.Diff(originVal, newVal)
			if diff == "" {
				continue
			}
			w.Logger.Infof("identified configuration changes on scope '%s' changes:\n%s", scope, diff)
			w.binder.Set(scope, newVal)

			for _, handler := range processHandlers {
				err := handler()
				if err != nil {
					w.Logger.Panicf("refresh property config scope '%s' error: %v", scope, err)
				}
			}
		}
	}
}
