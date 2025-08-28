package executor

import (
	"bytes"
	"context"
	"errors"
	"log"
	"os/exec"
	"sandbox/internal/model"
	"sandbox/internal/utils"
	"strings"
	"time"
)

type PythonExecutor struct{}

func (p PythonExecutor) Execute(
	ctx context.Context,
	code string,
	timeout time.Duration,
	dependencies ...string,
) model.ExecResult {

	logPrefix := "[sandbox-exec]"
	log.Println(logPrefix, "Preparing to execute Python code in Docker...")

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if !utils.IsDockerAvailable() {
		log.Println(logPrefix, "Docker is not available or not running")
		return model.ExecResult{Err: errors.New("docker is not available or not running")}
	}

	// If dependencies exist, prepare install command
	installCmd := ""
	if len(dependencies) > 0 {
		installCmd = "pip install " + strings.Join(dependencies, " ") + " && "
	}

	// Extra Docker sandboxing measures to consider:
	//
	// Security / Isolation:
	//   --network none              → disable network access
	//   --read-only                 → run with a read-only filesystem
	//   --tmpfs /tmp                → provide a safe writable temp directory
	//   --cap-drop=ALL              → drop all Linux capabilities
	//   --security-opt no-new-privileges → prevent privilege escalation
	//   --security-opt seccomp=...  → restrict syscalls via seccomp profile
	//   --user 1000:1000            → avoid running as root
	//
	// Resource Limits:
	//   --cpus="0.5"                → restrict to half a CPU core
	//   --memory="256m"             → cap memory usage
	//   --memory-swap="256m"        → prevent swap abuse
	//   --pids-limit=64             → limit process spawning (avoid fork bombs)
	//   --stop-timeout=5            → enforce container stop timeout
	//
	// Execution Environment Control:
	//   - Whitelist allowed pip packages only
	//   - Pre-build images with safe dependencies instead of installing on the fly
	//   - Never mount host filesystems; only ephemeral volumes if needed
	//   - Log execution metadata (time, resource usage) for auditing
	//
	// Example safer docker run could look like:
	//   docker run --rm -i \
	//     --network none \
	//     --cpus=0.5 \
	//     --memory=256m --memory-swap=256m \
	//     --pids-limit=64 \
	//     --read-only --tmpfs /tmp \
	//     --cap-drop=ALL \
	//     --security-opt no-new-privileges \
	//     --user 1000:1000 \
	//     python:3.11 sh -c "<installCmd> python -"

	dockerCmd := []string{
		"run", "--rm", "-i",
		"--network", "none",
		"--cpus", "0.5",
		"--memory", "256m", "--memory-swap", "256m",
		"--pids-limit", "64",
		"--read-only", "--tmpfs", "/tmp",
		"--cap-drop=ALL",
		"--security-opt", "no-new-privileges",
		"--user", "1000:1000",
		"python:3.11",
		"sh", "-c", installCmd + "python -",
	}

	log.Println(logPrefix, "Running:", "docker", dockerCmd)

	cmd := exec.CommandContext(ctx, "docker", dockerCmd...)
	cmd.Stdin = bytes.NewReader([]byte(code))

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		log.Println(logPrefix, "Execution timed out")
		return model.ExecResult{Err: errors.New("execution timed out")}
	}

	output := stdout.String()
	if stderr.Len() > 0 {
		if output != "" {
			output += "\n" + stderr.String()
		} else {
			output = stderr.String()
		}
	}

	if err != nil {
		log.Printf("%s Docker execution failed: %v\n", logPrefix, err)
		return model.ExecResult{Output: output, Err: errors.New("docker execution failed")}
	}

	log.Println(logPrefix, "Execution succeeded.")
	return model.ExecResult{Output: output}
}
