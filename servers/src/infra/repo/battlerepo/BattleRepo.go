package battlerepo

import (
	"database/sql"
	"pcc_card/infra/repo"
)

type BattleRepo interface {
	repo.Repo
}

type BattleRepoImpl struct {
	db *sql.DB
}

func (r *BattleRepoImpl) Set_db(db *sql.DB) {
	r.db = db
}
