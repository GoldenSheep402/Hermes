package pyroscope

import (
	"context"
	"github.com/GoldenSheep402/Hermes/conf"
	"github.com/GoldenSheep402/Hermes/core/kernel"
	"github.com/grafana/pyroscope-go"
	"os"
	"sync"
)

var _ kernel.Module = (*Mod)(nil)

type Mod struct {
	kernel.UnimplementedModule // 请为所有Module引入UnimplementedModule
	profiler                   *pyroscope.Profiler
}

func (m *Mod) Name() string {
	return "pyroscope"
}

func (m *Mod) PreInit(hub *kernel.Hub) error {
	pyroscopeConf := conf.Get().Pyroscope
	if pyroscopeConf.ServerAddress == "" {
		hub.Log.Info("pyroscope server address is empty, skip init pyroscope")
		return nil
	}
	var err error
	m.profiler, err = pyroscope.Start(pyroscope.Config{
		ApplicationName: pyroscopeConf.ApplicationName,

		Tags: map[string]string{
			"hostname": os.Getenv("HOSTNAME"),
		},

		// replace this with the address of pyroscope server
		ServerAddress: pyroscopeConf.ServerAddress,

		// you can disable logging by setting this to nil
		Logger: nil,

		// Optional HTTP Basic authentication (Grafana Cloud)
		BasicAuthUser:     pyroscopeConf.BasicAuthUser,
		BasicAuthPassword: pyroscopeConf.BasicAuthPass,
		// Optional Pyroscope tenant ID (only needed if using multi-tenancy). Not needed for Grafana Cloud.
		TenantID: pyroscopeConf.TenantID,

		// by default all profilers are enabled,
		// but you can select the ones you want to use:
		ProfileTypes: []pyroscope.ProfileType{
			pyroscope.ProfileCPU,
			pyroscope.ProfileAllocObjects,
			pyroscope.ProfileAllocSpace,
			pyroscope.ProfileInuseObjects,
			pyroscope.ProfileInuseSpace,
		},
	})
	if err != nil {
		hub.Log.Error(err)
		return err
	}
	return nil
}

func (m *Mod) Stop(wg *sync.WaitGroup, _ context.Context) error {
	defer wg.Done()
	if m.profiler != nil {
		return m.profiler.Stop()
	}
	return nil
}
