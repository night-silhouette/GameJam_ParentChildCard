package battlerepo

import (
	"database/sql"
	"encoding/json"
	"pcc_card/infra/repo"

	"github.com/redis/go-redis/v9"
)

type BattleRepo interface {
	repo.Repo
	ReadCardByID(ID int) map[string]any
}

type BattleRepoImpl struct {
	db *sql.DB
	rd *redis.Client
}

func (r *BattleRepoImpl) Set_db(db *sql.DB, rd *redis.Client) {
	r.db = db
	r.rd = rd
}

func (r *BattleRepoImpl) ReadCardByID(ID int) map[string]any {
	var info []byte
	var res map[string]any
	query := "Select info from cards where id = $1"
	data := r.db.QueryRow(query, ID)
	err := data.Scan(&info)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal([]byte(info), &res)
	if err != nil {
		panic(err)
	}
	return res
}
