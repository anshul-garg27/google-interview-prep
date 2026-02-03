package postgres

import (
	"coffee/constants"
	"coffee/core/appcontext"
	"coffee/core/persistence"
	"context"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func InitializePgSession(ctx context.Context) persistence.Session {
	session := PgSession{}
	session.db = session.Init()
	// session.db.Callback().Create().
	// 	After("gorm:commit_or_rollback_transaction").
	// 	Register("after_create_commit", session.afterCreateCommitCallback)

	// session.db.Callback().Update().
	// 	After("gorm:commit_or_rollback_transaction").
	// 	Register("after_update_commit", session.afterUpdateCommitCallback)

	// session.db.Callback().Delete().
	// 	After("gorm:commit_or_rollback_transaction").
	// 	Register("after_delete_commit", session.afterDeleteCommitCallback)
	return session
}

type PgSession struct {
	db *gorm.DB
}

func (s PgSession) PerformAfterCommit(ctx context.Context, fn func(ctx context.Context)) {
	appContext := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	appContext.SetValue("afterCommitCallback", fn)
}

func (s PgSession) Init() *gorm.DB {
	log.Debug("Starting transaction in DB")
	return DB().Session(&gorm.Session{SkipDefaultTransaction: false}).Begin()
}

func (s PgSession) Close(ctx context.Context) {
	if re := recover(); re != nil {
		s.Rollback(ctx)
		panic(re)
	} else {
		s.Commit(ctx)
	}
}

func (s PgSession) Rollback(ctx context.Context) {
	log.Debug("Rolling Back changes in DB")
	s.db = s.db.Rollback()
	if s.db.Error != nil {
		log.Error(s.db.Error)
	}
}

func (s PgSession) Commit(ctx context.Context) {
	log.Debug("Committing changes in DB")
	s.db = s.db.Commit()
	if s.db.Error != nil {
		// err := errors.Wrap(s.db.Error, "database commit failed")
		err := errors.Wrap(errors.New(s.db.Error.Error()), "database commit failed")
		log.Error(err)
	} else {
		s.performAfterCommitCallback(ctx)
	}
}

func (s PgSession) performAfterCommitCallback(ctx context.Context) {
	appContext := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	afterCommitFn := appContext.GetValue("afterCommitCallback")
	if afterCommitFn != nil {
		_fn := afterCommitFn.(func(ctx context.Context))
		_fn(ctx)
	}
}

func (s PgSession) afterCreateCommitCallback(db *gorm.DB) {
	log.Debug("afterCreateCommitCallback")
}

func (s PgSession) afterUpdateCommitCallback(db *gorm.DB) {
	log.Debug("afterUpdateCommitCallback")
}

func (s PgSession) afterDeleteCommitCallback(db *gorm.DB) {
	log.Debug("afterDeleteCommitCallback")
}
