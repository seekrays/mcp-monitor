package disk

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/shirou/gopsutil/v4/disk"
)

// NewTool creates a disk information tool
func NewTool() mcp.Tool {
	return mcp.NewTool("get_disk_info",
		mcp.WithDescription("Get disk usage information"),
		mcp.WithString("path",
			mcp.Description("Specify the disk path to query"),
			mcp.DefaultString("/"),
		),
		mcp.WithBoolean("all_partitions",
			mcp.Description("Whether to return information for all partitions"),
			mcp.DefaultBool(false),
		),
	)
}

// Handler handles disk information requests
func Handler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path, _ := request.Params.Arguments["path"].(string)
	allPartitions, _ := request.Params.Arguments["all_partitions"].(bool)

	// Get disk usage information
	usage, err := disk.Usage(path)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get disk usage information: %v", err)), nil
	}

	result := map[string]interface{}{
		"path":  path,
		"usage": usage,
	}

	// If all partitions information is requested
	if allPartitions {
		partitions, err := disk.Partitions(true)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get partition information: %v", err)), nil
		}
		result["partitions"] = partitions
	}

	// Get IO counters
	ioCounters, err := disk.IOCounters()
	if err == nil { // Only add if successful
		result["io_counters"] = ioCounters
	}

	// Convert result to JSON
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize result: %v", err)), nil
	}

	return mcp.NewToolResultText(string(data)), nil
} 