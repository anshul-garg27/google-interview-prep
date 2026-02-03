package clickhouse

import (
	"coffee/core/persistence"
	"context"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func InitializeClickhouseSession(ctx context.Context) persistence.Session {
	session := ClickhouseSession{}
	session.db = session.Init()
	return session
}

type ClickhouseSession struct {
	db *gorm.DB
}

func (s ClickhouseSession) Init() *gorm.DB {
	log.Debug("Starting transaction in DB")
	db := ClickhouseDB()
	if db != nil {
		return db.Session(&gorm.Session{SkipDefaultTransaction: false}).Begin()
	}
	return nil
}

func (s ClickhouseSession) PerformAfterCommit(ctx context.Context, fn func(ctx context.Context)) {
}

func (s ClickhouseSession) Close(ctx context.Context) {
	if re := recover(); re != nil {
		s.Rollback(ctx)
		panic(re)
	} else {
		s.Commit(ctx)
	}
}

func (s ClickhouseSession) Rollback(ctx context.Context) {
	log.Debug("Rolling Back changes in DB")
	s.db.Rollback()
}

func (s ClickhouseSession) Commit(ctx context.Context) {
	log.Debug("Committing changes in DB")
	s.db.Commit()
}

func (s ClickhouseSession) performAfterCommitCallback(ctx context.Context) {
}

func (s ClickhouseSession) afterCreateCommitCallback(db *gorm.DB) {
	log.Debug("afterCreateCommitCallback")
}

func (s ClickhouseSession) afterUpdateCommitCallback(db *gorm.DB) {
	log.Debug("afterUpdateCommitCallback")
}

func (s ClickhouseSession) afterDeleteCommitCallback(db *gorm.DB) {
	log.Debug("afterDeleteCommitCallback")
}
