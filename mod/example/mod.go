package example

import (
	"context"
	"errors"
	"fmt"
	"github.com/GoldenSheep402/Hermes/core/kernel"
	"github.com/juanjiTech/jin"
	"reflect"
	"sync"
)

var _ kernel.Module = (*Mod)(nil)

type Mod struct {
	kernel.UnimplementedModule // 请为所有Module引入UnimplementedModule
}

func (m *Mod) Name() string {
	return "example"
}

// 下面的方法皆为可选实现

func (m *Mod) PreInit(h *kernel.Hub) error {
	return nil
}

func (m *Mod) Init(h *kernel.Hub) error {
	h.Map("hello world") // 在内核注册这个依赖
	return nil
}

func (m *Mod) PostInit(h *kernel.Hub) error {
	return nil
}

func (m *Mod) Load(h *kernel.Hub) error {
	str := h.Value(reflect.TypeOf("string")).String() // 从内核获取上面注册的依赖
	fmt.Println(str)
	_, _ = h.Invoke(func(s string) { fmt.Println(s) }) // 也可以这样从内核获取上面注册的依赖
	var str2 string
	_ = h.Load(&str2) // 也可以这样从内核获取上面注册的依赖

	var http jin.Engine
	err := h.Load(&http)
	if err != nil {
		return errors.New("can't load jin from kernel")
	}
	http.GET("/ping", func(c *jin.Context) {
		_, _ = c.Writer.WriteString("pong")
	})
	return nil
}

func (m *Mod) Start(h *kernel.Hub) error {
	return nil
}

func (m *Mod) Stop(wg *sync.WaitGroup, ctx context.Context) error {
	defer wg.Done()
	return nil
}
