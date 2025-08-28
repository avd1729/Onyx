# Onyx

Onyx is a **MCP (Model Context Protocol) server** that executes code securely inside **Docker sandboxes**. It supports multiple programming languages and integrates seamlessly with **Claude Desktop** or other MCP clients. Onyx makes it safe and easy to execute arbitrary code in a **sandboxed MCP environment**. Perfect for **Claude Desktop** power users or anyone building AI workflows with executable code.

## Features

* 🔹 **Multi-language support**:

  * Python (`python:3.11`)
  * Java (`openjdk:17`)
  * C (`gcc:12`)
  * C++ (`gcc:12`)
  * JavaScript / Node.js (`node:20`)
  * Rust (`rust:1.72`)

* 🔹 **Docker sandboxing**:

  * Network disabled (`--network none`)
  * Read-only FS with tmpfs mounts for safe writes
  * Limited CPU, memory, and process count
  * Non-root execution (`--user 1000:1000`)

* 🔹 **MCP protocol**: Exposes a `run_code` tool for executing arbitrary code.

* 🔹 **Logging**: All execution logs go to `stderr` (safe for MCP clients).

* 🔹 **CI-ready**: Automated tests for executors via GitHub Actions.

---

## Project Structure

```
sandbox/
├── .github/
│   └── workflows/ci.yml       # GitHub Actions CI
├── cmd/
│   └── server/
│       └── main.go            # MCP server entrypoint
├── internal/
│   ├── model/                 # Code params, executor interface, result types
│   ├── executor/              # Language-specific executors
│   └── utils/                 # Docker availability checks
├── tests/
├── go.mod, go.sum
```

---

## Setup Instructions

### 1. Prerequisites

* [Go 1.20+](https://go.dev/dl/)
* [Docker Desktop](https://docs.docker.com/desktop/) (Windows/macOS) or Docker Engine (Linux)

Verify both:

```sh
go version
docker --version
```

### 2. Pull Required Docker Images

To avoid delays during the first execution, pre-pull all language runtimes:

```sh
docker pull python:3.11
docker pull openjdk:17
docker pull gcc:12
docker pull node:20
docker pull rust:1.72
```

### 3. Build and Run the Server

From the root of the repo:

```sh
# Run directly
go run ./cmd/server/main.go

# Or build a binary (recommended for Claude Desktop)
go build -o sandbox_server ./cmd/server/main.go
```

You should now have `sandbox_server` (or `sandbox_server.exe` on Windows).

---

## Connecting with Claude Desktop

1. Locate Claude Desktop’s config file:

   ```powershell
   code $env:AppData\Claude\claude_desktop_config.json
   ```

2. Add an MCP server entry for **Onyx**:

   ```json
   {
     "mcpServers": {
       "onyx": {
         "command": "<absolute path>/sandbox_server.exe",
         "args": []
       }
     }
   }
   ```

3. Restart Claude Desktop.

4. You should now see the `run_code` tool available.

---

## Usage

Call the `run_code` tool with JSON like:

```json
{
  "language": "python",
  "code": "print('Hello from Python!')"
}
```

Example for **Java**:

```json
{
  "language": "java",
  "code": "class Main { public static void main(String[] args) { System.out.println(\"Hello Java!\"); } }"
}
```

Example for **Rust**:

```json
{
  "language": "rust",
  "code": "fn main() { println!(\"Hello from Rust!\"); }"
}
```

---

## Extending to New Languages

Adding support for another language requires:

1. **Choose a Docker image** with the language’s compiler/runtime.
   Example: `golang:1.22` for Go, `php:8.2-cli` for PHP.

2. **Implement an Executor** in `internal/executor/<lang>.go`:

   * Define a struct (e.g., `GoExecutor`).
   * Implement the `Execute` method to:

     * Pipe source code into the container.
     * Compile/interpret it inside `/workspace`.
     * Run the output binary/script.

3. **Register it in `main.go`**:

   ```go
   if args.Language == "go" {
       runtime = executor.GoExecutor{}
   }
   ```

4. **Write a test** in `internal/executor/go_test.go`.

5. **Pull the Docker image** ahead of time (`docker pull golang:1.22`).

