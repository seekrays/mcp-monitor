package main

import (
	"flag"
	"fmt"
	"log"

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

	transport := flag.String("transport", "stdio", "Transport type (stdio or sse)")
	port := flag.Int("port", 8080, "TCP port for SSE transport")
	flag.Parse()

	fmt.Printf("Starting MCP System Monitor server with %s transport...\n", *transport)

	if *transport == "sse" {
		sseUrl := fmt.Sprintf("http://localhost:%d", *port)
		sseServer := server.NewSSEServer(s, server.WithBaseURL(sseUrl))
		if err := sseServer.Start(fmt.Sprintf(":%d", *port)); err != nil {
			log.Fatalf("Server error: %v", err)
		}
		log.Printf("SSE server listening on port %d", *port)
	} else {
		if err := server.ServeStdio(s); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}
}
