package config

import (
	_ "embed"
	"encoding/json"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

//go:embed info/secret_key.json
var secret_key_file []byte

type Secret_key struct {
	Key string `json:"secret_key"`
}

func Read_secret_key() string {
	config := Secret_key{}
	err := json.Unmarshal(secret_key_file, &config)
	if err != nil {
		fmt.Println("secret_key file err")
		panic(err)
	}
	bytes, _ := bcrypt.GenerateFromPassword([]byte(config.Key), bcrypt.DefaultCost)
	return string(bytes)
}
