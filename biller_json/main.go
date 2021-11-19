package main

import (
	"fmt"

	"github.com/adesupraptolaia/demo-biller/biller_json/controller"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	// Echo instance
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.BodyDumpWithConfig(middleware.BodyDumpConfig{
		Handler: func(c echo.Context, req []byte, res []byte) {
			fmt.Printf("request ======\n%s\n======\n", string(req))
			fmt.Printf("response %d ======\n%s\n======\n", c.Response().Status, string(res))

		},
	}))

	// Routes
	e.POST("/signature", controller.Signature)

	api := e.Group("api", controller.Middleware)

	api.POST("/inquiry", controller.Inquiry)
	api.POST("/purchase", controller.Purchase)
	api.POST("/advice", controller.Advice)

	// Start server
	e.Logger.Fatal(e.Start(":9000"))
}
