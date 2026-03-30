package repo

import (
	"database/sql"
	"reflect"

	"github.com/redis/go-redis/v9"
)

type Repo interface {
	Set_db(db *sql.DB, rd *redis.Client)
}

func New_repo[T Repo](db *sql.DB, rd *redis.Client) T {
	tType := reflect.TypeOf((*T)(nil)).Elem() // 拿到 *User_repo_impl 的类型信息
	newVal := reflect.New(tType.Elem())
	repo := newVal.Interface().(T)
	repo.Set_db(db, rd)
	return repo
}
