package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/usace/wat-api/config"
	handler "github.com/usace/wat-api/handlers"
	"github.com/usace/wat-api/utils"
	"gopkg.in/yaml.v2"
)

func TestPostCompute(t *testing.T) {
	wHandler, err := handler.CreateWatHandlerFromConfig(mockLoader().Config())
	if err != nil {
		t.Fail()
	}
	sj := handler.MockStochastic2dJob(wHandler.Config())
	fmt.Println(sj)
	byteblob, err := json.Marshal(sj)
	if err != nil {
		t.Fail()
	}
	fmt.Println(string(byteblob))
	response, err := http.Post("http://host.docker.internal:8001/wat-api/compute", "application/json", bytes.NewBuffer(byteblob))
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println(response)
}
func mockLoader() utils.ServicesLoader {
	cfg := config.WatConfig{
		APP_PORT:              "8080",
		SKIP_JWT:              true,
		AWS_ACCESS_KEY_ID:     "AKIAIOSFODNN7EXAMPLE",
		AWS_SECRET_ACCESS_KEY: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
		AWS_DEFAULT_REGION:    "us-east-1",
		AWS_S3_REGION:         "us-east-1",
		AWS_REGION:            "us-east-1",
		AWS_BUCKET:            "cloud-wat-dev",
		S3_MOCK:               true,
		S3_BUCKET:             "configs",
		S3_ENDPOINT:           "http://host.docker.internal:9000",
		S3_DISABLE_SSL:        true,
		S3_FORCE_PATH_STYLE:   true,
		REDIS_HOST:            "cache",
		REDIS_PORT:            "6379",
		REDIS_PASSWORD:        "eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81",
		SQS_ENDPOINT:          "http://host.docker.internal:9324",
	}
	ldr, err := utils.InitLoaderWithConfig("", cfg)
	if err != nil {
		fmt.Print(err)
	}
	return ldr
}
func TestSerializeStochasticJob(t *testing.T) {
	wHandler, err := handler.CreateWatHandler()
	if err != nil {
		t.Fail()
	}
	sj := handler.MockStochasticJob(wHandler.Config())
	byteblob, err := yaml.Marshal(sj)
	if err != nil {
		t.Fail()
	}
	fmt.Println(string(byteblob))
}
