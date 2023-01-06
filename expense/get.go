package expense

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func (h *handler) GetExpensesByIDHandler(c echo.Context) error {
	expense := Expense{}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error{Message: "ID is missing"})
	}

	stmt, err := h.DB.Prepare("SELECT id, title, amount, note, tags FROM expenses WHERE id = $1")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Error{Message: err.Error()})
	}

	rows := stmt.QueryRow(id)
	if rows.Err() != nil {
		return c.JSON(http.StatusInternalServerError, Error{Message: err.Error()})
	}

	err = rows.Scan(&expense.ID, &expense.Title, &expense.Amount, &expense.Note, pq.Array(&expense.Tags))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Error{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, expense)
}

func (h *handler) GetExpensesHandler(c echo.Context) error {
	expenses := []Expense{}
	stmt, err := h.DB.Prepare("SELECT id, title, amount, note, tags FROM expenses ORDER BY id ASC")
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error{Message: err.Error()})
	}

	rows, err := stmt.Query()
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error{Message: err.Error()})
	}

	for rows.Next() {
		expense := Expense{}
		err := rows.Scan(&expense.ID, &expense.Title, &expense.Amount, &expense.Note, pq.Array(&expense.Tags))
		if err != nil {
			return c.JSON(http.StatusBadRequest, Error{Message: err.Error()})
		}
		expenses = append(expenses, expense)
	}

	return c.JSON(http.StatusOK, expenses)
}
