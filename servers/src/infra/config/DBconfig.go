package config

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

//go:embed info/DB_info.json
var db_info_file []byte

type DB_info struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Ip       string `json:"ip"`
	Port     int    `json:"port"`
	Database string `json:"database"`
}

func Read_db_info() DB_info {
	config := DB_info{}
	err := json.Unmarshal(db_info_file, &config)
	if err != nil {
		fmt.Println("DB info file err")
		panic(err)
	}
	return config
}
