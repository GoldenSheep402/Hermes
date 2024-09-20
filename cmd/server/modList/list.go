package modList

import (
	"github.com/GoldenSheep402/Hermes/core/kernel"
	"github.com/GoldenSheep402/Hermes/mod/auth"
	"github.com/GoldenSheep402/Hermes/mod/casbinX"
	"github.com/GoldenSheep402/Hermes/mod/category"
	"github.com/GoldenSheep402/Hermes/mod/grpcGateway"
	"github.com/GoldenSheep402/Hermes/mod/jinPprof"
	"github.com/GoldenSheep402/Hermes/mod/jinx"
	"github.com/GoldenSheep402/Hermes/mod/pgsql"
	"github.com/GoldenSheep402/Hermes/mod/rds"
	"github.com/GoldenSheep402/Hermes/mod/torrent"
	"github.com/GoldenSheep402/Hermes/mod/user"
)

var ModList = []kernel.Module{
	&auth.Mod{},
	// &b2x.Mod{},
	&category.Mod{},
	&casbinX.Mod{},
	// &uptrace.Mod{},
	&grpcGateway.Mod{},
	&jinPprof.Mod{},
	&jinx.Mod{},
	&pgsql.Mod{},
	&rds.Mod{},
	&user.Mod{},
	&torrent.Mod{},
}
