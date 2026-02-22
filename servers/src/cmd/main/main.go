package main

import (
	"pcc_card/infra/db"
	"pcc_card/presentation/route"
)

func main() {
	db.Init()
	route.Init()
}
