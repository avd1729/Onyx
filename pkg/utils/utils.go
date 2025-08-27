package utils

import (
	"context"
	"os/exec"
	"time"
)

func IsDockerAvailable() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", "version")
	return cmd.Run() == nil
}
