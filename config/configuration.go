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
	ELS_PORT     		string
	API_GATEWAY     	string
	STREAM_DB_HOST     	string
	STREAM_DB_PORT     	string
	STREAM_DB_NAME     		string
	STREAM_DB_USERNAME     		string
	STREAM_DB_PASSWORD     		string
}

func AppConfig() Configuration {	
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
    if err != nil {
        panic("Config dir is invalid")
	}
	runes := []rune(dir)
	println(dir)
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
