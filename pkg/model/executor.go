package model

import (
	"context"
	"time"
)

type Executor interface {
	Execute(ctx context.Context, code string, timeout time.Duration, dependencies ...string) ExecResult
}
