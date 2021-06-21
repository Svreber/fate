package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func main() {
	port := flag.Int("port", 8080, "HTTP port")
	flag.Parse()

	// -- Basic echo --

	e := echo.New()

	// -- Add endpoints

	e.GET("/manage/health", health)

	rand.Seed(time.Now().UnixNano())
	e.GET("/roll", roll)

	// -- Start the server --

	address := fmt.Sprintf(":%d", *port)

	if err := e.Start(address); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

func health(c echo.Context) error {
	return c.String(http.StatusOK, "Toto")
}

func roll(c echo.Context) error {
	i := rollDice6()

	return c.String(http.StatusOK, fmt.Sprint(i))
}

func rollDice6() int {
	i := rand.Intn(6)
	i++

	return i
}
