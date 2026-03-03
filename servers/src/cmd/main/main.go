package main

import (
	"database/sql"
	"pcc_card/application/service"
	"pcc_card/infra/db"
	"pcc_card/infra/repo"
	"pcc_card/presentation/handler"
	"pcc_card/presentation/route"
)

func main() {
	DB := db.ConnectDB()
	defer DB.Close()
	route.Init()
	user(DB)
}

func user(DB *sql.DB) {
	user_repo := repo.New_repo[*repo.User_repo_impl](DB)
	user_service := service.New_service[*service.User_service_impl](user_repo)
	user_handler := handler.New_handler[*handler.User_handler_impl](user_service)
	route.Register_user_routes(user_handler)
}
