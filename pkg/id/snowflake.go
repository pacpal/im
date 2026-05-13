package id

import (
	"github.com/bwmarrin/snowflake"
)

// Generator Snowflake ID 生成器
type Generator struct {
	node *snowflake.Node
}

// NewGenerator 创建 ID 生成器
func NewGenerator(nodeID int64) (*Generator, error) {
	node, err := snowflake.NewNode(nodeID)
	if err != nil {
		return nil, err
	}
	return &Generator{node: node}, nil
}

// NextID 生成下一个唯一 ID
func (g *Generator) NextID() string {
	return g.node.Generate().String()
}
