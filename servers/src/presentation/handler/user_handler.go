package handler

import (
	"pcc_card/application/service"

	"github.com/gin-gonic/gin"
)

type User_handler interface {
	Handler
}
type User_handler_impl struct {
	s service.User_service
}

func (u *User_handler_impl) Set_service(svc service.Service) {
	u.s = svc.(service.User_service)
}

func (u *User_handler_impl) Get() gin.HandlerFunc {
	//TODO implement me
	panic("implement me")
}

func (u *User_handler_impl) Post() gin.HandlerFunc {
	//TODO implement me
	panic("implement me")
}

func (u *User_handler_impl) Put() gin.HandlerFunc {
	//TODO implement me
	panic("implement me")
}

func (u *User_handler_impl) Delete() gin.HandlerFunc {
	//TODO implement me
	panic("implement me")
}

func (u *User_handler_impl) Patch() gin.HandlerFunc {
	//TODO implement me
	panic("implement me")
}
