package handler

import (
	"github.com/gin-gonic/gin"
)

type Handler interface {
	Register(group *gin.RouterGroup)
	Get() gin.HandlerFunc
	Post() gin.HandlerFunc
	Put() gin.HandlerFunc
	Delete() gin.HandlerFunc
	Patch() gin.HandlerFunc
}
