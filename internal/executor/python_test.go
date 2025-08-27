package executor

import (
	"context"
	"testing"
	"time"
)

func TestExecute(t *testing.T) {
	ctx := context.Background()
	code := `print("Hello!")`
	runtime := PythonExecutor{}
	result := runtime.Execute(ctx, code, 10*time.Second)

	if result.Err != nil {
		t.Fatalf("Execution failed: %v", result.Err)
	}

	if want := "Hello!"; !contains(result.Output, want) {
		t.Errorf("Output does not contain expected string. Got: %q, want: %q", result.Output, want)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(s) > len(substr) && (contains(s[1:], substr) || contains(s[:len(s)-1], substr))))
}
