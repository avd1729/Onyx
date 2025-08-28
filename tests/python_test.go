package tests

import (
	"context"
	"sandbox/internal/executor"
	"sandbox/internal/utils"
	"testing"
	"time"
)

func TestExecute(t *testing.T) {
	ctx := context.Background()
	code := `print("Hello!")`
	runtime := executor.PythonExecutor{}
	result := runtime.Execute(ctx, code, 10*time.Second)

	if result.Err != nil {
		t.Fatalf("Execution failed: %v", result.Err)
	}

	if want := "Hello!"; !utils.Contains(result.Output, want) {
		t.Errorf("Output does not contain expected string. Got: %q, want: %q", result.Output, want)
	}
}
