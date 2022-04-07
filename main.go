package main

import (
	"context"
	"fmt"
	"log"

	"github.com/apache/airflow-client-go/airflow"
	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/usace/wat-api/config"
	"github.com/usace/wat-api/handler"
)

func main() {
	conf := airflow.NewConfiguration()
	conf.Host = "localhost:8080"
	conf.Scheme = "http"
	cli := airflow.NewAPIClient(conf)

	cred := airflow.BasicAuth{
		UserName: "username",
		Password: "password",
	}
	ctx := context.WithValue(context.Background(), airflow.ContextBasicAuth, cred)
	fmt.Println(ctx)
	variable, resp, err := cli.VariableApi.GetVariable(ctx, "foo").Execute()
	fmt.Println(resp)
	if err != nil {
		fmt.Println("we had an error" + err.Error())
	} else {
		fmt.Println(variable)
	}

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
	private.POST("wat-api/compute", wHandler.Version) //needs to post the job config
	//log.Print("starting server on port " + cfg.AppPort)
	//log.Fatal(e.Start(":" + cfg.AppPort))

}
