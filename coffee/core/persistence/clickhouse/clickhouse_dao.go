package clickhouse

import (
	"coffee/constants"
	"coffee/core/appcontext"
	"context"
	"gorm.io/gorm"
)

type ClickhouseDao struct {
	ctx context.Context
}

func NewClickhouseDao(ctx context.Context) *ClickhouseDao {
	return &ClickhouseDao{
		ctx: ctx,
	}
}

func (r *ClickhouseDao) SetCtx(ctx context.Context) {
	r.ctx = ctx
}

func (r *ClickhouseDao) GetSession(ctx context.Context) *gorm.DB {
	reqContext := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	session := reqContext.GetCHSession()
	reqContext.Mutex.Lock()
	if session == nil {
		// Create new txn if one already does not exist
		session = InitializeClickhouseSession(ctx)
		reqContext.SetCHSession(ctx, session)
	}
	reqContext.Mutex.Unlock()
	return session.(ClickhouseSession).db
}
