package utils

import (
	"log"

	"github.com/bwmarrin/snowflake"
)

var snowNode *snowflake.Node

// InitSnowflakeNode 初始化雪花算法节点
func InitSnowflakeNode(nodeId int64) {

	var err error
	snowNode, err = snowflake.NewNode(nodeId)
	if err != nil {
		log.Fatalf("初始化雪花算法节点失败: %s", err.Error())
	}
	log.Println("初始化雪花算法节点成功!")
}

// GenerateID 生成一个 int64 类型的雪花ID
func GenerateID() int64 {
	return snowNode.Generate().Int64()
}
