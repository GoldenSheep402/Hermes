package checker

import "github.com/GoldenSheep402/Hermes/mod/jinx/healthcheck"

var _ healthcheck.Checker = (*example)(nil)

func NewExample() healthcheck.Checker {
	return &example{}
}

type example struct {
}

func (e *example) Pass() bool {
	return true
}

func (e *example) Name() string {
	return "example"
}
