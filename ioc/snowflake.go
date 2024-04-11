package ioc

import (
	"github.com/bwmarrin/snowflake"
	"github.com/wx-up/go-book/internal/repository/dao"
)

func CreateSnowflake() *snowflake.Node {
	node, err := snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}
	return node
}

func IDGenProvider(node *snowflake.Node) dao.IDGen {
	return func() int64 {
		return node.Generate().Int64()
	}
}
