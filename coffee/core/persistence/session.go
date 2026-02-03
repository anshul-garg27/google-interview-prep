package persistence

import (
	"context"
)

type Session interface {
	Close(ctx context.Context)
	Rollback(ctx context.Context)
	Commit(ctx context.Context)
	PerformAfterCommit(ctx context.Context, fn func(ctx context.Context))
}
