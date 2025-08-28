package main

import (
	"context"
	"fmt"
	"log"
	"sandbox/internal/executor"
	"sandbox/internal/model"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func RunCode(ctx context.Context, req *mcp.CallToolRequest, args model.CodeParams) (*mcp.CallToolResult, any, error) {

	log.Printf("[sandbox-runner] Received run_code request: language=%s", args.Language)

	log.Printf("[sandbox-runner] Executing %s code in Docker sandbox...", args.Language)
	const timeout = 10 * time.Second

	var runtime executor.Executor
	switch args.Language {
	case "python":
		runtime = executor.PythonExecutor{}
	case "java":
		runtime = executor.JavaExecutor{}
	case "cpp":
		runtime = executor.CppExecutor{}
	default:
		log.Printf("[sandbox-runner] Unsupported language: %s", args.Language)
		return nil, nil, fmt.Errorf("unsupported language: %s (python, java and cpp supported)", args.Language)
	}

	result := runtime.Execute(ctx, args.Code, timeout)

	log.Printf("[sandbox-runner] Execution finished. Output length: %d, Error: %v", len(result.Output), result.Err)

	resp := &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: result.Output},
		},
	}

	if result.Err != nil {
		log.Printf("[sandbox-runner] Execution failed: %v", result.Err)
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error: %v\nOutput:\n%s", result.Err, result.Output)},
			},
		}, nil, fmt.Errorf("execution failed: %w", result.Err)
	}

	return resp, nil, nil
}

func main() {
	log.Println("[sandbox-runner] Starting MCP server...")

	// Create server
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "sandbox-runner",
		Version: "0.1.1",
	}, nil)

	// Register tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "run_code",
		Description: "Execute code inside a Docker sandbox",
	}, RunCode)

	log.Println("[sandbox-runner] MCP server is ready and waiting for requests.")

	// Run server over stdio
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
