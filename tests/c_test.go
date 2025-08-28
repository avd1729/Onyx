package tests

import (
	"context"
	"sandbox/internal/executor"
	"sandbox/internal/utils"
	"testing"
	"time"
)

func TestExecuteC(t *testing.T) {
	ctx := context.Background()
	code := `#include <stdio.h>
	int main() {
		printf("Hello!\n");
		return 0;
	}`

	runtime := executor.CppExecutor{}
	result := runtime.Execute(ctx, code, 10*time.Second)

	if result.Err != nil {
		t.Fatalf("Execution failed: %v", result.Err)
	}

	if want := "Hello!"; !utils.Contains(result.Output, want) {
		t.Errorf("Output does not contain expected string. Got: %q, want: %q", result.Output, want)
	}
}
