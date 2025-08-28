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

func Contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(s) > len(substr) && (Contains(s[1:], substr) || Contains(s[:len(s)-1], substr))))
}
