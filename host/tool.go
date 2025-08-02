package host

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/coreos/go-systemd/v22/login1"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/shirou/gopsutil/v4/host"
)

// NewTool creates a host information tool
func NewTool() mcp.Tool {
	return mcp.NewTool("get_host_info",
		mcp.WithDescription("Get host system information"),
	)
}

// Handler handles host information requests
func Handler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get host information
	info, err := host.Info()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get host information: %v", err)), nil
	}

	// Get user information from systemd, fallback to gopsiutils
	var users interface{}
	if systemdCon, err := login1.New(); err == nil {
		users, err = systemdCon.ListUsersContext(context.Background())
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get user information: %v", err)), nil
		}
	} else {
		users, err = host.Users()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get user information: %v", err)), nil
		}
	}
	// Get boot time
	bootTime, err := host.BootTime()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get system boot time: %v", err)), nil
	}
	bootTimeFormatted := time.Unix(int64(bootTime), 0).Format(time.RFC3339)

	// Build result
	result := map[string]interface{}{
		"info":                info,
		"users":               users,
		"boot_time":           bootTime,
		"boot_time_formatted": bootTimeFormatted,
		"uptime":              info.Uptime,
		"uptime_formatted":    formatUptime(info.Uptime),
	}

	// Convert result to JSON
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize result: %v", err)), nil
	}

	return mcp.NewToolResultText(string(data)), nil
}

// formatUptime formats uptime into human-readable format
func formatUptime(uptime uint64) string {
	days := uptime / (60 * 60 * 24)
	hours := (uptime % (60 * 60 * 24)) / (60 * 60)
	minutes := (uptime % (60 * 60)) / 60
	seconds := uptime % 60

	return fmt.Sprintf("%d days %d hours %d minutes %d seconds", days, hours, minutes, seconds)
}
