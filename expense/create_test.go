//go:build unit

package expense

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kinbiko/jsonassert"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

var (
	bodyJSON = `{
		"title": "Test",
		"amount": 10.99,
		"note": "note", 
		"tags": ["tag1", "tag2"]
	}`
	expectBodyJSON = `{
		"id": 1,
		"title": "Test",
		"amount": 10.99,
		"note": "note", 
		"tags": ["tag1", "tag2"]
	}`
)

func TestCreateExpenseHandler(t *testing.T) {

	ja := jsonassert.New(t)

	t.Run("Create expense should success", func(t *testing.T) {

		// Setup
		db, mock, _ := sqlmock.New()

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(bodyJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		h := &handler{db}

		// Assertions
		mock.ExpectQuery("INSERT INTO expenses").
			WithArgs("Test", 10.99, "note", pq.Array([]string{"tag1", "tag2"})).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).FromCSVString("1"))

		if assert.NoError(t, h.CreateExpenseHandler(c)) {
			assert.Equal(t, http.StatusCreated, rec.Code)
			ja.Assertf(expectBodyJSON, rec.Body.String())
		}

	})

	t.Run("Create expense should return fail when get error from DB ", func(t *testing.T) {

		// Setup
		db, mock, _ := sqlmock.New()

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(bodyJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		h := &handler{db}

		// Assertions
		mock.ExpectQuery("INSERT INTO expenses").
			WithArgs("Test", 10.99, "note", pq.Array([]string{"tag1", "tag2"})).WillReturnError(&pq.Error{Message: "Mock Error from DB"})

		h.CreateExpenseHandler(c)
		fmt.Println("Response: ", rec.Body.String())
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		ja.Assertf(`{"message": "pq: Mock Error from DB"}`, rec.Body.String())

	})

}
