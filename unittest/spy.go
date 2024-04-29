package unittest

import (
	"github.com/go-kid/properties"
	"gopkg.in/yaml.v3"
)

type tSpy struct {
	ch chan<- properties.Properties
}

func (s *tSpy) RegisterChannel(ch chan<- properties.Properties) {
	s.ch = ch
}

func (s *tSpy) Weight() int {
	return 999
}

func (s *tSpy) Update(c []byte) {
	p := properties.Properties{}
	err := yaml.Unmarshal(c, &p)
	if err != nil {
		panic(err)
	}
	s.ch <- p
}

func (s *tSpy) UpdateByPath(path string, val any) {
	p := properties.New()
	p.Set(path, val)
	s.ch <- p
}
