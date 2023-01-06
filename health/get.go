package health

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
)

type handler struct {
	DB *sql.DB
}

func InitHealthHandler(db *sql.DB) *handler {
	return &handler{db}
}

func (h *handler) GetHealthHandler(c echo.Context) error {
	err := h.DB.Ping()
	dbStatus := "UP"
	if err != nil {
		dbStatus = "DOWN"
	}

	return c.JSON(http.StatusOK, Health{Status: "UP", Database: DBHealth{Status: dbStatus}})
}
