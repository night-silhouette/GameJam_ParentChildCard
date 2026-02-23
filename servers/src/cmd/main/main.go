package main

import (
	"pcc_card/application/service"
	"pcc_card/infra/db"
	"pcc_card/presentation/route"
)

func main() {
	db.Init()
	service.Init()
	route.Init()
}
