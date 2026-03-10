package testhelper

import (
	"context"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/zuxt268/berry/internal/infrastructure"
	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// TestDB テスト用のデータベース接続情報を保持
type TestDB struct {
	Container *mysql.MySQLContainer
	DB        *gorm.DB
	DBDriver  infrastructure.DBDriver
	DSN       string
	ctx       context.Context
}

// SetupTestDB testcontainersを使ってMySQLを起動し、マイグレーションを実行
func SetupTestDB(t *testing.T) *TestDB {
	t.Helper()

	ctx := context.Background()

	// MySQLコンテナを起動
	mysqlContainer, err := mysql.Run(ctx,
		"mysql:8.0",
		mysql.WithDatabase("testdb"),
		mysql.WithUsername("testuser"),
		mysql.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("port: 3306").
				WithOccurrence(1),
		),
	)
	if err != nil {
		t.Fatalf("failed to start mysql container: %s", err)
	}

	// 接続文字列を取得
	dsn, err := mysqlContainer.ConnectionString(ctx, "charset=utf8mb4&parseTime=true")
	if err != nil {
		t.Fatalf("failed to get connection string: %s", err)
	}

	// GORMでDB接続
	db, err := gorm.Open(gormmysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to connect to database: %s", err)
	}

	// TODO: AutoMigrateでスキーマを作成
	// if err := db.AutoMigrate(&model.User{}); err != nil {
	//     t.Fatalf("failed to run auto migrate: %s", err)
	// }

	// DBDriverを作成
	dbDriver := infrastructure.NewDBDriver(db, db)

	return &TestDB{
		Container: mysqlContainer,
		DB:        db,
		DBDriver:  dbDriver,
		DSN:       dsn,
		ctx:       ctx,
	}
}

// Teardown テスト用DBをクリーンアップ
func (td *TestDB) Teardown(t *testing.T) {
	t.Helper()

	sqlDB, err := td.DB.DB()
	if err == nil {
		_ = sqlDB.Close()
	}

	if err := td.Container.Terminate(td.ctx); err != nil {
		t.Logf("failed to terminate container: %s", err)
	}
}
