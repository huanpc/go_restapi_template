package storage

import (
	"apistream/config"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type LiveTable struct {
	
}

type MySqlClient struct {
	Client *sql.DB
}

func NewMySqlClient(config config.Configuration) MySqlClient {
	client, err := sql.Open("mysql", fmt.Sprintf("%v:%v@/%v", config.MYSQL_USERNAME, config.MYSQL_PASSWORD, config.MYSQL_DB))
	if err != nil {
		panic(err.Error())
	}
	defer client.Close()
	return MySqlClient{
		Client: client,
	}
}

func (client *MySqlClient) ExecuteInsert(table string, values string) {
	// Prepare statement for inserting data
	stmtIns, err := client.Client.Prepare("INSERT INTO squareNum VALUES( ?, ? )")
	if err != nil {
		panic(err.Error())
	}
	defer stmtIns.Close()

	_, er := stmtIns.Exec("", "")
	if er != nil {
		panic(er.Error())
	}
}
