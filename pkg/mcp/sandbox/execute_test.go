package sandbox

import (
	"context"
	"testing"
	"time"
)

func TestExecutePythonSimpleDocker(t *testing.T) {
	ctx := context.Background()
	code := `print("Hello from test!")`
	result := ExecutePythonSimpleDocker(ctx, code, 10*time.Second)

	if result.Err != nil {
		t.Fatalf("Execution failed: %v", result.Err)
	}

	if want := "Hello from test!"; !contains(result.Output, want) {
		t.Errorf("Output does not contain expected string. Got: %q, want: %q", result.Output, want)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(s) > len(substr) && (contains(s[1:], substr) || contains(s[:len(s)-1], substr))))
}
