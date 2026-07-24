package userhandler

import (
	"fmt"
	"pcc_card/application/entity/BattleData"
	"pcc_card/application/entity/User_entity"
	"pcc_card/application/service"
	"pcc_card/application/service/UserService"
	"pcc_card/global"
	"pcc_card/presentation/handler"
	"pcc_card/presentation/response"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type User_handler interface {
	handler.Handler
	Get() gin.HandlerFunc
	Post() gin.HandlerFunc
	Put() gin.HandlerFunc
	Delete() gin.HandlerFunc
	Patch() gin.HandlerFunc
	GetAllOnePage() gin.HandlerFunc
	GetMailStatus() gin.HandlerFunc
	SendMail() gin.HandlerFunc
	ChangeMailStatus() gin.HandlerFunc
	DeleteMailByMailId() gin.HandlerFunc
	DeleteMailAll() gin.HandlerFunc
	UserVagueSearch() gin.HandlerFunc
	CreateFriendship() gin.HandlerFunc
	GetFriendships() gin.HandlerFunc
	DeleteFriendships() gin.HandlerFunc
	ChangeFriendshipsRequest() gin.HandlerFunc
	TimeSync() gin.HandlerFunc
	TimeDebug() gin.HandlerFunc
	GoldGet() gin.HandlerFunc

	DebugGiveCardByCardId() gin.HandlerFunc
	StarterPack() gin.HandlerFunc
	BagGet() gin.HandlerFunc
	CardSell() gin.HandlerFunc
	GetUserBattle() gin.HandlerFunc
	GetLoot() gin.HandlerFunc
	PostLoot() gin.HandlerFunc
	GoodsGet() gin.HandlerFunc
	GoodsPost() gin.HandlerFunc
	Refresh() gin.HandlerFunc
}

type CardSellDto struct {
	StuffIdList []int `json:"stuff_id_list"`
}

func (u *User_handler_impl) CardSell() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CardSellDto
		UserId := c.GetInt("id")
		if err := c.ShouldBindJSON(&req); err != nil {
			fmt.Println(err)
			response.Fail(c, global.ResponseInvalidReqParams)
			return
		}
		err := u.s.SellCard(c.Request.Context(), UserId, req.StuffIdList)
		if err != global.ResponseSuccess {
			response.Fail(c, err)
			return
		}
		response.Success(c, "卖出去了")
	}

}

type User_handler_impl struct {
	s UserService.User_service
}

func (u *User_handler_impl) GoldGet() gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := u.s.GetUserGold(c.Request.Context(), c.GetInt("id"))
		if err != global.ResponseSuccess {
			response.Fail(c, err)
		}
		response.Success(c, res)
	}
}

func (u *User_handler_impl) Set_service(svc service.Service) {
	u.s = svc.(UserService.User_service)
}

func (u *User_handler_impl) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		is_admin := c.GetBool("is_admin")
		user_id := c.GetInt("id")

		if !is_admin {
			// 传入 c.Request.Context()
			e, err := u.s.Find_user_by_id(c.Request.Context(), user_id)
			if err != global.ResponseSuccess {
				response.Fail(c, err)
				return
			}
			response.Success(c, e)
			return
		}
		var req UserSearchReqDto
		if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
			response.Fail(c, global.ResponseInvalidReqParams)
			return
		}

		if req.Id != 0 {
			// 传入 c.Request.Context()
			e, status := u.s.Find_user_by_id(c.Request.Context(), req.Id)
			if status == global.ResponseSuccess {
				response.Success(c, e)
				return
			} else {
				response.Fail(c, status)
				return
			}
		} else if req.Name != "" {
			// 传入 c.Request.Context()
			e, status := u.s.Find_user_by_name(c.Request.Context(), req.Name)
			if status == global.ResponseSuccess {
				response.Success(c, e)
				return
			} else {
				response.Fail(c, status)
				return
			}
		} else {
			response.Fail(c, global.ResponseInvalidReqParams)
			return
		}
	}
}

func (u *User_handler_impl) Post() gin.HandlerFunc {
	return func(c *gin.Context) {
		var data UserPostDto
		if err := c.ShouldBind(&data); err != nil {
			response.Fail(c, global.ResponseInvalidReqParams)
			return
		}
		if data.Name == "" || data.Password == "" {
			response.Fail(c, global.ResponseRequiredParamsMissing)
			return
		}
		e := &User_entity.User{Id: -1, Name: data.Name, Password: data.Password, Is_admin: false}
		// 传入 c.Request.Context()
		err := u.s.Create_user(c.Request.Context(), e)
		if err != global.ResponseSuccess {
			response.Fail(c, err)
			return
		}
		response.Success(c, "ok")
		return
	}
}

func (u *User_handler_impl) Put() gin.HandlerFunc {
	return func(c *gin.Context) {
		is_admin := c.GetBool("is_admin")
		user_id := c.GetInt("id")
		var req UserPutReqDto
		if !is_admin && user_id != req.Id {
			response.Fail(c, global.ResponseForbidden)
			return
		}
		if err := c.ShouldBind(&req); err != nil {
			response.Fail(c, global.ResponseInvalidReqParams)
			return
		}

		e := &User_entity.User{Id: req.Id, Name: req.Name, Password: req.Password}
		// 传入 c.Request.Context()
		err := u.s.Update_user(c.Request.Context(), e)
		if err != global.ResponseSuccess {
			response.Fail(c, err)
			return
		}
		response.Success(c, "ok")
		return
	}
}

