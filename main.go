package main

import (
	"fmt"
	"log"

	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/usace/wat-api/config"
	"github.com/usace/wat-api/handler"
)

func main() {
	var cfg config.Config
	if err := envconfig.Process("watapi", &cfg); err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(cfg)
	cfg.SkipJWT = true

	wHandler := handler.CreateWatHandler()
	e := echo.New()
	private := e.Group("")
	public := e.Group("")
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Public Routes
	public.GET("wat-api/version", wHandler.Version)
	public.GET("wat-api/plugins", wHandler.Plugins)
	//Private Routes
	private.POST("wat-api/compute", wHandler.Version)
	log.Print("starting server on port " + cfg.AppPort)
	log.Fatal(e.Start(":" + cfg.AppPort))

}
