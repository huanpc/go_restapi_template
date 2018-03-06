package handler

import (
	// "os"
	"apistream/view"
	"log"
	// "github.com/Sirupsen/logrus"
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

	// var logg = logrus.New()
	// logg.Out = os.Stdout
	// logg.Formatter = &logrus.JSONFormatter{
	// 	DisableTimestamp: true,
	// }
	// file, err := os.OpenFile("data.log", os.O_CREATE|os.O_APPEND, 0666)
  	// if err == nil {
	// 	logg.Out = file
  	// } else {
	// 	logg.Info("Failed to log to file, using default stderr")
	// }	
	
	// logg.Info(data)
	return view.ApiResponse{Code: http.StatusOK, Data: bh.r.Form}
}
