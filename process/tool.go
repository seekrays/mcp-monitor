package process

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/shirou/gopsutil/v4/process"
)

// NewTool creates a process information tool
func NewTool() mcp.Tool {
	return mcp.NewTool("get_process_info",
		mcp.WithDescription("Get process information"),
		mcp.WithNumber("pid",
			mcp.Description("Process ID. If not specified, returns summary information for all processes"),
		),
		mcp.WithNumber("limit",
			mcp.Description("Limit the number of processes returned"),
			mcp.DefaultNumber(10),
		),
		mcp.WithString("sort_by",
			mcp.Description("Sort field (cpu, memory, pid, name)"),
			mcp.DefaultString("cpu"),
		),
	)
}

// Handler handles process information requests
func Handler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	pidRaw, hasPID := request.Params.Arguments["pid"]
	limit, _ := request.Params.Arguments["limit"].(float64)
	sortBy, _ := request.Params.Arguments["sort_by"].(string)

	// If limit is not specified or is less than or equal to 0, use default value 10
	if limit <= 0 {
		limit = 10
	}

	// If PID is specified, return detailed information for that process
	if hasPID {
		// Convert PID type
		var pid int32
		switch v := pidRaw.(type) {
		case float64:
			pid = int32(v)
		case int:
			pid = int32(v)
		default:
			return mcp.NewToolResultError("PID must be a number"), nil
		}

		proc, err := process.NewProcess(pid)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get process information: %v", err)), nil
		}

		// Get detailed process information
		name, _ := proc.Name()
		cmdline, _ := proc.Cmdline()
		createTime, _ := proc.CreateTime()
		status, _ := proc.Status()
		memInfo, _ := proc.MemoryInfo()
		cpuPercent, _ := proc.CPUPercent()
		username, _ := proc.Username()
		numThreads, _ := proc.NumThreads()
		ioCounters, _ := proc.IOCounters()

		result := map[string]interface{}{
			"pid":         pid,
			"name":        name,
			"cmdline":     cmdline,
			"create_time": createTime,
			"status":      status,
			"cpu_percent": cpuPercent,
			"username":    username,
			"num_threads": numThreads,
		}

		if memInfo != nil {
			result["memory"] = map[string]interface{}{
				"rss":  memInfo.RSS,  // Resident Set Size
				"vms":  memInfo.VMS,  // Virtual Memory Size
				"swap": memInfo.Swap, // Swap space
			}
		}

		if ioCounters != nil {
			result["io_counters"] = ioCounters
		}

		// Convert result to JSON
		data, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize result: %v", err)), nil
		}

		return mcp.NewToolResultText(string(data)), nil
	}

	// Get all processes
	processes, err := process.Processes()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get process list: %v", err)), nil
	}

	// Build process information list
	var processList []map[string]interface{}
	for _, proc := range processes {
		pid := proc.Pid
		name, err := proc.Name()
		if err != nil {
			name = "unknown"
		}

		cpuPercent, err := proc.CPUPercent()
		if err != nil {
			cpuPercent = 0.0
		}

		memInfo, err := proc.MemoryInfo()
		var memPercent float64
		if err == nil && memInfo != nil {
			// Explicitly convert float32 to float64
			memPercentFloat32, _ := proc.MemoryPercent()
			memPercent = float64(memPercentFloat32)
		}

		processList = append(processList, map[string]interface{}{
			"pid":          pid,
			"name":         name,
			"cpu_percent":  cpuPercent,
			"mem_percent":  memPercent,
			"memory_bytes": memInfo,
		})
	}

	// Sort by specified field
	sort.Slice(processList, func(i, j int) bool {
		switch sortBy {
		case "cpu":
			return processList[i]["cpu_percent"].(float64) > processList[j]["cpu_percent"].(float64)
		case "memory":
			return processList[i]["mem_percent"].(float64) > processList[j]["mem_percent"].(float64)
		case "pid":
			return processList[i]["pid"].(int32) < processList[j]["pid"].(int32)
		case "name":
			return processList[i]["name"].(string) < processList[j]["name"].(string)
		default:
			// Default sort by CPU usage
			return processList[i]["cpu_percent"].(float64) > processList[j]["cpu_percent"].(float64)
		}
	})

	// Limit the number of results
	if int(limit) < len(processList) {
		processList = processList[:int(limit)]
	}

	result := map[string]interface{}{
		"total_count": len(processes),
		"returned_count": len(processList),
		"sort_by":     sortBy,
		"processes":   processList,
	}

	// Convert result to JSON
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize result: %v", err)), nil
	}

	return mcp.NewToolResultText(string(data)), nil
} 