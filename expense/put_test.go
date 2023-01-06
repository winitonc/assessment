//go:build unit

package expense

import (
	"net/http"
	"net/http/httptest"
	"strconv"
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
		"title": "apple smoothie",
		"amount": 89.0,
		"note": "note", 
		"tags": ["tag1", "tag2"]
	}`
	expectBodyJSON = `{
		"id": 1,
		"title": "apple smoothie",
		"amount": 89.0,
		"note": "note", 
		"tags": ["tag1", "tag2"]
	}`
)

func TestUpdateExpenseHandler(t *testing.T) {

	ja := jsonassert.New(t)

	t.Run("Update expense should success", func(t *testing.T) {

		// Setup
		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		id := 1
		e := echo.New()
		req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(bodyJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(id))
		h := &handler{db}

		result := sqlmock.NewResult(1, 1)

		// Assertions
		mock.ExpectPrepare("UPDATE expenses SET title=$2, amount=$3, note=$4, tags=$5 WHERE id = $1").
			ExpectExec().
			WithArgs(id, "apple smoothie", 89.0, "note", pq.Array([]string{"tag1", "tag2"})).
			WillReturnResult(result)

		if assert.NoError(t, h.UpdateExpenseHandler(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			ja.Assertf(expectBodyJSON, rec.Body.String())
		}

	})

}
