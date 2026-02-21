package modList

import (
	"github.com/GoldenSheep402/Hermes/core/kernel"
	"github.com/GoldenSheep402/Hermes/mod/auth"
	"github.com/GoldenSheep402/Hermes/mod/bonus"
	"github.com/GoldenSheep402/Hermes/mod/casbinX"
	"github.com/GoldenSheep402/Hermes/mod/category"
	"github.com/GoldenSheep402/Hermes/mod/comment"
	"github.com/GoldenSheep402/Hermes/mod/grpcGateway"
	"github.com/GoldenSheep402/Hermes/mod/invite"
	"github.com/GoldenSheep402/Hermes/mod/jinx"
	"github.com/GoldenSheep402/Hermes/mod/message"
	"github.com/GoldenSheep402/Hermes/mod/pgsql"
	"github.com/GoldenSheep402/Hermes/mod/rds"
	"github.com/GoldenSheep402/Hermes/mod/resource"
	"github.com/GoldenSheep402/Hermes/mod/system"
	"github.com/GoldenSheep402/Hermes/mod/torrent"
	"github.com/GoldenSheep402/Hermes/mod/tracker"
	"github.com/GoldenSheep402/Hermes/mod/traffic"
	"github.com/GoldenSheep402/Hermes/mod/user"
)

var ModList = []kernel.Module{
	// Infrastructure
	&pgsql.Mod{},
	&rds.Mod{},
	&jinx.Mod{},
	// &jinPprof.Mod{},
	&grpcGateway.Mod{},
	&casbinX.Mod{},
	&auth.Mod{},

	// Core business
	&user.Mod{},
	&category.Mod{},
	&torrent.Mod{},
	&resource.Mod{},
	&tracker.Mod{},
	&traffic.Mod{},
	&bonus.Mod{},
	&invite.Mod{},
	&message.Mod{},
	&comment.Mod{},
	&system.Mod{},
}
