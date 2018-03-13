package config

import (
	"strings"
	"encoding/json"
	"os"
	"path/filepath"
)

type Configuration struct {
	HOST_ADDRESS		string
	HOST_NAME 			string
	Port     			string
	ELS_HOST     		string
	ES_PORT     		string
	API_GATEWAY     	string
	MYSQL_USERNAME     	string
	MYSQL_PASSWORD     	string
	MYSQL_DB     		string
}

func AppConfig() Configuration {	
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
    if err != nil {
        panic("Config dir is invalid")
	}
	runes := []rune(dir)	
	rootPath := string(runes[0:strings.Index(dir, "apistream")]) + "apistream"
	absPath:= rootPath + "/config/httpd.json"
	println(absPath)
	file, ok := os.Open(absPath)
	defer file.Close()

	if ok != nil {
		panic("Can't open config file: httpd.json")
	}

	decoder := json.NewDecoder(file)
	cfg := Configuration{}

	if ok := decoder.Decode(&cfg); ok != nil {
		panic(ok)
	}
	return cfg
}
