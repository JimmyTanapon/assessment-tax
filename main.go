package main

import (
	"github.com/JimmyTanapon/assessment-tax/tax"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.POST("/tax/calculations", tax.TaxHanler)
	e.Logger.Fatal(e.Start(":8080"))
}
