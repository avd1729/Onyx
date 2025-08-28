package executor

import (
	"bytes"
	"context"
	"errors"
	"log"
	"os/exec"
	"sandbox/internal/model"
	"sandbox/internal/utils"
	"time"
)

type JavaScriptExecutor struct{}

func (js JavaScriptExecutor) Execute(
	ctx context.Context,
	code string,
	timeout time.Duration,
	dependencies ...string, // For future: support adding external deps
) model.ExecResult {

	logPrefix := "[sandbox-exec]"
	log.Println(logPrefix, "Preparing to execute Js code in Docker...")

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if !utils.IsDockerAvailable() {
		log.Println(logPrefix, "Docker is not available or not running")
		return model.ExecResult{Err: errors.New("docker is not available or not running")}
	}

	if len(dependencies) > 0 {
		log.Println(logPrefix, "Warning: Js dependencies are not yet supported. Ignoring...")
	}

	dockerCmd := []string{
		"run", "--rm", "-i",
		"--network", "none",
		"--cpus", "0.5",
		"--memory", "256m", "--memory-swap", "256m",
		"--pids-limit", "64",
		"--read-only",
		"--tmpfs", "/workspace:rw,exec,uid=1000,gid=1000",
		"-w", "/workspace",
		"--cap-drop=ALL",
		"--security-opt", "no-new-privileges",
		"--user", "1000:1000",
		"node:20",
		"sh", "-c", "cat > main.js && node main.js",
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
		return model.ExecResult{Output: output, Err: errors.New(output)}
	}

	log.Println(logPrefix, "Execution succeeded.")
	return model.ExecResult{Output: output}
}
