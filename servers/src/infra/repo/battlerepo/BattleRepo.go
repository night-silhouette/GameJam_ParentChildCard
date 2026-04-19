package battlerepo

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"pcc_card/infra/repo"

	"github.com/redis/go-redis/v9"
)

type BattleRepo interface {
	repo.Repo
	ReadCardByID(ctx context.Context, db repo.SQLQueryer, ID int) map[string]any
}

type BattleRepoImpl struct {
	db *sql.DB
	rd *redis.Client
}

func (r *BattleRepoImpl) Get_db() *sql.DB {
	return r.db
}

func (r *BattleRepoImpl) Set_db(db *sql.DB, rd *redis.Client) {
	r.db = db
	r.rd = rd
}

func (r *BattleRepoImpl) ReadCardByID(ctx context.Context, db repo.SQLQueryer, ID int) map[string]any {
	var info []byte
	var res map[string]any
	query := "select info from cards where id = $1"
	// 使用 QueryRowContext 传递 ctx
	data := db.QueryRowContext(ctx, query, ID)
	err := data.Scan(&info)
	if err != nil {
		log.Println(err)
	}
	err = json.Unmarshal(info, &res)
	if err != nil {
		log.Println(err)
	}
	return res
}
