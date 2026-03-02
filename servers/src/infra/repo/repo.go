package repo

import (
	"database/sql"
	"reflect"
)

type Repo interface {
	Set_db(db *sql.DB)
}

func New_repo[T Repo](db *sql.DB) T {
	tType := reflect.TypeOf((*T)(nil)).Elem() // 拿到 *User_repo_impl 的类型信息
	newVal := reflect.New(tType.Elem())
	repo := newVal.Interface().(T)
	repo.Set_db(db)
	return repo
}
