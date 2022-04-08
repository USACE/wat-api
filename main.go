package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/apache/airflow-client-go/airflow"
)

func main() {

	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		rw.Write([]byte(`{"key": "food", "value": "bar"}`))
	}))
	// Close the server when test finishes
	defer server.Close()
	url, err := url.Parse(server.URL)
	//assert.Equal(t, nil, err)

	conf := airflow.NewConfiguration()
	conf.Host = url.Host //"host.docker.internal:8080"
	conf.Scheme = "http"
	cli := airflow.NewAPIClient(conf)

	cred := airflow.BasicAuth{
		UserName: "username",
		Password: "password",
	}

	//foo := "foo"
	//bar := "bar"
	ctx := context.WithValue(context.Background(), airflow.ContextBasicAuth, cred)
	//fmt.Println(ctx)
	//variable := airflow.Variable{}

	//variable.Key = &foo
	//variable.Value = &bar
	//cli.VariableApi.PostVariables(ctx).Variable(variable)
	variable, _, err := cli.VariableApi.GetVariable(ctx, "abc").Execute()
	//fmt.Println(resp)
	if err != nil {
		fmt.Println("we had an error" + err.Error())
	} else {
		//fmt.Println(resp)
		fmt.Println(variable.GetKey())
		fmt.Println(variable.GetValue())
	}
	/*
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
	*/
}
