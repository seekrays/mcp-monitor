package cpu

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/shirou/gopsutil/v4/cpu"
)

// NewTool creates a CPU information tool
func NewTool() mcp.Tool {
	return mcp.NewTool("get_cpu_info",
		mcp.WithDescription("Get CPU information and usage"),
		mcp.WithBoolean("per_cpu",
			mcp.Description("Whether to return data for each core"),
			mcp.DefaultBool(false),
		),
	)
}

// Handler handles CPU information requests
func Handler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	perCPU, _ := request.Params.Arguments["per_cpu"].(bool)

	// Get CPU usage percentage
	percent, err := cpu.Percent(0, perCPU)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get CPU usage: %v", err)), nil
	}

	// Get CPU information
	info, err := cpu.Info()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get CPU information: %v", err)), nil
	}

	// Get CPU count information
	counts, err := cpu.Counts(true)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get CPU core count: %v", err)), nil
	}

	// Build result
	result := map[string]interface{}{
		"usage_percent": percent,
		"info":          info,
		"core_count":    counts,
	}

	// Convert result to JSON
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize result: %v", err)), nil
	}

	return mcp.NewToolResultText(string(data)), nil
} 