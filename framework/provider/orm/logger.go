package orm

import (
	"context"
	"time"

	"github.com/awaketai/goweb/framework/contract"
	gormLogger "gorm.io/gorm/logger"
)

// OrmLogger implement Gorm Logger.interface
type OrmLogger struct {
	logger contract.Log
}

func NewOrmLogger(logger contract.Log) *OrmLogger {
	return &OrmLogger{logger: logger}
}

func (l *OrmLogger) Info(ctx context.Context, s string, i ...any) {
	fields := map[string]any{
		"fields": i,
	}
	l.logger.Info(ctx, s, fields)
}

func (l *OrmLogger) Warn(ctx context.Context, s string, i ...any) {
	fields := map[string]any{
		"fields": i,
	}
	l.logger.Warn(ctx, s, fields)
}

func (l *OrmLogger) Error(ctx context.Context, s string, i ...any) {
	fields := map[string]any{
		"fields": i,
	}
	l.logger.Error(ctx, s, fields)
}

func (l *OrmLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, rows := fc()
	elapsed := time.Since(begin)

	fields := map[string]any{
		"begin": begin,
		"error": err,
		"sql":   sql,
		"rows":  rows,
		"time":  elapsed,
	}
	s := "orm trace sql"
	l.logger.Trace(ctx, s, fields)
}

func (l *OrmLogger) LogMode(logLevel gormLogger.LogLevel) gormLogger.Interface {
	return l
}
