package modList

import (
	"github.com/GoldenSheep402/Hermes/core/kernel"
	"github.com/GoldenSheep402/Hermes/mod/casbinX"
	"github.com/GoldenSheep402/Hermes/mod/category"
	"github.com/GoldenSheep402/Hermes/mod/grpcGateway"
	"github.com/GoldenSheep402/Hermes/mod/jinPprof"
	"github.com/GoldenSheep402/Hermes/mod/jinx"
	"github.com/GoldenSheep402/Hermes/mod/pgsql"
	"github.com/GoldenSheep402/Hermes/mod/rds"
	"github.com/GoldenSheep402/Hermes/mod/torrent"
)

var ModList = []kernel.Module{
	// &b2x.Mod{},
	&category.Mod{},
	&casbinX.Mod{},
	// &uptrace.Mod{},
	&grpcGateway.Mod{},
	&jinPprof.Mod{},
	&jinx.Mod{},
	&pgsql.Mod{},
	&rds.Mod{},
	&torrent.Mod{},
}
