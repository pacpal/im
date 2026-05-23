// Package id 提供基于 snowflake 的全局唯一 ID 生成器封装。
package id

import (
	"github.com/bwmarrin/snowflake"
)

// Generator 封装了 snowflake.Node，用于生成唯一字符串 ID。
type Generator struct {
	node *snowflake.Node
}

// SnowflakeGenerator 为 Generator 的别名，向后兼容使用。
type SnowflakeGenerator = Generator

// NewGenerator 使用给定的 nodeID 创建一个 Generator（返回错误以便调用者处理）。
func NewGenerator(nodeID int64) (*Generator, error) {
	node, err := snowflake.NewNode(nodeID)
	if err != nil {
		return nil, err
	}
	return &Generator{node: node}, nil
}

// NewSnowflakeGenerator 快速创建一个 Generator，忽略创建过程中的错误（适用于非关键路径或测试）。
func NewSnowflakeGenerator(nodeID int64) *Generator {
	node, _ := snowflake.NewNode(nodeID)
	return &Generator{node: node}
}

// NextID 返回下一个唯一 ID 的字符串表示。
func (g *Generator) NextID() string {
	return g.node.Generate().String()
}

// Generate 与 NextID 等价，返回新生成的唯一字符串 ID。
func (g *Generator) Generate() string {
	return g.node.Generate().String()
}
