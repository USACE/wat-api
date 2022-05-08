package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	handler "github.com/usace/wat-api/handlers"
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
	//defer response.Body.Close()
	/*
		b, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s", b)
	*/
}
