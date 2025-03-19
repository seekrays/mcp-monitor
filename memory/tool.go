package memory

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/shirou/gopsutil/v4/mem"
)

// NewTool creates a memory information tool
func NewTool() mcp.Tool {
	return mcp.NewTool("get_memory_info",
		mcp.WithDescription("Get system memory usage information"),
	)
}

// Handler handles memory information requests
func Handler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get virtual memory information
	v, err := mem.VirtualMemory()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get memory information: %v", err)), nil
	}

	// Get swap memory information
	s, err := mem.SwapMemory()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get swap memory information: %v", err)), nil
	}

	// Build result
	result := map[string]interface{}{
		"virtual": map[string]interface{}{
			"total":        v.Total,
			"available":    v.Available,
			"used":         v.Used,
			"used_percent": v.UsedPercent,
			"free":         v.Free,
		},
		"swap": map[string]interface{}{
			"total":        s.Total,
			"used":         s.Used,
			"used_percent": s.UsedPercent,
			"free":         s.Free,
		},
	}

	// Convert result to JSON
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize result: %v", err)), nil
	}

	return mcp.NewToolResultText(string(data)), nil
} 