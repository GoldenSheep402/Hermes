package jinPprof

import (
	"encoding/base64"
	"github.com/GoldenSheep402/Hermes/core/kernel"
	"github.com/juanjiTech/jin"
	"github.com/pkg/errors"
	"net/http"
	"net/http/pprof"
)

var _ kernel.Module = (*Mod)(nil)

type Mod struct {
	kernel.UnimplementedModule // 请为所有Module引入UnimplementedModule
}

func (m *Mod) Name() string {
	return "jinPprof"
}

func (m *Mod) Load(hub *kernel.Hub) error {
	var jinE jin.Engine
	err := hub.Load(&jinE)
	if err != nil {
		return errors.New("can't load jin.Engine from kernel")
	}

	authStr := "Basic " + base64.StdEncoding.EncodeToString([]byte("pprof:jframe"))
	g := jinE.Group("/debug/pprof", func(c *jin.Context) {
		auth := c.Request.URL.Query().Get("Authorization")
		if auth == "" {
			auth = c.Request.Header.Get("Authorization")
		}
		if auth != authStr {
			c.Writer.Header().Set("WWW-Authenticate", "Basic")
			c.Writer.WriteHeader(http.StatusUnauthorized)
			c.Abort()
			return
		}
		c.Next()
	})
	{
		g.GET("/", pprof.Index)
		g.GET("/cmdline", pprof.Cmdline)
		g.GET("/profile", pprof.Profile)
		g.POST("/symbol", pprof.Symbol)
		g.GET("/symbol", pprof.Symbol)
		g.GET("/trace", pprof.Trace)
		g.GET("/allocs", pprof.Handler("allocs").ServeHTTP)
		g.GET("/block", pprof.Handler("block").ServeHTTP)
		g.GET("/goroutine", pprof.Handler("goroutine").ServeHTTP)
		g.GET("/heap", pprof.Handler("heap").ServeHTTP)
		g.GET("/mutex", pprof.Handler("mutex").ServeHTTP)
		g.GET("/threadcreate", pprof.Handler("threadcreate").ServeHTTP)
	}

	return nil
}
