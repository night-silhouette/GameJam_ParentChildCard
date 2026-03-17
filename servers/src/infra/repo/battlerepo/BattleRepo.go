package battlerepo

import (
	"database/sql"
	"fmt"
	"pcc_card/application/entity/Card/CardAbstract"
	"pcc_card/infra/repo"
)

type BattleRepo interface {
	repo.Repo
	InsertCard(data CardAbstract.Card)
}

type BattleRepoImpl struct {
	db *sql.DB
}

func (r *BattleRepoImpl) Set_db(db *sql.DB) {
	r.db = db
}

func (r *BattleRepoImpl) InsertCard(data CardAbstract.Card) {
	query := "insert into cards (id, info) values ($1, $2)" +
		"on conflict (id) do update set info = cards.info || excluded.info"
	_, err := r.db.Exec(query, data.GetID(), data.GetInfo())
	if err != nil {
		fmt.Println(err)
	}
}
