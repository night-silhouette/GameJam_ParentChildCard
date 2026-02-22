package route

import (
	"fmt"
	"pcc_card/infra/config"

	"github.com/gin-gonic/gin"
)

func Init() {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {})
	user_v1 := r.Group("v1/user")
	group_user(user_v1)

	route_config := config.Read_route_info()
	err := r.Run(fmt.Sprintf("%s:%d", route_config.Ip, route_config.Port))
	if err != nil {
		fmt.Println("路由启动失败")
		panic(err)
	}
}

func group_user(r *gin.RouterGroup) {
	r.GET("", func(c *gin.Context) {})
}
