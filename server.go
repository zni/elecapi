package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

type Config struct {
	DBUser string `json:"db_user"`
	DBPass string `json:"db_pass"`
	DBUri  string `json:"db_uri"`
	DBName string `json:"db_name"`
}

type Ping struct {
	Now int64 `json:"now"`
}

type Resistor struct {
	Value string `json:"value"`
	Type  string `json:"type"`
	Count int64  `json:"count"`
}

type Capacitor struct {
	Value string `json:"value"`
	Type  string `json:"type"`
	Count int64  `json:"count"`
}

func main() {
	e := echo.New()

	var config Config
	f, err := ioutil.ReadFile("config.json")
	if err != nil {
		e.Logger.Fatal(err)
	}
	json.Unmarshal(f, &config)

	conn := fmt.Sprintf(
		"postgresql://%s:%s@%s/%s?sslmode=disable",
		config.DBUser,
		config.DBPass,
		config.DBUri,
		config.DBName,
	)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		e.Logger.Fatal(err)
	}

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
	e.GET("/api/v1/resistors", func(c echo.Context) error {
		return listResistors(c, db)
	}).Name = "api.v1.resistors"
	e.POST("/api/v1/resistors", func(c echo.Context) error {
		return addResistors(c, db)
	}).Name = "api.v1.add-resistors"

	// Capacitor endpoints
	e.GET("/api/v1/capacitors", func(c echo.Context) error {
		return listCapacitors(c, db)
	}).Name = "api.v1.capacitors"
	e.POST("/api/v1/capacitors", func(c echo.Context) error {
		return addCapacitors(c, db)
	}).Name = "api.v1.add-capacitors"

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

func listResistors(c echo.Context, db *sql.DB) error {
	rows, err := db.Query("SELECT value, type, count FROM resistors")
	if err != nil {
		c.Echo().Logger.Fatal(err)
		return err
	}

	var resistors []Resistor
	var value string
	var type_ string
	var count int64
	var resistor *Resistor
	for rows.Next() {
		rows.Scan(&value, &type_, &count)
		resistor = new(Resistor)
		resistor.Value = value
		resistor.Type = type_
		resistor.Count = count
		resistors = append(resistors, *resistor)
	}
	rows.Close()

	return c.JSON(http.StatusOK, resistors)
}

func addResistors(c echo.Context, db *sql.DB) error {
	resistor := new(Resistor)
	if err := c.Bind(resistor); err != nil {
		c.Echo().Logger.Fatal(err)
		return err
	}

	stmt, err := db.Prepare("INSERT INTO resistors (value, type, count) VALUES (($1) , ($2) , ($3))")
	if err != nil {
		c.Echo().Logger.Fatal(err)
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(resistor.Value, resistor.Type, resistor.Count); err != nil {
		c.Echo().Logger.Fatal(err)
		return err
	}

	return c.JSON(http.StatusCreated, resistor)
}

func listCapacitors(c echo.Context, db *sql.DB) error {
	rows, err := db.Query("SELECT value, type, count FROM capacitors")
	if err != nil {
		c.Echo().Logger.Fatal(err)
		return err
	}

	var capacitors []Capacitor
	var value string
	var type_ string
	var count int64
	var capacitor *Capacitor
	for rows.Next() {
		rows.Scan(&value, &type_, &count)
		capacitor = new(Capacitor)
		capacitor.Value = value
		capacitor.Type = type_
		capacitor.Count = count
		capacitors = append(capacitors, *capacitor)
	}
	rows.Close()

	return c.JSON(http.StatusOK, capacitors)
}

func addCapacitors(c echo.Context, db *sql.DB) error {
	capacitor := new(Capacitor)
	if err := c.Bind(capacitor); err != nil {
		c.Echo().Logger.Fatal(err)
		return err
	}

	stmt, err := db.Prepare("INSERT INTO capacitors (value, type, count) VALUES (($1) , ($2) , ($3))")
	if err != nil {
		c.Echo().Logger.Fatal(err)
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(capacitor.Value, capacitor.Type, capacitor.Count); err != nil {
		c.Echo().Logger.Fatal(err)
		return err
	}

	return c.JSON(http.StatusCreated, capacitor)
}
