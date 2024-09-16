package rds

import (
	"context"
	"errors"
	"fmt"
	"github.com/juanjiTech/jframe/conf"
	"github.com/juanjiTech/jframe/core/kernel"
	"github.com/juanjiTech/jframe/mod/jinx/healthcheck"
	rds "github.com/redis/go-redis/v9"
	"time"
)

var _ kernel.Module = (*Mod)(nil)

type Mod struct {
	kernel.UnimplementedModule // 请为所有Module引入UnimplementedModule
}

func (m *Mod) Name() string {
	return "rds"
}

func (m *Mod) PreInit(hub *kernel.Hub) error {
	rdb := rds.NewClient(&rds.Options{
		Addr:     fmt.Sprintf("%s:%s", conf.Get().Redis.Addr, conf.Get().Redis.PORT),
		Password: conf.Get().Redis.PASSWORD,
		DB:       conf.Get().Redis.DB,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return err
	}

	err = healthcheck.RegisterHealthChecker(&healthChecker{rdb: rdb})
	if err != nil {
		return err
	}
	hub.Map(&rdb)
	return nil
}

func (m *Mod) Init(hub *kernel.Hub) error {
	var rdb *rds.Client
	if hub.Load(&rdb) != nil {
		return errors.New("can't load rds client from kernel")
	}

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return err
	}
	return nil
}

var _ healthcheck.Checker = (*healthChecker)(nil)

type healthChecker struct {
	rdb *rds.Client
}

func (h *healthChecker) Pass() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := h.rdb.Ping(ctx).Result()
	return err == nil
}

func (h *healthChecker) Name() string {
	return "redis"
}

//var builderPool = sync.Pool{
//	New: func() interface{} {
//		return &strings.Builder{}
//	},
//}

//这里最初是想实现一个自动为业务的key加前缀的功能，但是目前看起来有点麻烦，暂缓实现
//
//type Client struct {
//	businessPool map[string]struct{}
//	*rds.Client
//}
//
//func (c *Client) SetBusiness(business string, indexes ...string) {
//	builder := builderPool.Get().(*strings.Builder)
//	defer builderPool.Put(builder)
//	builder.WriteString("au")
//	builder.WriteString(":")
//	builder.WriteString(strings.Join([]string{"t", "test"}, ":"))
//
//	genHook := func(next rds.ProcessHook) rds.ProcessHook {
//		return func(ctx context.Context, cmd rds.Cmder) error {
//			cmd.Args() // set the key
//			_ = next(ctx, cmd)
//			return nil
//		}
//	}
//
//	rds.Command()
//	c.Client.AddHook(genHook)
//}
