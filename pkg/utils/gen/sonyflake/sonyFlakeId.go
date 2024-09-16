package sonyflake

import (
	"github.com/juanjiTech/jframe/core/logx"
	"github.com/sony/sonyflake"
)

var flake *sonyflake.Sonyflake

func init() {
	flake = sonyflake.NewSonyflake(sonyflake.Settings{})
}

func GenSonyFlakeId() (int64, error) {
	id, err := flake.NextID()
	if err != nil {
		logx.NameSpace("sonyFlakeId").Warn("flake NextID failed: ", err)
		return 0, err
	}
	return int64(id), nil
}
