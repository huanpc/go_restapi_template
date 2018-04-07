package main

import(
	"apistream/storage"
	"apistream/config"
	"fmt"
	// "strings"
)

func main(){
	// index to elastic search
	cfg := config.AppConfig()
	client := storage.NewMySqlClient(cfg.STREAM_DB_USERNAME, cfg.STREAM_DB_PASSWORD, cfg.STREAM_DB_NAME, cfg.STREAM_DB_HOST, cfg.STREAM_DB_PORT)
	
	table := storage.Table{
		Name: "channel",
		DateTimeColumns: []string{"time_start", "time_end"},
		NotNullColumns: []string{"channel_name", "channel_alias_name", "owner_id", "storage"},
		AutoUpdateDateTimeColumns: []string{"time_start", "time_end"},
	}

	sql2 := storage.PrepareCount(client.Client, table)
	ret, err := sql2.ExecuteCount("SELECT count(*) as count FROM "+table.Name + " WHERE channel_name=?", "NgoC Trinh LiveStream 3")
	if err == nil && ret == 0 {
		// insert
		in := &storage.ChannelTable{
			ChannelName: "NgoC Trinh LiveStream 3",
			ChannelAliasName: "test22",
			OwnerId: 1234124,
			Password: "3w434fsd",
			Storage: "test.com",
		}
		sqlIn:= storage.PrepareInsert(client.Client, table, in)
		sqlIn.ExecuteInsert()
	}
	
	// select
	dest := &storage.ChannelTable{}
	res := make([]storage.ChannelTable, 0)	
	sql := storage.PrepareSelect(client.Client, table, dest)
	println(sql.ExecuteSelect(&res, "SELECT * FROM "+table.Name))
	// count
	sql3 := storage.PrepareCount(client.Client, table)
	ret2, err2 := sql3.ExecuteCount("SELECT count(*) as count FROM "+table.Name)
	fmt.Println(ret2,err2)

}
