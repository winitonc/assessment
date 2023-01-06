package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"github.com/winitonc/assessment/expense"
	"github.com/winitonc/assessment/health"
)

func main() {
	fmt.Println("Stating application... ")

	db := expense.InitDB()

	serv := echo.New()
	serv.Use(middleware.Logger())
	serv.Use(middleware.Recover())

	healthHl := health.InitHealthHandler(db)
	serv.GET("/health", healthHl.GetHealthHandler)

	go func() {
		if err := serv.Start(":" + os.Getenv("PORT")); err != nil && err != http.ErrServerClosed {
			serv.Logger.Fatal("Shutting down server...")
		}
	}()

	// Gracefully Shutdown
	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, os.Interrupt, syscall.SIGTERM)
	signal.Notify(gracefulShutdown, os.Interrupt, syscall.SIGINT)

	<-gracefulShutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := serv.Shutdown(ctx); err != nil {
		log.Fatal("Error shutting down server", err)
	} else {
		log.Fatal("Server gracefully stopped")
	}

	if err := db.Close(); err != nil {
		log.Fatal("Error closing db connection", err)
	} else {
		log.Fatal("DB connection gracefully closed")
	}

}
