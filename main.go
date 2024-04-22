package main

import (
	"github.com/JimmyTanapon/assessment-tax/postgres"
	"github.com/JimmyTanapon/assessment-tax/tax"
	"github.com/labstack/echo/v4"
)

func main() {
	p, err := postgres.New()
	if err != nil {
		panic(err)
	}

	e := echo.New()
	handler := tax.New(p)
	e.POST("/tax/calculations", handler.TaxHandler)
	e.Logger.Fatal(e.Start(":8080"))
}
