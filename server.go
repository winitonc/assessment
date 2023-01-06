package main

import (
	"context"
	"database/sql"
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
	"github.com/winitonc/assessment/authen"
	"github.com/winitonc/assessment/expense"
	"github.com/winitonc/assessment/health"
)

func main() {
	fmt.Println("Stating application... ")
	db := expense.InitDB()
	server := setupRoute(db)
	start(server)
	gracefulShutdown(server, db)
}

func start(server *echo.Echo) {
	go func() {
		if err := server.Start(":" + os.Getenv("PORT")); err != nil && err != http.ErrServerClosed {
			server.Logger.Fatal("Shutting down server...")
		}
	}()
}

func setupRoute(db *sql.DB) *echo.Echo {
	serv := echo.New()
	serv.Use(middleware.Logger())
	serv.Use(middleware.Recover())
	serv.Use(authen.UserAuth())

	healthHl := health.InitHealthHandler(db)
	serv.GET("/health", healthHl.GetHealthHandler)

	expenseHl := expense.InitHandler(db)
	serv.POST("/expenses", expenseHl.CreateExpenseHandler)
	serv.PUT("/expenses/:id", expenseHl.UpdateExpenseHandler)
	serv.GET("/expenses/:id", expenseHl.GetExpensesByIDHandler)
	serv.GET("/expenses", expenseHl.GetExpensesHandler)

	return serv
}

func gracefulShutdown(serv *echo.Echo, db *sql.DB) {
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
