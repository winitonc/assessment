package expense

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kinbiko/jsonassert"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

var (
	expectBodyJSON = `{
		"id": 1,
		"title": "Test",
		"amount": 10.99,
		"note": "note",
		"tags": ["tag1", "tag2"]
	}`
)

func TestGetExpensesByIDHandler(t *testing.T) {

	ja := jsonassert.New(t)

	t.Run("Get expenses by ID should success and return record", func(t *testing.T) {

		// Setup
		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		ID := 1

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/expenses/"+strconv.Itoa(ID), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(ID))
		h := &handler{db}

		// Assertions
		mockRows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
			AddRow(1, "Test", 10.99, "note", pq.Array([]string{"tag1", "tag2"}))

		mock.ExpectPrepare("SELECT id, title, amount, note, tags FROM expenses WHERE id = $1").
			ExpectQuery().
			WithArgs(ID).
			WillReturnRows(mockRows)

		if assert.NoError(t, h.GetExpensesByIDHandler(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			ja.Assertf(expectBodyJSON, rec.Body.String())
		}

	})

	// t.Run("Get expenses by ID should return fail", func(t *testing.T) {

	// 	// Setup
	// 	db, mock, _ := sqlmock.New()
	// 	ID := 1

	// 	e := echo.New()
	// 	req := httptest.NewRequest(http.MethodGet, "/expenses/"+strconv.Itoa(ID), nil)
	// 	rec := httptest.NewRecorder()
	// 	c := e.NewContext(req, rec)
	// 	c.SetParamNames("id")
	// 	c.SetParamValues(strconv.Itoa(ID))
	// 	h := &handler{db}

	// 	// Assertions
	// 	mock.ExpectPrepare("SELECT id, title, amount, note, tags FROM expenses WHERE id = $1").
	// 		ExpectQuery().
	// 		WithArgs(ID).
	// 		WillReturnError(sql.ErrNoRows)

	// 	h.GetExpensesByIDHandler(c)
	// 	fmt.Println("Response: ", rec.Body.String())
	// 	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	// 	// ja.Assertf(`{"message": "pq: Mock Error from DB"}`, rec.Body.String())

	// })

}
