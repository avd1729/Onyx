// main.go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"sandbox/pkg/mcp/sandbox"
	"sandbox/pkg/model"
)

func RunCode(ctx context.Context, req *mcp.CallToolRequest, args model.CodeParams) (*mcp.CallToolResult, any, error) {
	log.Printf("[sandbox-runner] Received run_code request: language=%s", args.Language)
	if args.Language != "python" {
		log.Printf("[sandbox-runner] Unsupported language: %s", args.Language)
		return nil, nil, fmt.Errorf("unsupported language: %s (only python supported now)", args.Language)
	}

	log.Printf("[sandbox-runner] Executing Python code in Docker sandbox...")
	const timeout = 10 * time.Second
	result := sandbox.ExecutePythonSimpleDocker(ctx, args.Code, timeout)

	log.Printf("[sandbox-runner] Execution finished. Output length: %d, Error: %v", len(result.Output), result.Err)

	resp := &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: result.Output},
		},
	}

	if result.Err != nil {
		log.Printf("[sandbox-runner] Execution failed: %v", result.Err)
		return resp, nil, fmt.Errorf("execution failed: %w", result.Err)
	}

	return resp, nil, nil
}

func main() {
	log.Println("[sandbox-runner] Starting MCP server...")

	// Create server
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "sandbox-runner",
		Version: "0.1.0",
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
