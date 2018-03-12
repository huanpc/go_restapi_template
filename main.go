package main

import(
	"apistream/utils"
	"net/url"
	"fmt"
)

func main(){
	request := &utils.HttpRequest{
		Method: "GET",
		Domain: "http://10.240.152.228:5001",
		Path: "/chats/contact",
		Body: &url.Values{},
		Authen: &utils.AuthenData{
			Token: "Bearer eyJhbGciOiJIUzI1NiJ9.eyJleHAiOjE1MzYwNTA4NDQsInVzZXIiOiJ7XCJpZFwiOlwiOTcyNDk2NzU5MTEwMTYzMzFcIixcInJvbGVzXCI6W1wiVXNlclwiXX0ifQ.zxSl1GqAFa115CV56zRmceq8x9LIfs8U0HljN1M5UBU",
		},
	}
	response := request.MakeRequest()
	fmt.Println(response.Code)
	fmt.Println(response.Data)
}

```
curl -i -X POST --url http://localhost:8001/apis/ --data 'name=dev.apistream.ereka'   --data 'hosts=dev.apistream.ereka.vn'   --data 'upstream_url=http://10.60.150.116:8080/apis/auth'
```