package utils

import (
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

	client := &http.Client{}
	r, _ := http.NewRequest(httpRequest.Method, url, strings.NewReader(httpRequest.Body.Encode()))
	r.Header.Add("Authorization", httpRequest.Authen.Token)
	r.Header.Add("Host", httpRequest.Hosts)
	r.Header.Add("Content-Type", "application/x-wwww-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(httpRequest.Body.Encode())))

	resp, _ := client.Do(r)
	buffer := make([]byte, resp.ContentLength)
	resp.Body.Read(buffer)
	stringBuffer := string(buffer[:resp.ContentLength])
	return &HttpResponse{
		Code: resp.StatusCode,
		Data: stringBuffer,
	}
}