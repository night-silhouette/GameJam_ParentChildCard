package route

import (
	"fmt"
	"pcc_card/infra/config"
	"pcc_card/presentation/handler/user_handler"
	"pcc_card/presentation/response"

	"github.com/gin-gonic/gin"
)

var R *gin.Engine

func Init() {
	R = gin.Default()
	R.GET("/ping", func(c *gin.Context) {
		response.Success(c, "pong")
	})
}

func Run() {
	route_config := config.Read_route_info()
	err := R.Run(fmt.Sprintf("%s:%d", route_config.Ip, route_config.Port))
	if err != nil {
		fmt.Println("路由启动失败")
		panic(err)
	}
}

func Register_user_routes(h user_handler.User_handler) {
	v1_user := R.Group("/v1/user")
	v1_user.GET("/", h.Get())
	v1_user.POST("/", h.Post())
	v1_user.PATCH("/", h.Patch())
	v1_user.DELETE("/", h.Delete())
	v1_user.PUT("/", h.Put())
}
