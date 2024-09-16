package snowflake

import (
	"github.com/bwmarrin/snowflake"
)

var Node *snowflake.Node

func init() {
	Node, _ = snowflake.NewNode(1)
	//if err != nil {
	//	panic(err)
	//}
}

func MakeUUID() string {
	return Node.Generate().String()
}
