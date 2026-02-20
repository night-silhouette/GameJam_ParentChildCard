package db

import (
	"database/sql"
	"fmt"
	"pcc_card/infra/config"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var DB *sql.DB

func Init() {
	ConnectDB()
}

func ConnectDB() {
	DB_info := config.Read_db_info()
	var dsn = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", DB_info.User, DB_info.Password, DB_info.Ip, DB_info.Port, DB_info.Database)
	var err error
	DB, err = sql.Open("pgx", dsn)
	if err != nil {
		panic(err)
	}
	err = DB.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected to db")
}
