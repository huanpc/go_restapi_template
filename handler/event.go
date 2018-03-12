package handler

import (
	"strings"
	"log"
	"net/http"
	"apistream/view"
	"context"
	"apistream/storage"
	"apistream/config"
	"fmt"
)

func Event(bh BaseHandler) view.ApiResponse{
	// Read config file
	cfg := config.AppConfig()

	bh.r.ParseForm()
	data := bh.r.Form
	log.Println(strings.Repeat("#", 20))
	log.Println("received data from authen service")
	log.Println(data)
	log.Println(strings.Repeat("#", 20))
	if len(data) == 0 {
		return view.ApiResponse{Code: http.StatusBadRequest, Data: bh.r.Form}
	}

	/* 
	map[token:[123] user_id:[141\channel_1] jwt:[eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiI3V2hKNk9CUFBSRWp1ZW9WRXYxT2p2SHZlb282c1lTUyJ9.2-CxE7GroMHqTTQ8jFyqD9w3smjwJ4zmQfUziZMwo0k]]
	2018/03/10 17:50:29 ------------------------
	{"http_method":"POST","http_proto":"HTTP/1.1","http_scheme":"http","level":"info","msg":"request complete","remote_addr":"10.60.150.117:34066","resp_bytes_length":215,"resp_elasped_ms":0.430467,"resp_status":200,"ts":"Sat, 10 Mar 2018 10:50:29 UTC","uri":"http://10.60.150.116:8080/event?token=123\u0026user_id=141%5Cchannel_1\u0026jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiI3V2hKNk9CUFBSRWp1ZW9WRXYxT2p2SHZlb282c1lTUyJ9.2-CxE7GroMHqTTQ8jFyqD9w3smjwJ4zmQfUziZMwo0k","user_agent":"PostmanRuntime/6.1.6"}
	*/

	/* 
		2018/03/08 09:30:02 map[swfurl:[] tcurl:[rtmp://10.240.152.231:1935/ereka_live] pageurl:[] addr:[10.240.152.233] epoch:[59093425] call:[connect] app:[ereka_live] flashver:[FMLE/3.0 (compatible; Lavf57.71]]
		2018/03/08 09:30:02 ------------------------
		{"http_method":"POST","http_proto":"HTTP/1.0","http_scheme":"http","level":"info","msg":"request complete","remote_addr":"10.240.152.231:60404","resp_bytes_length":235,"resp_elasped_ms":0.23192,"resp_status":200,"ts":"Thu, 08 Mar 2018 02:30:02 UTC","uri":"http://dev.apistream.ereka.vn/apis/auth","user_agent":""}
		2018/03/08 09:30:02 ------------------------
		2018/03/08 09:30:02 map[swfurl:[] tcurl:[rtmp://10.240.152.231:1935/ereka_live] pageurl:[] addr:[10.240.152.231] epoch:[59093560] call:[connect] app:[ereka_live] flashver:[LNX 9,0,124,2]]
		2018/03/08 09:30:02 ------------------------
		{"http_method":"POST","http_proto":"HTTP/1.0","http_scheme":"http","level":"info","msg":"request complete","remote_addr":"10.240.152.231:60408","resp_bytes_length":217,"resp_elasped_ms":0.200713,"resp_status":200,"ts":"Thu, 08 Mar 2018 02:30:02 UTC","uri":"http://dev.apistream.ereka.vn/apis/auth","user_agent":""}
	*/
	event := &storage.OnConnectEvent{}
	if data["call"] != nil {
		event.Call = data["call"][0]
	}
	if data["tc_url"] != nil {
		event.Call = data["tc_url"][0]
	}
	if data["addr"] != nil {
		event.Call = data["addr"][0]
	}
	if data["app"] != nil {
		event.Call = data["app"][0]
	}
	if data["flash_ver"] != nil {
		event.Call = data["flash_ver"][0]
	}
	if data["page_url"] != nil {
		event.Call = data["page_url"][0]
	}
	// index to elastic search
	ctx := context.Background()
	client := storage.NewEsClient(fmt.Sprintf("http://%v:%v", cfg.ELS_HOST, cfg.ES_PORT), ctx)
	item := &storage.Event {
		Id: "2",
		IndexName: "event",
		Type: "streaming_event",
		Data: event,
	}
	client.IndexItem(item)

	return view.Ok(bh.r.Form)
}