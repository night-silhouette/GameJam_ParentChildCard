package route

import (
	"fmt"
	"pcc_card/application/service"
	"pcc_card/infra/config"
	"pcc_card/presentation/response"

	"github.com/gin-gonic/gin"
)

func Init() {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		response.Success(c, "pong")
	})
	user_v1 := r.Group("v1/user")

	route_config := config.Read_route_info()
	err := r.Run(fmt.Sprintf("%s:%d", route_config.Ip, route_config.Port))
	if err != nil {
		fmt.Println("路由启动失败")
		panic(err)
	}
}
