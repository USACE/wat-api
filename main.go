package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	handler "github.com/usace/wat-api/handlers"
)

func main() {
	wHandler, err := handler.CreateWatHandler()
	if err != nil {
		log.Fatal(err.Error())
	}
	e := echo.New()
	private := e.Group("")
	public := e.Group("")
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// Public Routes
	public.GET("wat-api/version", wHandler.Version)
	public.GET("wat-api/plugins", wHandler.Plugins)
	//Private Routes
	private.GET("wat-api/compute", wHandler.ExecuteJob) //needs to be a post and post the job config
	log.Print("starting server on port " + wHandler.AppPort)
	log.Fatal(e.Start(":" + wHandler.AppPort))
}
