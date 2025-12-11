package main

import (
	"context"

	"github.com/jmoiron/sqlx"

	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/config"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/services/server"
	"code.mipt.ru/fullstack2025a/serdechnyjgl-project/internal/utils/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// main - точка входа в приложение.
func main() {
	// Инициализация глобального логгера.
	logger.InitLogger("")
	// Создание корневого контекста.
	ctx := context.Background()
	// Загрузка конфигурации из файла.
	cfg, err := config.LoadConfig(ctx, "config/config.yaml")
	if err != nil {
		logger.Fatalf(ctx, "failed to load config: %v", err)
	}

	// Получение URL базы данных из конфигурации.
	dbURL := cfg.Database.GetDBUrl()
	logger.Infof(ctx, "Database URL: %v", dbURL)
	if dbURL == "" {
		logger.Fatalf(ctx, "Database URL is empty")
	}

	// Применение миграций к базе данных.
	if err := runMigrations(ctx, dbURL); err != nil {
		logger.Fatalf(ctx, "Migrations error: %v", err)
	}

	// Инициализация подключения к базе данных.
	db, err := sqlx.Open("postgres", dbURL)
	// Создание экземпляра сервера.
	// Передаются:
	// - Новый экземпляр Fiber
	// - Базовый *sql.DB из sqlx (для совместимости)
	// - Загруженная конфигурация
	// - Контекст
	srvr, err := server.NewServer(fiber.New(), db.DB, cfg, ctx)
	if err != nil {
		logger.Fatalf(context.Background(), "Server init error: %v", err)
	}

	// Запуск HTTP-сервера на порту, указанном в конфигурации.
	// Метод Start блокирует выполнение, пока сервер не будет остановлен.
	logger.Fatalf(context.Background(), "%v", srvr.Start(":"+cfg.Server.Port))
}

// runMigrations применяет миграции базы данных из директории 'migrations'.
// Она ищет SQL-файлы в папке 'migrations' и применяет их к БД по указанному URL.
func runMigrations(ctx context.Context, dbURL string) error {
	logger.Infof(ctx, "Migration start...")
	// Создание экземпляра мигратора.
	// 'file://migrations' указывает на локальную папку с миграциями.
	// dbURL - строка подключения к целевой БД (например, postgresql://user:pass@host:port/dbname)
	m, err := migrate.New(
		"file://migrations",
		dbURL,
	)
	if err != nil {
		return err
	}

	// Применение всех накативных (up) миграций, которые еще не были применены.
	err = m.Up()
	if err != nil {
		// ErrNoChange - это не ошибка, а ожидаемая ситуация, когда все миграции уже применены.
		if err == migrate.ErrNoChange {
			logger.Infof(ctx, "No new migrations to apply")
			return nil
		}
		return err
	}

	logger.Infof(ctx, "Migration completed successfully")
	return nil
}
