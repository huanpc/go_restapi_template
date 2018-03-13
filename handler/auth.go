package handler

import (
	"strings"
	// "os"
	"log"
	// "github.com/Sirupsen/logrus"
	// "context"
	// "fmt"
	"net/http"

	"apistream/view"
	// "apistream/storage"
	"apistream/config"
	"apistream/utils"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}


// Authentication handler
func Auth(bh BaseHandler) view.ApiResponse {
	// Read config file
	cfg := config.AppConfig()

	bh.r.ParseForm()
	data := bh.r.Form
	log.Println(strings.Repeat("#", 20))
	log.Println("received data from nginx-rtmp")
	log.Println(data)
	log.Println(strings.Repeat("#", 20))
	if len(data) == 0 {
		return view.ApiResponse{Code: http.StatusBadRequest, Data: bh.r.Form}
	}

	if res, ok := data["token"]; !ok || len(data["token"]) == 0{
		log.Println(res)
		return view.ApiResponse{Code: http.StatusBadRequest, Data: bh.r.Form}
	}
	log.Println("Do request")
	// forward loopback to getkong to authenticate
	request := &utils.HttpRequest{
		Method: "POST",
		Domain: cfg.API_GATEWAY,
		Path: "/api-stream/apis/event",
		Body: &data,
		Authen: &utils.AuthenData{
			Token: "Bearer " + string(data["token"][0]),
		},
		Hosts: cfg.HOST_NAME,
	}
	response := request.MakeRequest()
	log.Println("Response " + response.Data)
	return view.ApiResponse{Code: response.Code, Data: bh.r.Form}
}
