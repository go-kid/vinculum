package vinculum

import (
	"github.com/go-kid/ioc/app"
)

var Refresher app.SettingOption = func(s *app.App) {
	s.AddScanPolicies(&scanRefreshScopePolicy{})
	injectProcessor := newInjector()
	s.AddCustomizedInjectors(injectProcessor)
	s.Register(NewDistributionCenter(s, injectProcessor))
}
