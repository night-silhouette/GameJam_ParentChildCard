package service

import (
	"pcc_card/infra/repo"
	"reflect"
)

type Service interface {
	Set_repo(repo repo.Repo)
}

func New_service[T Service](repo repo.Repo) T {
	tType := reflect.TypeOf((*T)(nil)).Elem() // 拿到 *User_repo_impl 的类型信息
	newVal := reflect.New(tType.Elem())
	service := newVal.Interface().(T)
	service.Set_repo(repo)
	return service
}
