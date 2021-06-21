package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestRoll(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/roll", nil)

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

func TestRollWithArbitraryFaces(t *testing.T) {
	e := echo.New()

	n := rand.Int()
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/roll?faces=%d", n), nil)

	for i := 0; i < 10000; i++ {
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		assert.NoError(t, roll(ctx))

		payload, err := io.ReadAll(rec.Result().Body)
		assert.NoError(t, err)

		res, err := strconv.Atoi(string(payload))
		assert.NoError(t, err)

		assert.GreaterOrEqual(t, res, 1)
		assert.LessOrEqual(t, res, n)
	}
}

func TestRollDice(t *testing.T) {
	n := rand.Int()

	for i := 0; i < 10000; i++ {
		res := rollDice(n)

		assert.GreaterOrEqual(t, res, 1)
		assert.LessOrEqual(t, res, n)
	}
}

func TestFlip(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/flip", nil)

	for i := 0; i < 10000; i++ {
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		assert.NoError(t, flip(ctx))

		payload, err := io.ReadAll(rec.Result().Body)
		assert.NoError(t, err)

		assert.Contains(t, []string{"Head", "Tails"}, string(payload))
	}
}

func TestMagicBall(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/magic-ball", nil)

	for i := 0; i < 10000; i++ {
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		assert.NoError(t, magicBall(ctx))

		payload, err := io.ReadAll(rec.Result().Body)
		assert.NoError(t, err)

		assert.Contains(t, magicBallAnswers, string(payload))
	}
}
