package main

import (
	"fmt"

	"github.com/mark3labs/mcp-go/server"
	"github.com/seekrays/mcp-monitor/cpu"
	"github.com/seekrays/mcp-monitor/disk"
	"github.com/seekrays/mcp-monitor/host"
	"github.com/seekrays/mcp-monitor/memory"
	"github.com/seekrays/mcp-monitor/network"
	"github.com/seekrays/mcp-monitor/process"
)

func main() {
	// Create MCP server
	s := server.NewMCPServer(
		"System Monitor",
		"1.0.0",
	)

	// Add CPU tool
	s.AddTool(cpu.NewTool(), cpu.Handler)

	// Add memory tool
	s.AddTool(memory.NewTool(), memory.Handler)

	// Add disk tool
	s.AddTool(disk.NewTool(), disk.Handler)

	// Add network tool
	s.AddTool(network.NewTool(), network.Handler)

	// Add host tool
	s.AddTool(host.NewTool(), host.Handler)

	// Add process tool
	s.AddTool(process.NewTool(), process.Handler)

	// Start the stdio server
	fmt.Println("Starting MCP System Monitor server...")
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
} 