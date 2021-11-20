package main

import (
	"gocommerce/controllers"
	"gocommerce/core"

	_ "gocommerce/docs"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func main() {
	defer core.App.Close()

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"*"},
		AllowHeaders: []string{"*"},
	}))
	e.GET("/docs/*", echoSwagger.WrapHandler)
	api := e.Group("/v1")
	{
		users := api.Group("/users")
		// mid.SetClientJWTmiddlewares(users, "user")
		users.GET("", controllers.UserList)
		users.POST("/create", controllers.UserCreate)
		users.PATCH("/update/:id", controllers.UserUpdate)
		users.DELETE("/delete/:id", controllers.UserDelete)
	}
	var host string
	if core.App.Config.DB_NAME == "db_example" {
		host = "127.0.0.1"
	}

	e.Logger.Fatal(e.Start(host + ":" + core.App.Port))
	os.Exit(0)

}
