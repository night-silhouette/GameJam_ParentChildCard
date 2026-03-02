package main

import (
	"pcc_card/infra/db"
	"pcc_card/infra/repo"
)

func main() {
	DB := db.ConnectDB()
	user_repo := repo.New_repo[*repo.User_repo_impl](DB)
	
}
