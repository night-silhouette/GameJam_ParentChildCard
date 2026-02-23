package handler

import (
	"pcc_card/application/service"

	"github.com/gin-gonic/gin"
)

type User_handler struct {
	service service.User_service
}

func New_user_handler(service service.User_service) *User_handler {
	return &User_handler{service: service}
}

func (u User_handler) Register(group *gin.RouterGroup) {
	//TODO implement me
	panic("implement me")
}

func (u User_handler) Get() gin.HandlerFunc {
	//TODO implement me
	panic("implement me")
}

func (u User_handler) Post() gin.HandlerFunc {
	//TODO implement me
	panic("implement me")
}

func (u User_handler) Put() gin.HandlerFunc {
	//TODO implement me
	panic("implement me")
}

func (u User_handler) Delete() gin.HandlerFunc {
	//TODO implement me
	panic("implement me")
}

func (u User_handler) Patch() gin.HandlerFunc {
	//TODO implement me
	panic("implement me")
}
