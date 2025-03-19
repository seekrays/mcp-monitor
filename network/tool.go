package network

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/shirou/gopsutil/v4/net"
)

// NewTool creates a network information tool
func NewTool() mcp.Tool {
	return mcp.NewTool("get_network_info",
		mcp.WithDescription("Get network interface and traffic information"),
		mcp.WithString("interface",
			mcp.Description("Specify the network interface name to query (optional)"),
		),
	)
}

// Handler handles network information requests
func Handler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	interfaceName, hasInterface := request.Params.Arguments["interface"].(string)

	// Get network interface information
	interfaces, err := net.Interfaces()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get network interface information: %v", err)), nil
	}

	// Get network connection count
	connectionStats, err := net.Connections("all")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get network connection information: %v", err)), nil
	}

	// Get IO counters
	ioCounters, err := net.IOCounters(true)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get network IO counters: %v", err)), nil
	}

	// Build result
	result := map[string]interface{}{
		"connection_count": len(connectionStats),
	}

	// If a specific interface is specified, only return information for that interface
	if hasInterface && interfaceName != "" {
		var selectedInterface net.InterfaceStat
		var selectedIO net.IOCountersStat
		interfaceFound := false
		ioFound := false

		// Find the specified interface
		for _, iface := range interfaces {
			if iface.Name == interfaceName {
				selectedInterface = iface
				interfaceFound = true
				break
			}
		}

		// Find IO counters for the specified interface
		for _, io := range ioCounters {
			if io.Name == interfaceName {
				selectedIO = io
				ioFound = true
				break
			}
		}

		if !interfaceFound {
			return mcp.NewToolResultError(fmt.Sprintf("Network interface not found: %s", interfaceName)), nil
		}

		result["interface"] = selectedInterface
		if ioFound {
			result["io_counter"] = selectedIO
		}
	} else {
		// Return information for all interfaces
		result["interfaces"] = interfaces
		result["io_counters"] = ioCounters
	}

	// Convert result to JSON
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize result: %v", err)), nil
	}

	return mcp.NewToolResultText(string(data)), nil
} 