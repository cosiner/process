package process

import (
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
)

type SignalHandler func(sig os.Signal) (continu bool)

func SigIgnore(sig os.Signal) bool { return true }
func SigExit(sig os.Signal) bool   { return false }

type Signal struct {
	closed    int32
	closeChan chan struct{}
	c         chan os.Signal

	mu             sync.RWMutex
	defaultHandler SignalHandler
	handlers       map[string]SignalHandler
	notified       map[string]struct{}
}

func NewSignal() *Signal {
	return &Signal{
		c:              make(chan os.Signal, 1),
		closeChan:      make(chan struct{}),
		defaultHandler: SigExit,
		handlers:       make(map[string]SignalHandler),
		notified:       make(map[string]struct{}),
	}
}

func (s *Signal) Ignore(sigs ...os.Signal) *Signal {
	return s.Handle(SigIgnore, sigs...)
}

func (s *Signal) Exit(sigs ...os.Signal) *Signal {
	return s.Handle(SigExit, sigs...)
}

func (s *Signal) Default(h SignalHandler) *Signal {
	s.mu.Lock()
	s.defaultHandler = h
	s.mu.Unlock()
	return s
}

func (s *Signal) Handle(handler SignalHandler, sigs ...os.Signal) *Signal {
	s.mu.Lock()
	for _, sig := range sigs {
		name := sig.String()
		s.handlers[name] = handler

		if _, has := s.notified[name]; !has {
			s.notified[name] = struct{}{}
			signal.Notify(s.c, sig)
		}
	}
	s.mu.Unlock()
	return s
}

func (s *Signal) Close() {
	if atomic.CompareAndSwapInt32(&s.closed, 0, 1) {
		close(s.closeChan)
	}
}

func (s *Signal) handler(signal os.Signal) SignalHandler {
	name := signal.String()

	s.mu.RLock()
	handler := s.handlers[name]
	if handler == nil {
		handler = s.defaultHandler
	}
	s.mu.RUnlock()

	if handler == nil {
		handler = SigExit
	}
	return handler
}

func (s *Signal) Wait() (os.Signal, bool) {
	if atomic.LoadInt32(&s.closed) == 1 {
		return nil, false
	}
	select {
	case sig := <-s.c:
		return sig, s.handler(sig)(sig)
	case <-s.closeChan:
		return nil, false
	}
}

func (s *Signal) Loop() os.Signal {
	for {
		sig, continu := s.Wait()
		if !continu {
			return sig
		}
	}
}

func Kill(pid int, sig os.Signal) error {
	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return process.Signal(sig)
}