func (u *User_handler_impl) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		is_admin := c.GetBool("is_admin")
		user_id := c.GetInt("id")

		var id int
		var req UserDeleteReqDto
		if err := c.ShouldBindQuery(&req); err != nil {
			response.Fail(c, global.ResponseInvalidReqParams)
			return
		}
		id = req.Id
		if !is_admin && id != 0 {
			response.Fail(c, global.ResponseForbidden)
			return
		}

		// 这里原本就是 c.Request.Context()，保持
		err := u.s.Delete_user(c.Request.Context(), user_id)

		if err != global.ResponseSuccess {
			response.Fail(c, err)
			return
		}
		response.Success(c, "ok")
		return
	}
}

func (u *User_handler_impl) Patch() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.GetInt("id")
		var req UserPatchDto
		if err := c.ShouldBind(&req); err != nil {
			response.Fail(c, global.ResponseInvalidReqParams)
			return
		}
		if req.Name == "" && req.Password == "" {
			response.Fail(c, global.ResponseInvalidReqParams)
			return
		}
		// 传入 c.Request.Context()
		e, err := u.s.Find_user_by_id(c.Request.Context(), user_id)
		if err != global.ResponseSuccess {
			response.Fail(c, err)
			return
		}
		if req.Name != "" {
			if req.Name == e.Name {
				response.Fail(c, global.ResponseInvalidReqParams)
				return
			}
			e.Name = req.Name
		}
		if req.Password != "" {
			// 传入 c.Request.Context()
			OK := u.s.Check_password(c.Request.Context(), user_id, req.Password)
			if OK == global.ResponseSuccess {
				response.Fail(c, global.ResponseInvalidReqParams)
				return
			} else {
				e.Password = req.Password
			}
		} else {
			e.Password = ""
		}
		// 传入 c.Request.Context()
		err = u.s.Update_user(c.Request.Context(), e)
		if err != global.ResponseSuccess {
			response.Fail(c, err)
			return
		}
		response.Success(c, "ok")
	}
}

func (u *User_handler_impl) UserVagueSearch() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := UserVagueSearchReq{}
		if err := c.ShouldBind(&req); err != nil {
			fmt.Println(err)
			return
		}
		// 传入 c.Request.Context()
		err, uList := u.s.UserSearch(c.Request.Context(), req.VagueName)
		if err != global.ResponseSuccess {
			response.Fail(c, err)
			return
		}
		response.Success(c, uList)
	}
}
func (u *User_handler_impl) BagGet() gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := u.s.GetBags(c.Request.Context(), c.GetInt("id"))
		if err != global.ResponseSuccess {
			response.Fail(c, err)
			return
		}
		response.Success(c, res)
	}
}

func (u *User_handler_impl) StarterPack() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := u.s.GiveInitCardBag(c.Request.Context(), c.GetInt("id"))
		if err != global.ResponseSuccess {
			response.Fail(c, err)
			return
		}
		response.Success(c, "给了")
	}
}

type DebugGiveCardByCardIdDto struct {
	CardId float64 `json:"card_id"`
}

func (u *User_handler_impl) DebugGiveCardByCardId() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req DebugGiveCardByCardIdDto
		UserId := c.GetInt("id")
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Fail(c, global.ResponseInvalidReqParams)
			fmt.Println(err)
			return
		}
		fmt.Println(req)
		err := u.s.GiveCardByCardId(c.Request.Context(), UserId, int(req.CardId))
		if err != global.ResponseSuccess {
			response.Fail(c, err)
			return
		}
		response.Success(c, "增加成功")
	}
}
func (u *User_handler_impl) GetUserBattle() gin.HandlerFunc {
	return func(c *gin.Context) {
		UserId := c.GetInt("id")
		err := u.s.IsInBattle(c.Request.Context(), UserId)
		if err == global.ResponseSuccess {
			response.Success(c, true)
			return
		} else {
			fmt.Println("查看是否在战斗不在err:", err)
			response.Success(c, false)
			return
		}
	}
}

func (u *User_handler_impl) GetLoot() gin.HandlerFunc {
	return func(c *gin.Context) {
		UserId := c.GetInt("id")
		err, LootDto := u.s.GetLoot(UserId, c.Request.Context())
		if err != global.ResponseSuccess {
			response.Fail(c, err)
			return
		}
		response.Success(c, LootDto)
	}
}

func (u *User_handler_impl) PostLoot() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req BattleData.LootDto
		UserId := c.GetInt("id")
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Fail(c, global.ResponseInvalidReqParams)
			fmt.Println(err)
			return
		}
		if len(req.Data) > 5 {
			response.Fail(c, global.BattleCardNumErr)
			return
		}

		err1 := u.s.CreateStuffByLootCardId(&req, UserId, c.Request.Context())
		if err1 != global.ResponseSuccess {
			response.Fail(c, err1)
			return
		}
		response.Success(c, "ok")
		return
	}
}

func (u *User_handler_impl) GoodsGet() gin.HandlerFunc {
	return func(c *gin.Context) {
		UserId := c.GetInt("id")
		err, List := u.s.GetGoods(UserId, c.Request.Context())
		if err != global.ResponseSuccess {
			response.Fail(c, err)
			return
		}
		if len(List) == 0 {
			err2 := u.s.CreateGoods(UserId, c.Request.Context())
			if err2 != global.ResponseSuccess {
				response.Fail(c, err2)
				return
			}
			_, List2 := u.s.GetGoods(UserId, c.Request.Context())
			response.Success(c, List2)
		} else {
			response.Success(c, List)
		}
	}
}
func (u *User_handler_impl) GoodsPost() gin.HandlerFunc {
	return func(c *gin.Context) {}
}
func (u *User_handler_impl) Refresh() gin.HandlerFunc {
	return func(c *gin.Context) {
		UserId := c.GetInt("id")
		err2 := u.s.CreateGoods(UserId, c.Request.Context())
		if err2 != global.ResponseSuccess {
			response.Fail(c, err2)
			return
		}
		response.Success(c, "ok")
	}
}
