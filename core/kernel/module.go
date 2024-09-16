package kernel

import (
	"context"
	"github.com/juanjiTech/inject/v2"
	"go.uber.org/zap"
	"sync"
)

type Hub struct {
	inject.Injector
	Log *zap.SugaredLogger
}

type Module interface {
	Name() string
	Config() any
	PreInit(*Hub) error
	Init(*Hub) error
	PostInit(*Hub) error
	Load(*Hub) error
	Start(*Hub) error
	Stop(wg *sync.WaitGroup, ctx context.Context) error

	mustEmbedUnimplementedModule()
}

func (e *Engine) RegMod(mods ...Module) {
	e.modulesMu.Lock()
	defer e.modulesMu.Unlock()
	for _, mod := range mods {
		if mod.Name() == "" {
			panic("name of module can't be empty")
		}
		if _, ok := e.modules[mod.Name()]; ok {
			panic("module " + mod.Name() + " already exists")
		}
		e.modules[mod.Name()] = mod
	}
}

var _ Module = (*UnimplementedModule)(nil)

// UnimplementedModule 由于Module接口中的方法除Name外都是可选的，所以这里提供一个默认实现，方便开发者只实现需要的方法
type UnimplementedModule struct {
}

func (u *UnimplementedModule) Name() string {
	panic("name of module should be defined")
}

func (u *UnimplementedModule) Config() any {
	return nil
}

func (u *UnimplementedModule) PreInit(*Hub) error {
	return nil
}

func (u *UnimplementedModule) Init(*Hub) error {
	return nil
}

func (u *UnimplementedModule) PostInit(*Hub) error {
	return nil
}

func (u *UnimplementedModule) Load(*Hub) error {
	return nil
}

func (u *UnimplementedModule) Start(*Hub) error {
	return nil
}

func (u *UnimplementedModule) Stop(wg *sync.WaitGroup, _ context.Context) error {
	defer wg.Done()
	return nil
}

func (u *UnimplementedModule) mustEmbedUnimplementedModule() {}
