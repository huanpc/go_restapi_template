package handler

import (
	"apistream/utils"
	"reflect"
	"strconv"
	"time"
	"strings"
	"log"
	"net/http"
	"apistream/view"
	"context"
	"apistream/storage"
	"apistream/config"
	"fmt"
	"net/url"
	"errors"
)

func storeDB (client storage.MySqlClient, input *storage.ChannelTable){
	table := storage.Table{
		Name: "channel",
		DateTimeColumns: []string{"time_start", "time_end"},
		NotNullColumns: []string{"channel_name", "channel_alias_name", "owner_id", "storage"},
		AutoUpdateDateTimeColumns: []string{"time_start", "time_end"},
	}
	// insert
	sqlIn:= storage.PrepareInsert(client.Client, table, input)
	sqlIn.ExecuteInsert()
}

func MapJsonKey2Field(dest interface{}) (key2Field map[string]string, keyNotNil []string){
	t := reflect.TypeOf(dest).Elem()
	for index := 0; index < t.NumField(); index++ {
		f := t.Field(index)
		jsonKey := f.Tag.Get("json")
		tokens := strings.Split(jsonKey, ",")
		key2Field[tokens[0]] = f.Name
		if !utils.IsStringSliceContains(tokens, "omitempty") {
			keyNotNil = append(keyNotNil, tokens[0])
		}
	}
	return
}

func ParseToObject(form url.Values, dest interface{}) (err error){
	errs := make([]string, 0)	
	destVal := reflect.ValueOf(dest).Elem()
	for k, v := range form {
		key2Field, keyNotNil := MapJsonKey2Field(dest)		
		val := destVal.FieldByName(key2Field[k])		
		if !val.IsNil() {
			switch val.Type().Kind() {
			case reflect.String:
				val.SetString(v[0])
			case reflect.Bool:
				if s, err := strconv.ParseBool(v[0]); err == nil {
					val.SetBool(s)
				}				
			case reflect.Float64:
				if s, err := strconv.ParseFloat(v[0], 64); err == nil {
					val.SetFloat(s)
				}				
			case reflect.Int64:
				if s, err := strconv.ParseInt(v[0], 10, 64); err == nil {
					val.SetInt(s)
				}
			}
		} else {
			if utils.IsStringSliceContains(keyNotNil, k) {
				errs = append(errs, k)
			}
		}		
	}
	if len(errs) > 0 {
		return errors.New(fmt.Sprintf("missing request params: [%s]", errs))
	}
	return nil
}


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

	// index to elastic search
	// parse to event
	event := &storage.OnConnectEvent{}
	err := ParseToObject(data, event)
	if err != nil {
		log.Panic(err)
		return view.BadRequest(err)
	}

	event.Created = time.Now().UTC()
	ctx := context.Background()
	client := storage.NewEsClient(fmt.Sprintf("%v:%v", cfg.ELS_HOST, cfg.ELS_PORT), ctx)
	item := &storage.Event {
		Id: strconv.FormatInt(time.Now().UnixNano(), 10),
		IndexName: "event",
		Type: "streaming_event",
		Data: event,
	}
	client.IndexItem(item)

	// store do db
	// parse to channel
	channel := &storage.ChannelTable{}
	err = ParseToObject(data, channel)
	if err != nil {
		log.Panic(err)
		return view.BadRequest(err)
	}
	mysqlClient := storage.NewMySqlClient(cfg.STREAM_DB_USERNAME, cfg.STREAM_DB_PASSWORD, cfg.STREAM_DB_NAME, cfg.STREAM_DB_HOST, cfg.STREAM_DB_PORT)
	storeDB(mysqlClient, channel)

	return view.Ok(bh.r.Form)
}