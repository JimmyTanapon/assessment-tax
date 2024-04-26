package main

import (
	"os"

	"github.com/JimmyTanapon/assessment-tax/postgres"
	"github.com/JimmyTanapon/assessment-tax/tax"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type User struct {
	ID       int
	Username string
	Password string
	Role     string
}

func AuthMiddleware(username, password string, c echo.Context) (bool, error) {

	if username != os.Getenv("ADMIN_USERNAME") || password != os.Getenv("ADMIN_PASSWORD") {

		return false, nil
	}
	return true, nil
}

func main() {
	p, err := postgres.New()
	if err != nil {
		panic(err)
	}
	godotenv.Load()

	e := echo.New()
	handler := tax.New(p)
	e.POST("/tax/calculations", handler.TaxHandler)
	e.POST("tax/calculations/upload-csv", handler.TaxCSVUploadHandler)
	e.POST("/admin/deductions/:type", handler.UpDeductionHandler, middleware.BasicAuth(AuthMiddleware))

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
