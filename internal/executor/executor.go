package executor

import (
	"context"
	"sandbox/internal/model"
	"time"
)

type Executor interface {
	Execute(ctx context.Context, code string, timeout time.Duration, dependencies ...string) model.ExecResult
}
