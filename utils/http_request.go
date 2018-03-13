package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"net/http"
	"net/url"
)

type AuthenData struct {
	UserId		string
	Token		string
}

type HttpResponse struct {
	Data 		string
	Code		int
}

type HttpRequest struct {
	Method		string 			`json:""`
	Domain		string 			`json:""`
	Path		string 			`json:""`
	Body		*url.Values
	Authen 		*AuthenData
	Hosts		string
}

func (httpRequest *HttpRequest) MakeRequest() *HttpResponse {
	u, _ := url.ParseRequestURI(httpRequest.Domain)
	u.Path = httpRequest.Path
	url := u.String()

	inputBody := httpRequest.Body.Encode()
	r, _ := http.NewRequest(httpRequest.Method, url, strings.NewReader(inputBody))
	r.Header.Add("authorization", httpRequest.Authen.Token)
	r.Header.Add("host", httpRequest.Hosts)
	r.Header.Add("content-type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(httpRequest.Body.Encode())))
	
	resp, err := http.DefaultClient.Do(r)	

	if err != nil {
		log.Println(resp.StatusCode)
		return &HttpResponse{
			Code: resp.StatusCode,
		}
	}
	defer resp.Body.Close()
	
	httpResponse := &HttpResponse{
		Code: resp.StatusCode,
	}
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(resp)
	fmt.Println(string(body))
	if err != nil {
		httpResponse.Data = string(body)
	}
	return httpResponse
}