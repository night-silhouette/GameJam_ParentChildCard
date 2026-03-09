package main

import (
	"database/sql"
	"pcc_card/application/service"
	"pcc_card/application/service/UserService"
	"pcc_card/infra/db"
	"pcc_card/infra/repo"
	"pcc_card/presentation/handler"
	"pcc_card/presentation/handler/token_handler"
	"pcc_card/presentation/handler/user_handler"
	"pcc_card/presentation/route"
)

func main() {
	DB := db.ConnectDB()
	defer DB.Close()
	route.Init()
	user(DB)

	route.Run()
}

func user(DB *sql.DB) {
	user_repo := repo.New_repo[*repo.User_repo_impl](DB)
	user_service := service.New_service[*UserService.User_service_impl](user_repo)
	user_handler := handler.New_handler[*user_handler.User_handler_impl](user_service)
	token_handler := handler.New_handler[*token_handler.Token_handler_impl](user_service)
	route.Register_token_routes(token_handler)
	route.Register_user_routes(user_handler)
}
