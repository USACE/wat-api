package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	handler "github.com/usace/wat-api/handlers"
	"gopkg.in/yaml.v2"
)

func TestPostCompute(t *testing.T) {
	wHandler, err := handler.CreateWatHandler()
	if err != nil {
		t.Fail()
	}
	sj := handler.MockStochasticJob(wHandler.Config())
	byteblob, err := json.Marshal(sj)
	if err != nil {
		t.Fail()
	}
	response, err := http.Post("http://host.docker.internal:8001/wat-api/compute", "application/json", bytes.NewBuffer(byteblob))
	if err != nil {
		t.Fail()
	}
	fmt.Println(response)
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
