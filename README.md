# Onyx

This project provides a minimal MCP (Model Context Protocol) server for executing Python code securely in a Docker sandbox. It is designed for easy integration with Claude Desktop and other MCP clients.

## Features
- **Python runtime support**: Execute arbitrary Python code in a sandboxed Docker container.
- **MCP protocol**: Exposes a `run_code` tool for remote code execution.
- **Security**: Uses Docker for process isolation.
- **Logging**: Prints execution and error logs to stderr.
- **CI**: Automated tests run on every push/PR via GitHub Actions.

## Project Structure
```
sandbox/
├── .github/
│   └── workflows/ci.yml      # GitHub Actions CI
├── cmd/
│   └── server/
│       └── main.go           # MCP server entrypoint
├── internal/
│   ├── model/                # Code params, executor interface, result types
|   ├── executor/
│   └── utils/                # Docker availability check
├── go.mod, go.sum
```

## Setup Instructions

### Prerequisites
- Go 1.20+
- Docker (Desktop or Engine) installed and running
- (Optional) Claude Desktop for interactive use

### Build and Run the Server
```sh
cd /path/to/sandbox
# Run directly
GO111MODULE=on go run ./cmd/server/main.go
# Or build a binary
GO111MODULE=on go build -o sandbox_server ./cmd/server/main.go
./sandbox_server
```

## Connecting with Claude Desktop
1. In Claude Desktop, add a new MCP server configuration:
   - **Command**: `go`
   - **Args**: `["run", "/absolute/path/to/cmd/server/main.go"]`
   - (Or use the built binary as the command)
2. Start the server from Claude Desktop.
3. Use the `run_code` tool with a request like:
   ```json
   {
     "language": "python",
     "code": "print('Hello from Claude!')"
   }
   ```
4. The response will contain the output of your Python code.

## Current Capabilities
- **Supported language**: Python (via Docker `python:3.11` image)
- **Isolation**: Each request runs in a fresh container
- **Error handling**: Returns errors for unsupported languages, Docker issues, or code errors
- **Logging**: All logs go to stderr (do not break MCP protocol)

## Limitations / TODO
- No support for languages other than Python
- No persistent state or file access
- No advanced resource limits (CPU/memory) beyond Docker defaults
- Dependency installation is basic (via pip if specified)
- No authentication or user isolation

---

Feel free to open issues or PRs for improvements!
