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

	// 仅仅修改了 query，利用 json_build_object 把字段打包
	// 这样数据库返回的就是一个完整的 JSON 字符串，刚好能塞进你的 info []byte
	query := `select json_build_object(
        'damage', damage, 
        'initHp', "initHp", 
        'maxHp', "maxHp", 
        'price', price, 
        'skillCharge', "skillCharge", 
        'skillcardUseNum', "skillcardUseNum", 
        'category', category
    ) from newcards where id = $1`

	data := db.QueryRowContext(ctx, query, ID)
	err := data.Scan(&info)
	if err != nil {
		log.Println(err)
	}
	err = json.Unmarshal(info, &res)
	if err != nil {
		log.Println(err)
	}
	if res["category"] == 1 || res["category"] == 2 {
		res["is_parent"] = true
	}
	if res["category"] == 3 || res["category"] == 4 {
		res["is_parent"] = false
	}

	return res
}
