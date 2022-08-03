package main

import (
	"bytes"
	"github.com/labstack/echo/v4"
	"html/template"
	"net/http"
	"time"
)

type Ping struct {
	Now int64 `json:"now"`
}

type Resistor struct {
	Value string `json:"value"`
	Count int64  `json:"count"`
}

type Capacitor struct {
	Value string `json:"value"`
	Type  string `json:"type"`
	Count int64  `json:"count"`
}

func main() {
	e := echo.New()

	// Utility and index endpoints
	e.GET("/ping", ping).Name = "ping"
	e.GET("/", func(c echo.Context) error {
		if c.Request().Header.Get("Content-Type") == "application/json" {
			return c.JSON(http.StatusOK, e.Routes())
		}

		html, err := generateIndex(e.Routes())
		if err != nil {
			return err
		}

		return c.HTML(http.StatusOK, html)
	}).Name = "index"

	// Resistor endpoints
	e.GET("/api/v1/resistors", listResistors).Name = "api.v1.resistors"

	e.Logger.Fatal(e.Start(":1323"))
}

func generateIndex(routes []*echo.Route) (string, error) {
	templates := []string{"templates/index.html"}
	t, err := template.New("index.html").ParseFiles(templates...)
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	err = t.Execute(&buffer, routes)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func ping(c echo.Context) error {
	p := Ping{Now: time.Now().Unix()}
	return c.JSON(http.StatusOK, p)
}

func listResistors(c echo.Context) error {
	// TODO Eventually get this from a database.
	resistors := []Resistor{
		Resistor{
			Value: "100k",
			Count: 100,
		},
		Resistor{
			Value: "1M",
			Count: 100,
		},
	}

	return c.JSON(http.StatusOK, resistors)
}
