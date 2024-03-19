package unittest

import (
	"github.com/go-kid/ioc/configure"
	"github.com/go-kid/vinculum"
)

type tSpy struct {
	ch chan vinculum.UpdateHandler
}

func (s *tSpy) Change() <-chan vinculum.UpdateHandler {
	return s.ch
}

func (s *tSpy) Init() error {
	s.ch = make(chan vinculum.UpdateHandler, 0)
	return nil
}

func (s *tSpy) Weight() int {
	return 999
}

func (s *tSpy) Update(c []byte) {
	s.ch <- func(binder configure.Binder) error {
		return binder.SetConfig(c)
	}
}

func (s *tSpy) UpdateByPath(path string, val any) {
	s.ch <- func(binder configure.Binder) error {
		binder.Set(path, val)
		return nil
	}
}

func (s *tSpy) Close() error {
	close(s.ch)
	return nil
}
