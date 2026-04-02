package config

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

//go:embed info/route_info.json
var route_info_file []byte

type Route_info struct {
	Ip   string `json:"ip"`
	Port int    `json:"port"`
}

func Read_route_info() Route_info {
	config := Route_info{}
	err := json.Unmarshal(route_info_file, &config)
	if err != nil {
		fmt.Println("route info file err")
		panic(err)
	}
	return config
}
