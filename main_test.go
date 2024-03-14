package main

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/go-resty/resty/v2"
)

type ResponseTODO struct {
	Status       string `json:"string"`
	ErrorMessage any    `json:"error_message,omitempty"`
	Data         []Todo `json:"data"`
}

const BASE_URL = "http://localhost:3000"
const DEBUG = false

var accessToken = ""

func TestLogin(t *testing.T) {
	res, err := resty.New().
		R().
		SetBasicAuth("noval", "agung").
		Post(BASE_URL + "/api/login")
	if err != nil {
		t.Error(err)
	}

	if DEBUG {
		log.Println("raw response body", string(res.Body()))
	}

	resBody := Response{}
	err = json.Unmarshal(res.Body(), &resBody)
	if err != nil {
		t.Error(err)
	}

	accessToken = resBody.Data.(string)
}

func TestAddTodo(t *testing.T) {
	payload := Todo{
		Name:    "do math homework",
		Details: "due date is tomorrow! do not forget",
		Done:    false,
	}

	res, err := resty.New().
		R().
		SetAuthToken(accessToken).
		SetHeader("Content-type", "application/json").
		SetBody(payload).
		Post(BASE_URL + "/api/todo")
	if err != nil {
		t.Error(err)
	}

	if DEBUG {
		log.Println("raw response body", string(res.Body()))
	}

	resBody := ResponseTODO{}
	err = json.Unmarshal(res.Body(), &resBody)
	if err != nil {
		t.Error(err)
	}

	if len(resBody.Data) == 0 {
		t.Error("empty todo list")
	}
	if resBody.Data[0].ID == "" {
		t.Error("incorrect ID")
	}
	if resBody.Data[0].Name != payload.Name {
		t.Error("incorrect name")
	}
	if resBody.Data[0].Details != payload.Details {
		t.Error("incorrect details")
	}
	if resBody.Data[0].Done != payload.Done {
		t.Error("incorrect done")
	}
}

func TestGetTodo(t *testing.T) {
	res, err := resty.New().
		R().
		SetAuthToken(accessToken).
		Get(BASE_URL + "/api/todo")
	if err != nil {
		t.Error(err)
	}

	if DEBUG {
		log.Println("raw response body", string(res.Body()))
	}

	resBody := ResponseTODO{}
	err = json.Unmarshal(res.Body(), &resBody)
	if err != nil {
		t.Error(err)
	}

	if len(resBody.Data) == 0 {
		t.Error("empty todo list")
	}
	if resBody.Data[0].ID == "" {
		t.Error("incorrect ID")
	}
	if resBody.Data[0].Name == "" {
		t.Error("incorrect name")
	}
	if resBody.Data[0].Details == "" {
		t.Error("incorrect details")
	}
}

func TestUpdateTodo(t *testing.T) {
	res, err := resty.New().
		R().
		SetAuthToken(accessToken).
		Get(BASE_URL + "/api/todo")
	if err != nil {
		t.Error(err)
	}

	if DEBUG {
		log.Println("raw response body", string(res.Body()))
	}

	resBody := ResponseTODO{}
	err = json.Unmarshal(res.Body(), &resBody)
	if err != nil {
		t.Error(err)
	}

	sample := resBody.Data[0]
	sample.Done = true
	sample.Details = "updated details"

	res, err = resty.New().
		R().
		SetAuthToken(accessToken).
		SetHeader("Content-type", "application/json").
		SetBody(sample).
		Put(BASE_URL + "/api/todo/" + sample.ID)
	if err != nil {
		t.Error(err)
	}

	if DEBUG {
		log.Println("raw response body", string(res.Body()))
	}

	resBody = ResponseTODO{}
	err = json.Unmarshal(res.Body(), &resBody)
	if err != nil {
		t.Error(err)
	}

	if len(resBody.Data) == 0 {
		t.Error("empty todo list")
	}
	if resBody.Data[0].ID == "" {
		t.Error("incorrect ID")
	}
	if resBody.Data[0].Name != sample.Name {
		t.Error("incorrect name")
	}
	if resBody.Data[0].Details != sample.Details {
		t.Error("incorrect details")
	}
	if resBody.Data[0].Done != sample.Done {
		t.Error("incorrect done")
	}
}

func TestDeleteTodo(t *testing.T) {
	res, err := resty.New().
		R().
		SetAuthToken(accessToken).
		Get(BASE_URL + "/api/todo")
	if err != nil {
		t.Error(err)
	}

	if DEBUG {
		log.Println("raw response body", string(res.Body()))
	}

	resBody := ResponseTODO{}
	err = json.Unmarshal(res.Body(), &resBody)
	if err != nil {
		t.Error(err)
	}
	totalData := len(resBody.Data)

	sample := resBody.Data[0]

	res, err = resty.New().
		R().
		SetAuthToken(accessToken).
		Delete(BASE_URL + "/api/todo/" + sample.ID)
	if err != nil {
		t.Error(err)
	}

	if DEBUG {
		log.Println("raw response body", string(res.Body()))
	}

	resBody = ResponseTODO{}
	err = json.Unmarshal(res.Body(), &resBody)
	if err != nil {
		t.Error(err)
	}
	if xx := len(resBody.Data); xx == totalData {
		t.Error("error on deleting todo")
	}
}
