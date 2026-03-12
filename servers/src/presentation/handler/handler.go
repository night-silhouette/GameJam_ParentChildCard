package handler

import (
	"pcc_card/application/service"
	"reflect"
)

type Handler interface {
	Set_service(service service.Service)
}

func New_handler[T Handler](service service.Service) T {
	tType := reflect.TypeOf((*T)(nil)).Elem() // 拿到 *User_repo_impl 的类型信息
	newVal := reflect.New(tType.Elem())
	h := newVal.Interface().(T)
	h.Set_service(service)
	return h
}
