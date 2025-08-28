package executor

import (
	"context"
	"sandbox/internal/utils"
	"testing"
	"time"
)

func TestExecuteJava(t *testing.T) {
	ctx := context.Background()
	code := `class Main {
		public static void main(String[] args) {
			System.out.println("Hello!");
		}
	}`

	runtime := JavaExecutor{}
	result := runtime.Execute(ctx, code, 10*time.Second)

	if result.Err != nil {
		t.Fatalf("Execution failed: %v", result.Err)
	}

	if want := "Hello!"; !utils.Contains(result.Output, want) {
		t.Errorf("Output does not contain expected string. Got: %q, want: %q", result.Output, want)
	}
}
