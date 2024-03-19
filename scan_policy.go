package vinculum

import (
	"github.com/go-kid/ioc/scanner"
	"github.com/go-kid/ioc/scanner/meta"
	"reflect"
)

type scanRefreshScopePolicy struct {
}

func (r *scanRefreshScopePolicy) Group() meta.NodeType {
	return meta.NodeTypeComponent
}

func (r *scanRefreshScopePolicy) Tag() string {
	return Tag
}

func (r *scanRefreshScopePolicy) ExtHandler() scanner.ExtTagHandler {
	return func(field reflect.StructField, value reflect.Value) (tag string, tagVal string, ok bool) {
		if configuration, infer := value.Interface().(RefreshScopeConfiguration); infer {
			tag = Tag
			tagVal = configuration.RefreshScope()
			ok = true
		}
		return
	}
}
