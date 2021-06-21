package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	port := flag.Int("port", 8080, "HTTP port")
	flag.Parse()

	// -- Basic echo --

	e := echo.New()

	// -- Instrument echo --

	e.Use(middleware.Logger())

	p := prometheus.NewPrometheus("fate", nil)
	p.Use(e)

	// -- Add endpoints

	e.GET("/manage/health", health)

	rand.Seed(time.Now().UnixNano())
	e.GET("/roll", roll)
	e.GET("/flip", flip)
	e.GET("/magic-ball", magicBall)
	e.POST("/custom", custom)

	// -- Start the server --

	address := fmt.Sprintf(":%d", *port)

	if err := e.Start(address); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

func health(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

func roll(c echo.Context) error {
	facesStr := c.QueryParam("faces")

	if len(facesStr) == 0 {
		facesStr = fmt.Sprint(6)
	}

	faces, err := strconv.Atoi(facesStr)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	dicesStr := c.QueryParam("dices")

	if len(dicesStr) < 2 {
		// Only one dice to roll
		return c.String(http.StatusOK, fmt.Sprint(rollDice(faces)))
	}

	// Roll all the dices concurrently

	dices, err := strconv.Atoi(dicesStr)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	ch := make(chan int, dices)

	for i := 0; i < dices; i++ {
		go func() {
			ch <- rollDice(faces)
		}()
	}

	res := []string{}
	for i := 0; i < dices; i++ {
		res = append(res, fmt.Sprint(<-ch))
	}

	return c.String(http.StatusOK, strings.Join(res, ","))
}

func rollDice(faces int) int {
	i := rand.Intn(faces)
	i++

	return i
}

func flip(c echo.Context) error {
	res := rollDice(2)

	if res == 2 {
		return c.String(http.StatusOK, "Head")
	}

	return c.String(http.StatusOK, "Tails")
}

var magicBallAnswers = []string{
	"It is certain",
	"It is decidedly so",
	"Without a doubt",
	"Yes definitely",
	"You may rely on it",
	"As I see it yes",
	"Most likely",
	"Outlook good",
	"Yes",
	"Signs point to yes",
	"Reply hazy try again",
	"Ask again later",
	"Better not tell you now",
	"Cannot predict now",
	"Concentrate and ask again",
	"Don't count on it",
	"My reply is no",
	"My sources say no",
	"Outlook not so good",
	"Very doubtful",
}

func magicBall(c echo.Context) error {
	return c.String(http.StatusOK, randomElement(magicBallAnswers))
}

type payload struct {
	Outcomes []string `json:"outcomes"`
}

func custom(c echo.Context) error {
	var p payload
	if err := c.Bind(&p); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.String(http.StatusOK, randomElement(p.Outcomes))
}

func randomElement(s []string) string {
	if len(s) == 0 {
		return ""
	}

	k := rand.Intn(len(s))
	return s[k]
}
