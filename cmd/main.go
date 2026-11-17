package main

import (
	"awesomeProject1/internal/calculationService"
	"awesomeProject1/internal/db"
	"awesomeProject1/internal/handlers"
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func main() {
	// 1. Инициализация БД
	database, err := db.InitDB()
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}

	// 2. Настройка Echo
	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.RequestLogger())

	// 3. Сборка зависимостей (Dependency Injection)
	calcRepo := calculationService.NewCalculationRepository(database)
	calcService := calculationService.NewCalculationService(calcRepo)
	calcHandlers := handlers.NewCalculationHandler(calcService)

	// 4. Маршруты
	e.GET("/calculations", calcHandlers.GetCalculations)
	e.POST("/calculations", calcHandlers.PostCalculation)
	e.PATCH("/calculations/:id", calcHandlers.PatchCalculation)
	e.DELETE("/calculations/:id", calcHandlers.DeleteCalculation)

	// 5. Запуск сервера в горутине (чтобы не блокировать основной поток)
	go func() {
		if err := e.Start("localhost:8080"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// 6. Ожидание сигнала для Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	// 7. Контекст с таймаутом для остановки
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	// 8. Закрытие БД
	sqlDB, err := database.DB()
	if err == nil {
		sqlDB.Close()
		log.Info("Database connection closed.")
	}

	log.Info("Server exited properly")
}
