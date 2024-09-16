package kernel

import (
	"context"
	"github.com/GoldenSheep402/Hermes/core/logx"
	"github.com/juanjiTech/inject/v2"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"reflect"
	"sync"
)

type Engine struct {
	config Config

	Ctx    context.Context
	Cancel context.CancelFunc

	inject.Injector
	modules   map[string]Module
	modulesMu sync.Mutex
}

type Config struct {
	EnableSentry bool
}

func New(config ...Config) *Engine {
	if len(config) == 0 {
		panic("config can't be empty")
	}
	return &Engine{
		config:   config[0],
		Injector: inject.New(),
		modules:  make(map[string]Module),
	}
}

func (e *Engine) Init() {
	e.Ctx, e.Cancel = context.WithCancel(context.Background())
}

func (e *Engine) StartModule() error {
	hub := Hub{
		Injector: e.Injector,
	}
	for _, module := range e.modules {
		c := module.Config()
		if c == nil {
			continue
		}
		zap.S().Info("Module " + module.Name() + " has config, try to load it")
		ct := reflect.TypeOf(c)
		if ct.Kind() != reflect.Pointer {
			zap.S().Errorf("The config exported by module %s is not a pointer.", module.Name())
		}
		if err := viper.UnmarshalKey(module.Name(), module.Config()); err != nil {
			zap.S().Error("Config Unmarshal failed: " + err.Error())
		}
	}
	for _, m := range e.modules {
		h4m := hub
		h4m.Log = logx.NameSpace("module." + m.Name())
		if err := m.PreInit(&h4m); err != nil {
			h4m.Log.Error(err)
			panic(err)
		}
	}
	for _, m := range e.modules {
		h4m := hub
		h4m.Log = logx.NameSpace("module." + m.Name())
		if err := m.Init(&h4m); err != nil {
			h4m.Log.Error(err)
			panic(err)
		}
	}
	for _, m := range e.modules {
		h4m := hub
		h4m.Log = logx.NameSpace("module." + m.Name())
		if err := m.PostInit(&h4m); err != nil {
			h4m.Log.Error(err)
			panic(err)
		}
	}
	for _, m := range e.modules {
		h4m := hub
		h4m.Log = logx.NameSpace("module." + m.Name())
		if err := m.Load(&h4m); err != nil {
			h4m.Log.Error(err)
			panic(err)
		}
	}
	for _, m := range e.modules {
		go func(m Module) {
			h4m := hub
			h4m.Log = logx.NameSpace("module." + m.Name())
			if err := m.Start(&h4m); err != nil {
				h4m.Log.Error(err)
				panic(err)
			}
		}(m)
	}
	return nil
}

func (e *Engine) Serve() {
}

func (e *Engine) Stop() error {
	wg := sync.WaitGroup{}
	wg.Add(len(e.modules))
	for _, m := range e.modules {
		err := m.Stop(&wg, e.Ctx)
		if err != nil {
			return err
		}
	}
	wg.Wait()

	return nil
}
