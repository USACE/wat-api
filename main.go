package main

import (
	"fmt"
	"log"

	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/usace/wat-api/handler"
	"github.com/usace/wat-api/model"
)

func main() {
	var cfg model.Config
	if err := envconfig.Process("watapi", &cfg); err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(cfg)
	cfg.SkipJWT = true

	wHandler := handler.CreateWatHandler()
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Public Routes
	e.GET("wat-api/version", wHandler.Version)
	//Private Routes

	log.Print("starting server on port " + cfg.AppPort)
	log.Fatal(e.Start(":" + cfg.AppPort))

}
