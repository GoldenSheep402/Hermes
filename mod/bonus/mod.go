package bonus

import (
	"github.com/GoldenSheep402/Hermes/core/kernel"
)

var _ kernel.Module = (*Mod)(nil)

type Mod struct {
	kernel.UnimplementedModule
	config Config
}

type Config struct {
}

func (m *Mod) Config() any {
	return &m.config
}

func (m *Mod) Name() string {
	return "bonus"
}

func (m *Mod) Load(hub *kernel.Hub) error {
	return nil
}
