package unittest

import (
	"github.com/go-kid/ioc"
	"github.com/go-kid/ioc/app"
	"github.com/go-kid/ioc/configure/loader"
	"github.com/go-kid/vinculum"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type TestAppComp struct {
	ApiConfig       *ApiConfig `refreshScope:""`
	Proxy           string     `prop:"Proxy.url" refreshScope:""`
	ApiConfigChange bool
	proxyChange     bool
}

type ApiConfig struct {
	Host string `yaml:"host"`
}

func (a *ApiConfig) Prefix() string {
	return "Api"
}

func (t *TestAppComp) OnScopeChange(path string) error {
	switch path {
	case t.ApiConfig.Prefix():
		t.ApiConfigChange = true
	case "Proxy.url":
		t.proxyChange = true
	}
	return nil
}

func initTestComp(t *testing.T, tapp any) *tSpy {
	var (
		config = []byte(`
Api:
  host: localhost:8080
Proxy:
  url: localhost:3333
`)
		spy = &tSpy{}
	)
	ioc.RunTest(t,
		app.SetConfigLoader(loader.NewRawLoader(config)),
		app.SetComponents(tapp, spy, vinculum.New()),
	)
	return spy
}

func TestRefreshComponent(t *testing.T) {
	var tapp = &TestAppComp{}
	spy := initTestComp(t, tapp)
	spy.Update([]byte(`
Api:
  host: localhost:8888
Proxy:
  url: localhost:3333
`))
	time.Sleep(time.Millisecond * 100)
	assert.Equal(t, true, tapp.ApiConfigChange)

	spy.Update([]byte(`
Api:
  host: localhost:8888
Proxy:
  url: localhost:4444
`))
	time.Sleep(time.Millisecond * 100)
	assert.Equal(t, true, tapp.proxyChange)
}
