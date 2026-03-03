package handler

import (
	"pcc_card/application/service"
	"reflect"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	Set_service(service service.Service)
	Get() gin.HandlerFunc
	Post() gin.HandlerFunc
	Put() gin.HandlerFunc
	Delete() gin.HandlerFunc
	Patch() gin.HandlerFunc
}

func New_handler[T Handler](service service.Service) T {
	tType := reflect.TypeOf((*T)(nil)).Elem() // 拿到 *User_repo_impl 的类型信息
	newVal := reflect.New(tType.Elem())
	h := newVal.Interface().(T)
	h.Set_service(service)
	return h
}
