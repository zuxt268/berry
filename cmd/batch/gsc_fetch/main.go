package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zuxt268/berry/internal/infrastructure"
	"github.com/zuxt268/berry/internal/interface/adapter"
	"github.com/zuxt268/berry/internal/repository"
	"github.com/zuxt268/berry/internal/usecase"
)

func main() {
	slog.Info("starting GSC data fetch batch")

	// タイムアウト + シグナル制御
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	ctx, stop := signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	// DB接続
	db, err := infrastructure.NewMySQLConnection()
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer func() {
		sqlDB, err := db.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
	}()

	dbDriver := infrastructure.NewDBDriver(db, db)

	// 依存組み立て
	gscConnRepo := repository.NewGSCConnectionRepository(dbDriver)
	gscDailyReportRepo := repository.NewGSCDailyReportRepository(dbDriver)
	gscDataAdapter := adapter.NewGSCDataAdapter()
	gscFetchUseCase := usecase.NewGSCFetchUseCase(gscConnRepo, gscDailyReportRepo, gscDataAdapter)

	// 対象日: デフォルトは3日前（GSCは2-3日のラグあり）、引数で指定可能
	targetDate := time.Now().AddDate(0, 0, -3).Truncate(24 * time.Hour)
	if len(os.Args) > 1 {
		parsed, err := time.Parse("2006-01-02", os.Args[1])
		if err != nil {
			slog.Error("invalid date format, use YYYY-MM-DD", "input", os.Args[1], "error", err)
			os.Exit(1)
		}
		targetDate = parsed
	}

	// 実行
	if err := gscFetchUseCase.Execute(ctx, targetDate); err != nil {
		slog.Error("batch execution failed", "error", err)
		os.Exit(1)
	}

	slog.Info("GSC data fetch batch completed")
}
