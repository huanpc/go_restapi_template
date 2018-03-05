package handler

import (
	"apistream/view"
	"log"
	"net/http"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Authentication handler
func Auth(bh BaseHandler) view.ApiResponse {
	bh.r.ParseForm()
	data := bh.r.Form
	log.Println("------------------------")
	log.Println(data)
	log.Println("------------------------")
	return view.ApiResponse{Code: http.StatusOK, Data: bh.r.Form}
}
