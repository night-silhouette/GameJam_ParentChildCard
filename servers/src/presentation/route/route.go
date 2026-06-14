package route

import (
	"fmt"
	"io"
	"os"
	"pcc_card/infra/config"
	"pcc_card/presentation/handler/battlehandler"
	"pcc_card/presentation/handler/tokenhandler"
	"pcc_card/presentation/handler/userhandler"
	"pcc_card/presentation/response"
	"strings"

	"github.com/gin-gonic/gin"
)

var R *gin.Engine

type FilterWriter struct {
	Output io.Writer
}

func (f *FilterWriter) Write(p []byte) (n int, err error) {
	s := string(p)
	if strings.Contains(s, "/ping") && strings.Contains(s, " 200 ") {
		return len(p), nil
	}
	return f.Output.Write(p)
}

func Init() {
	filter := &FilterWriter{Output: os.Stdout}
	gin.DefaultWriter = filter
	R = gin.Default()

	R.GET("/ping", func(c *gin.Context) {
		response.Success(c, "pong")
	})
	//-------------------------------------------------------------

}

func Run() {
	route_config := config.Read_route_info()
	err := R.Run(fmt.Sprintf("%s:%d", route_config.Ip, route_config.Port))
	if err != nil {
		fmt.Println("路由启动失败")
		panic(err)
	}
}

func Register_user_routes(h userhandler.User_handler) {
	v1_user := R.Group("/v1/user")
	v1_user.GET("/", h.Get())
	v1_user.POST("/", h.Post())
	v1_user.PATCH("/", h.Patch())
	v1_user.DELETE("/", h.Delete())
	v1_user.PUT("/", h.Put())
	v1_user.GET("/vague/", h.UserVagueSearch())
	v1_user.GET("/gold/", h.GoldGet())

	R.GET("/v1/time/", h.TimeSync())
	R.GET("/v1/debug/time/", h.TimeDebug())

	R.GET("v1/start_pack/", h.StarterPack())
	R.POST("v1/debug/Card/", h.DebugGiveCardByCardId())
	R.GET("v1/bags/", h.BagGet())
	R.POST("v1/bags/sell/", h.CardSell())
}

func Register_token_routes(h tokenhandler.Token_handler) {
	R.Use(h.Middleware_token_check())
	v1_user := R.Group("/v1/token")
	v1_user.GET("/", h.Get())
	v1_user.POST("/", h.Post())
	v1_user.PATCH("/", h.Patch())
	v1_user.DELETE("/", h.Delete())
	v1_user.PUT("/", h.Put())
}

func RegisterBattleWS(h battlehandler.BattleHandler) {
	R.GET("/v1/ws/", h.BattleWs())
	R.GET("/v1/debug/match_pool/", h.DebugGetMachData())
	R.GET("/v1/debug/battle_container", h.DebugBattleContainer())
}

func RegisterMailRoute(h userhandler.User_handler) {
	R.POST("/v1/mail/", h.SendMail())
	R.GET("/v1/mail/", h.GetAllOnePage())
	R.DELETE("v1/mail/All/", h.DeleteMailAll())
	R.DELETE("/v1/mail/", h.DeleteMailByMailId())
	R.POST("/v1/mail/status/", h.ChangeMailStatus())
	R.GET("/v1/mail/status/", h.GetMailStatus())
	R.POST("v1/mail/friendships/", h.ChangeFriendshipsRequest())
}

func RegisterFriendshipsRoutes(h userhandler.User_handler) {
	R.GET("v1/friendships/", h.GetFriendships())
	R.POST("v1/friendships/", h.CreateFriendship())
	R.DELETE("v1/friendships/", h.DeleteFriendships())
}
