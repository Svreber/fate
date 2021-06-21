package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestRoll(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	for i := 0; i < 10000; i++ {
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		assert.NoError(t, roll(ctx))

		payload, err := io.ReadAll(rec.Result().Body)
		assert.NoError(t, err)

		res, err := strconv.Atoi(string(payload))
		assert.NoError(t, err)

		assert.GreaterOrEqual(t, res, 1)
		assert.LessOrEqual(t, res, 6)
	}
}
