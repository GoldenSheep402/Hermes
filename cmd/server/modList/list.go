package modList

import (
	"github.com/juanjiTech/jframe/core/kernel"
	"github.com/juanjiTech/jframe/mod/b2x"
	"github.com/juanjiTech/jframe/mod/grpcGateway"
	"github.com/juanjiTech/jframe/mod/jinPprof"
	"github.com/juanjiTech/jframe/mod/jinx"
	"github.com/juanjiTech/jframe/mod/myDB"
	"github.com/juanjiTech/jframe/mod/rds"
)

var ModList = []kernel.Module{
	&b2x.Mod{},
	//&uptrace.Mod{},
	&grpcGateway.Mod{},
	&jinPprof.Mod{},
	&jinx.Mod{},
	&myDB.Mod{},
	&rds.Mod{},
}
