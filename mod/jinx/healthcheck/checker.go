package healthcheck

import (
	"strings"
	"sync"
)

type Checker interface {
	Pass() bool
	Name() string
}

type StatusCode int

const (
	StatusOK StatusCode = iota
	StatusExcluded
	StatusError
)

var stringBuilderPool = sync.Pool{
	New: func() any {
		return &strings.Builder{}
	},
}

type Status struct {
	sync.Mutex
	m map[string]StatusCode
}

func NewStatus(n int) *Status {
	return &Status{
		m: make(map[string]StatusCode, n),
	}
}
func (s *Status) Get(k string) (StatusCode, bool) {
	s.Lock()
	defer s.Unlock()
	v, ok := s.m[k]
	return v, ok
}

func (s *Status) Set(k string, v StatusCode) {
	s.Lock()
	defer s.Unlock()
	s.m[k] = v
}

func (s *Status) Each(f func(string, StatusCode)) {
	s.Lock()
	defer s.Unlock()
	for k, v := range s.m {
		f(k, v)
	}
}

func (s *Status) String(verbose bool) string {
	allPass := true
	if verbose {
		b := stringBuilderPool.Get().(*strings.Builder)
		defer stringBuilderPool.Put(b)
		defer b.Reset()
		s.Each(func(name string, status StatusCode) {
			switch status {
			case StatusOK:
				b.WriteString("[+] " + name + " ok\n")
			case StatusError:
				b.WriteString("[-] " + name + " fail\n")
				allPass = false
			case StatusExcluded:
				b.WriteString("[+] " + name + " excluded: ok\n")
			}
		})

		if allPass {
			b.WriteString("healthz check passed")
		} else {
			b.WriteString("healthz check failed")
		}
		return b.String()
	}

	s.Each(func(name string, status StatusCode) {
		switch status {
		case StatusOK:
			allPass = false
		case StatusError:
			allPass = false
		case StatusExcluded:
			allPass = false
		}
	})

	if allPass {
		return "OK"
	}
	return "Fail"
}
