package handler

import (
	"strings"
	"log"
	"net/http"
	"apistream/view"
	// "context"
	"apistream/storage"
	// "apistream/config"
	// "fmt"
)

func Event(bh BaseHandler) view.ApiResponse{
	// Read config file
	// cfg := config.AppConfig()

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
	2018/03/13 10:13:56 map[jwt:[eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiI3V2hKNk9CUFBSRWp1ZW9WRXYxT2p2SHZlb282c1lTUyJ9.2-CxE7GroMHqTTQ8jFyqD9w3smjwJ4zmQfUziZMwo0k] swfurl:[] pageurl:[] addr:[10.240.152.231] user_id:[141] channel:[xyz] call:[connect] token:[123] app:[live] flashver:[FMLE/3.0 (compatible; Lavf57.71] tcurl:[rtmp://10.240.152.231:1935/live?token=123&user_id=141&channel=xyz&jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiI3V2hKNk9CUFBSRWp1ZW9WRXYxT2p2SHZlb282c1lTUyJ9.2-CxE7GroMHqTTQ8jFyqD9w3smjwJ4zmQfUziZMwo0k] epoch:[493718001]]
	2018/03/13 10:13:56 ####################
	{"http_method":"POST","http_proto":"HTTP/1.0","http_scheme":"http","level":"info","msg":"request complete","remote_addr":"10.240.152.231:34740","resp_bytes_length":0,"resp_elasped_ms":0.595929,"resp_status":0,"ts":"Tue, 13 Mar 2018 03:13:56 UTC","uri":"http://10.60.150.116/apis/auth","user_agent":""}

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
	// ctx := context.Background()
	// client := storage.NewEsClient(fmt.Sprintf("http://%v:%v", cfg.ELS_HOST, cfg.ES_PORT), ctx)
	// item := &storage.Event {
	// 	Id: "2",
	// 	IndexName: "event",
	// 	Type: "streaming_event",
	// 	Data: event,
	// }
	// client.IndexItem(item)

	return view.Ok(bh.r.Form)
}