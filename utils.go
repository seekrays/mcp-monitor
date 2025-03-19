package main

import (
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

// NewToolResultJSON creates a new CallToolResult with JSON content
// It converts the result to a JSON string and wraps it in a text content
func NewToolResultJSON(v interface{}) *mcp.CallToolResult {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize JSON: %v", err))
	}
	return mcp.NewToolResultText(string(data))
} 