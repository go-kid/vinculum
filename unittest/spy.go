package unittest

type tSpy struct {
	ch chan []byte
}

func (s *tSpy) Init() error {
	s.ch = make(chan []byte, 0)
	return nil
}

func (s *tSpy) Weight() int {
	return 999
}

func (s *tSpy) Change() <-chan []byte {
	return s.ch
}

func (s *tSpy) Update(c []byte) {
	s.ch <- c
}

func (s *tSpy) Close() error {
	close(s.ch)
	return nil
}
