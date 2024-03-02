package vinculum

import (
	"github.com/go-kid/ioc/app"
)

var Refresher app.SettingOption = func(s *app.App) {
	s.Scanner.AddTags([]string{Tag})
	injectProcessor := newInjector()
	s.AddCustomizedInjectors(injectProcessor)
	s.Register(NewDistributionCenter(s, injectProcessor))
}
