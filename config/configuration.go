package config

import "encoding/json"
import "os"
import "path/filepath"

type Configuration struct {
	HostName string
	Port     string
}

func AppConfig() Configuration {
	absPath, _ := filepath.Abs("httpd.json")
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
