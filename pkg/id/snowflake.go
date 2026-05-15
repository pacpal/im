package id

import (
	"github.com/bwmarrin/snowflake"
)

type Generator struct {
	node *snowflake.Node
}

type SnowflakeGenerator = Generator

func NewGenerator(nodeID int64) (*Generator, error) {
	node, err := snowflake.NewNode(nodeID)
	if err != nil {
		return nil, err
	}
	return &Generator{node: node}, nil
}

func NewSnowflakeGenerator(nodeID int64) *Generator {
	node, _ := snowflake.NewNode(nodeID)
	return &Generator{node: node}
}

func (g *Generator) NextID() string {
	return g.node.Generate().String()
}

func (g *Generator) Generate() string {
	return g.node.Generate().String()
}
