// Package common 通用函数
package common

import (
	"fmt"

	"github.com/bwmarrin/snowflake"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

// Initialize a Snowflake node
var node *snowflake.Node

func init() {
	var err error
	// Assuming a node number, for example, 1. In a real-world scenario, this would be a unique number for each node.
	node, err = snowflake.NewNode(1)
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize Snowflake node: %v", err))
	}
}

// GenerateSnowflakeID generates a new unique ID using the Snowflake algorithm
func GenerateSnowflakeID() string {
	return fmt.Sprintf("%d", node.Generate().Int64())
}
func GenerateNanoId(length int) (string, error) {
	const customCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	nanoid, err := gonanoid.Generate(customCharset, 21) // Generate a 21-character ID
	if err != nil {
		return "", fmt.Errorf("failed to generate NanoID: %w", err)
	}
	return nanoid, nil
}
