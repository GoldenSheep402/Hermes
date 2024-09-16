package modList

import (
	"github.com/GoldenSheep402/Hermes/core/kernel"
	"github.com/GoldenSheep402/Hermes/mod/grpcGateway"
	"github.com/GoldenSheep402/Hermes/mod/jinPprof"
	"github.com/GoldenSheep402/Hermes/mod/jinx"
	"github.com/GoldenSheep402/Hermes/mod/pgsql"
	"github.com/GoldenSheep402/Hermes/mod/rds"
)

var ModList = []kernel.Module{
	// &b2x.Mod{},
	// &uptrace.Mod{},
	&grpcGateway.Mod{},
	&jinPprof.Mod{},
	&jinx.Mod{},
	&pgsql.Mod{},
	&rds.Mod{},
}
