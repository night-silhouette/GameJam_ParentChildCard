package main

import (
	"database/sql"
	"fmt"
	"pcc_card/application/service"
	"pcc_card/application/service/UserService"
	"pcc_card/application/service/battleservice"
	"pcc_card/global"
	"pcc_card/infra/db"
	"pcc_card/infra/repo"
	"pcc_card/infra/repo/battlerepo"
	"pcc_card/infra/repo/userrepo"
	"pcc_card/presentation/handler"
	"pcc_card/presentation/handler/battlehandler"
	"pcc_card/presentation/handler/tokenhandler"
	"pcc_card/presentation/handler/userhandler"
	"pcc_card/presentation/route"
)

func main() {
	if global.Isdebug == "debug" {
		fmt.Println("Start Debug")
	} else {
		fmt.Println("Production")
	}

	DB := db.ConnectDB()
	defer DB.Close()
	route.Init()
	user(DB)
	Battle(DB)

	battleservice.BC.AddBattle(2, 3)
	bt := battleservice.BC.GetBattle(3)
	fmt.Println("---------------------------------------------------")
	fmt.Println(bt.BattleID)
	fmt.Println(bt.Ctx.DataA.ID)
	fmt.Println(bt.Ctx.DataA.CardInHand)
	fmt.Println(bt.Ctx.DataB.CardInHand)
	route.Run()

}

func user(DB *sql.DB) {
	user_repo := repo.New_repo[*userrepo.User_repo_impl](DB)
	user_service := service.New_service[*UserService.User_service_impl](user_repo)
	user_handler := handler.New_handler[*userhandler.User_handler_impl](user_service)
	token_handler := handler.New_handler[*tokenhandler.Token_handler_impl](user_service)
	route.Register_token_routes(token_handler)
	route.Register_user_routes(user_handler)

}

func Battle(DB *sql.DB) {
	BattleRepo := repo.New_repo[*battlerepo.BattleRepoImpl](DB)
	BattleService := service.New_service[*battleservice.BattleServiceImpl](BattleRepo)
	BattleHandler := handler.New_handler[*battlehandler.BattleHandlerImpl](BattleService)
	battleservice.NewMatchManager()
	battleservice.InitBattleContainer()
	route.RegisterBattleRoutes(BattleHandler)
	battleservice.InitCardList(BattleService)

}
