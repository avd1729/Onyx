package sandbox

import (
	"bytes"
	"context"
	"errors"
	"log"
	"os/exec"
	"sandbox/pkg/model"
	"time"
)

// ExecutePythonSimpleDocker runs Python code in a Docker container using python:3.11 and feeds code via stdin.
func ExecutePythonSimpleDocker(ctx context.Context, code string, timeout time.Duration) model.ExecResult {
	logPrefix := "[sandbox-exec]"
	log.Println(logPrefix, "Preparing to execute Python code in Docker...")

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if !isDockerAvailable() {
		log.Println(logPrefix, "Docker is not available or not running")
		return model.ExecResult{Err: errors.New("docker is not available or not running")}
	}

	log.Println(logPrefix, "Running docker run --rm -i python:3.11 python -")
	cmd := exec.CommandContext(ctx, "docker", "run", "--rm", "-i", "python:3.11", "python", "-")
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

func isDockerAvailable() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", "version")
	return cmd.Run() == nil
}
