package expense

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func (h *handler) UpdateExpenseHandler(c echo.Context) error {
	var expense Expense
	err := c.Bind(&expense)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error{Message: err.Error()})
	}

	expense.ID, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error{Message: "ID is missing"})
	}

	stmt, err := h.DB.Prepare("UPDATE expenses SET title=$2, amount=$3, note=$4, tags=$5 WHERE id = $1")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Error{Message: err.Error()})
	}

	if _, err := stmt.Exec(expense.ID, expense.Title, expense.Amount, expense.Note, pq.Array(expense.Tags)); err != nil {
		return c.JSON(http.StatusInternalServerError, Error{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, expense)
}
