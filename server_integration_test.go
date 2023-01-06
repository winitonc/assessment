//go:build integration

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/winitonc/assessment/expense"
)

const serverPort = 2565

func init() {
	// Setup server
	eh := echo.New()
	go func(e *echo.Echo) {
		db, err := sql.Open("postgres", "postgres://postgres:password@db:5432/postgres?sslmode=disable")
		if err != nil {
			log.Fatal(err)
		}

		h := expense.InitHandler(db)

		e.POST("/expenses", h.CreateExpenseHandler)
		e.GET("/expenses", h.GetExpensesHandler)
		e.PUT("/expenses/:id", h.UpdateExpenseHandler)
		e.GET("/expenses/:id", h.GetExpensesByIDHandler)
		e.Start(fmt.Sprintf(":%d", serverPort))
	}(eh)
	for {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", serverPort), 30*time.Second)
		if err != nil {
			log.Println(err)
		}
		if conn != nil {
			conn.Close()
			break
		}
	}
}
func TestCreateExpenseAPI(t *testing.T) {

	// Arrange
	reqBody := `{
		"title": "strawberry smoothie",
        "amount": 13.26,
        "note": "night market promotion discount 10 bath",
        "tags": [
            "food",
            "beverage"
        ]
	}`
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:%d/expenses", serverPort), strings.NewReader(reqBody))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "November 10, 2009")

	// Act
	var resp *http.Response
	client := http.Client{}
	resp, err = client.Do(req)
	assert.NoError(t, err)

	var byteBody []byte
	byteBody, err = io.ReadAll(resp.Body)
	assert.NoError(t, err)
	resp.Body.Close()

	// Assertions
	var exp expense.Expense
	err = json.Unmarshal(byteBody, &exp)
	fmt.Println("Response >>>> ", string(byteBody))
	if assert.NoError(t, err) {
		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.Equal(t, "strawberry smoothie", exp.Title)
		assert.Equal(t, 13.26, exp.Amount)
		assert.Equal(t, "night market promotion discount 10 bath", exp.Note)
		assert.Equal(t, []string{"food", "beverage"}, exp.Tags)
	}

}

func TestGetExpenseByID(t *testing.T) {
	// Arrange
	reqBody := ``
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:%d/expenses/2", serverPort), strings.NewReader(reqBody))
	assert.NoError(t, err)

	// Act
	var resp *http.Response
	client := http.Client{}
	resp, err = client.Do(req)
	assert.NoError(t, err)

	var byteBody []byte
	byteBody, err = io.ReadAll(resp.Body)
	assert.NoError(t, err)
	resp.Body.Close()

	expect := `{"id":2,"title":"Title IT","amount":10.99,"note":"Note IT","tags":["tagsIT1","tagsIT2"]}` + "\n"
	if assert.NoError(t, err) {
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, expect, string(byteBody))
	}
}

func TestUpdateExpenseAPI(t *testing.T) {
	// Arrange
	reqBody := `{
		"title": "strawberry",
        "amount": 14.65,
        "note": "Test update",
        "tags": [ "tags1" ]
	}`
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("http://localhost:%d/expenses/1", serverPort), strings.NewReader(reqBody))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Act
	var resp *http.Response
	client := http.Client{}
	resp, err = client.Do(req)
	assert.NoError(t, err)

	var byteBody []byte
	byteBody, err = io.ReadAll(resp.Body)
	assert.NoError(t, err)
	resp.Body.Close()

	var exp expense.Expense
	err = json.Unmarshal(byteBody, &exp)
	if assert.NoError(t, err) {
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, 1, exp.ID)
		assert.Equal(t, "strawberry", exp.Title)
		assert.Equal(t, 14.65, exp.Amount)
		assert.Equal(t, "Test update", exp.Note)
		assert.Equal(t, []string{"tags1"}, exp.Tags)
	}
}

func TestGetAllExpensesAPI(t *testing.T) {
	// Arrange
	reqBody := ``
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:%d/expenses", serverPort), strings.NewReader(reqBody))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Act
	var resp *http.Response
	client := http.Client{}
	resp, err = client.Do(req)
	assert.NoError(t, err)

	var byteBody []byte
	byteBody, err = io.ReadAll(resp.Body)
	assert.NoError(t, err)
	resp.Body.Close()

	// Assertions
	expect := `[{"id":1,"title":"strawberry","amount":14.65,"note":"Test update","tags":["tags1"]},{"id":2,"title":"Title IT","amount":10.99,"note":"Note IT","tags":["tagsIT1","tagsIT2"]}]` + "\n"
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Greater(t, len(byteBody), 0)
		assert.Equal(t, expect, string(byteBody))
	}
}
