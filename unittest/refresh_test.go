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

type TestApp struct {
	StringProp   string  `prop:"Test.refresh.stringProp" refreshScope:""`
	StructConfig *Config `refreshScope:""`
}

type TestApp2 struct {
	StringProp   string  `prop:"Test.refresh.stringProp" refreshScope:"Refresh.stringProp"`
	StructConfig *Config `refreshScope:"Refresh.structConfig"`
}

type Config struct {
	ScopeS  string   `yaml:"scopeS"`
	ScopeI  int64    `yaml:"scopeI"`
	ScopeL  []string `yaml:"scopeL"`
	ScopeST struct {
		FieldA string `yaml:"fieldA"`
	} `yaml:"scopeST"`
}

func (c *Config) Prefix() string {
	return "Test.refresh.structConfig"
}

func initTest(t *testing.T, tapp any) *tSpy {
	var (
		config = `
Test:
  refresh:
    structConfig:
      scopeS: "hello"
      scopeI: 201
      scopeL: [A,B,C,D]
      scopeST:
        fieldA: "world"
    stringProp: "foo"
`
		spy = &tSpy{}
	)
	ioc.RunTest(t,
		app.SetConfig(config),
		app.SetConfigLoader(loader.NewRawLoader()),
		app.SetComponents(tapp, spy),
		vinculum.Refresher,
	)
	return spy
}

func TestRefreshConfig(t *testing.T) {
	t.Run("TestUnSpecifyScope", func(t *testing.T) {
		t.Run("StringProp", func(t *testing.T) {
			var tapp = &TestApp{}
			spy := initTest(t, tapp)
			spy.Update([]byte(`
Test:
  refresh:
    structConfig:
      scopeS: "hello"
      scopeI: 201
      scopeL: [A,B,C,D]
      scopeST:
        fieldA: "world"
    stringProp: bar
`))
			time.Sleep(time.Millisecond * 10)
			assert.Equal(t, "bar", tapp.StringProp)
		})
		t.Run("StructConfigure", func(t *testing.T) {
			var tapp = &TestApp{}
			spy := initTest(t, tapp)
			spy.Update([]byte(`
Test:
  refresh:
    structConfig:
      scopeS: "hello1"
      scopeI: 20100
      scopeL: [B,C,D,E]
      scopeST:
        fieldA: "world1"
    stringProp: foo
`))
			time.Sleep(time.Millisecond * 10)
			assert.Equal(t, "hello1", tapp.StructConfig.ScopeS)
			assert.Equal(t, int64(20100), tapp.StructConfig.ScopeI)
			assert.Equal(t, []string{"B", "C", "D", "E"}, tapp.StructConfig.ScopeL)
			assert.Equal(t, "world1", tapp.StructConfig.ScopeST.FieldA)
		})
	})

	t.Run("TestSpecifyScope", func(t *testing.T) {
		var tapp = &TestApp2{}
		spy := initTest(t, tapp)
		spy.Update([]byte(`
Test:
  refresh:
    structConfig:
      scopeS: "hello"
      scopeI: 20
      scopeL: [A,B,C,D]
      scopeST:
        fieldA: "world"
    stringProp: foo
Refresh:
  structConfig:
    scopeS: "hello1"
    scopeI: 20100
    scopeL: [B,C,D,E]
    scopeST:
      fieldA: "world1"
  stringProp: bar
`))
		time.Sleep(time.Millisecond * 10)
		assert.Equal(t, "bar", tapp.StringProp)
		assert.Equal(t, "hello1", tapp.StructConfig.ScopeS)
		assert.Equal(t, int64(20100), tapp.StructConfig.ScopeI)
		assert.Equal(t, []string{"B", "C", "D", "E"}, tapp.StructConfig.ScopeL)
		assert.Equal(t, "world1", tapp.StructConfig.ScopeST.FieldA)
	})
}