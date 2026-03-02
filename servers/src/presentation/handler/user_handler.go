package handler

import (
	"pcc_card/application/service"

	"github.com/gin-gonic/gin"
)

type User_handler_impl struct {
	service service.User_service_impl
}

func New_user_handler(service service.User_service_impl) *User_handler_impl {
	return &User_handler_impl{service: service}
}

func (u User_handler_impl) Register(group *gin.RouterGroup) {
	//TODO implement me
	panic("implement me")
}

func (u User_handler_impl) Get() gin.HandlerFunc {
	//TODO implement me
	panic("implement me")
}

func (u User_handler_impl) Post() gin.HandlerFunc {
	//TODO implement me
	panic("implement me")
}

func (u User_handler_impl) Put() gin.HandlerFunc {
	//TODO implement me
	panic("implement me")
}

func (u User_handler_impl) Delete() gin.HandlerFunc {
	//TODO implement me
	panic("implement me")
}

func (u User_handler_impl) Patch() gin.HandlerFunc {
	//TODO implement me
	panic("implement me")
}
