package UserService

import (
	"context"
	"fmt"
	"pcc_card/Util"
	"pcc_card/application/entity/BattleData"
	"pcc_card/application/entity/User_entity"
	"pcc_card/application/entity/mail"
	"pcc_card/application/service"
	"pcc_card/global"
	"pcc_card/infra/repo"
	"pcc_card/infra/repo/userrepo"

	"golang.org/x/crypto/bcrypt"
)

type User_service interface {
	service.Service
	// 新增 ctx 参数
	Find_user_by_name(ctx context.Context, name string) (*User_entity.User, global.ResponseStatusCode)
	Find_user_by_id(ctx context.Context, id int) (*User_entity.User, global.ResponseStatusCode)
	Check_password(ctx context.Context, id int, password string) global.ResponseStatusCode
	Release_token(userID int, ctx context.Context) (string, global.ResponseStatusCode)
	Is_valid_token(tokenString string, ctx context.Context) (int, bool, global.ResponseStatusCode)
	Create_user(ctx context.Context, user *User_entity.User) global.ResponseStatusCode
	ChangeMailStatus(ctx context.Context, AcceptId int, MailId []int, status int) global.ResponseStatusCode
	Update_user(ctx context.Context, e *User_entity.User) global.ResponseStatusCode
	Delete_user(ctx context.Context, id int) global.ResponseStatusCode
	GetAllOnePage(ctx context.Context, AcceptId int, page int) ([]*mail.Mail, global.ResponseStatusCode)
	GetMailStatus(ctx context.Context, id int) (int, global.ResponseStatusCode)
	SendMail(ctx context.Context, SendId int, body string, AcceptId int, Category string) global.ResponseStatusCode
	DeleteMailByMailId(ctx context.Context, MailId []int, AcceptId int) global.ResponseStatusCode
	DeleteMailAll(ctx context.Context, AcceptId int) global.ResponseStatusCode
	UserSearch(ctx context.Context, NameVague string) (global.ResponseStatusCode, []*User_entity.User)
	CreateFriendships(ctx context.Context, id1 int, id2 int) global.ResponseStatusCode
	DeleteFriendships(ctx context.Context, id1 int, id2 int) global.ResponseStatusCode
	FindFriendships(ctx context.Context, id int) (global.ResponseStatusCode, map[int]string)
	AddFriendshipsRequest(ctx context.Context, userId int, requestId int) global.ResponseStatusCode
	ChangeFriendshipsRequest(ctx context.Context, RequestId int, IsFriend bool, mailId int) global.ResponseStatusCode

	GiveCardByCardId(ctx context.Context, UserId int, CardId int) global.ResponseStatusCode        //给指定卡，没有保护，不可以用没有的cardId
	GiveInitCardBag(ctx context.Context, UserId int) global.ResponseStatusCode                     //给1级初始卡包
	GetBags(ctx context.Context, UserId int) ([]BattleData.BagStuffDto, global.ResponseStatusCode) //背包渲染

	GetUserGold(ctx context.Context, userId int) (int, global.ResponseStatusCode)
	SellCard(ctx context.Context, userId int, stuffIdList []int) global.ResponseStatusCode
	CheckBtDataIsValid(ctx context.Context, userId int, data []BattleData.BagStuffDto, gold int) global.ResponseStatusCode
}

type User_service_impl struct {
	repo userrepo.User_repo
}

func (u *User_service_impl) Set_repo(r repo.Repo) {
	u.repo = r.(userrepo.User_repo)
	u.Init_key()
}

func (u *User_service_impl) Find_user_by_name(ctx context.Context, name string) (*User_entity.User, global.ResponseStatusCode) {
	// 传入 ctx 和 u.repo.Get_db()
	e, err := u.repo.Get_by_name(ctx, u.repo.Get_db(), name)
	if err != global.ResponseSuccess {
		return &User_entity.User{}, err
	} else {
		return e, global.ResponseSuccess
	}
}

func (u *User_service_impl) Find_user_by_id(ctx context.Context, id int) (*User_entity.User, global.ResponseStatusCode) {
	// 传入 ctx 和 u.repo.Get_db()
	e, err := u.repo.Get_by_id(ctx, u.repo.Get_db(), id)
	if err != global.ResponseSuccess {
		return &User_entity.User{}, err
	} else {
		return e, global.ResponseSuccess
	}
}

func (u *User_service_impl) hash_password(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes)
}

func (u *User_service_impl) Check_password(ctx context.Context, id int, password string) global.ResponseStatusCode {
	// 调用内部方法需传递 ctx
	e, err := u.Find_user_by_id(ctx, id)
	if err != global.ResponseSuccess {
		return err
	}
	hash := e.Password
	ok := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if ok != nil {
		return global.ResponseIncorrectPassword
	}
	return global.ResponseSuccess
}

func (u *User_service_impl) Create_user(ctx context.Context, user *User_entity.User) global.ResponseStatusCode {
	tx, errDb := u.repo.Get_db().BeginTx(ctx, nil)
	if errDb != nil {
		return global.ResponseInternalServersError
	}
	defer tx.Rollback()

	user.Password = u.hash_password(user.Password)
	err := u.repo.Create(ctx, tx, user)
	if err != global.ResponseSuccess {

		return err
	}

	temp, _ := u.repo.Get_by_name(ctx, tx, user.Name)

	UserId := temp.Id

	err = u.repo.CreateAsset(ctx, tx, UserId)
	if err != global.ResponseSuccess {

		return err
	}
	err = u.repo.UpdateAssetGold(ctx, tx, UserId, 500)
	if err != global.ResponseSuccess {

		return err
	}
	tx.Commit()
	return global.ResponseSuccess
}

func (u *User_service_impl) Delete_user(ctx context.Context, id int) global.ResponseStatusCode {
	e := User_entity.User{
		Id:   id,
		Name: fmt.Sprintf("已注销用户_%d", id),
	}
	// 传入 ctx 和 u.repo.Get_db()
	err := u.repo.ChangeUserNameByID(ctx, u.repo.Get_db(), e.Id, e.Name)
	if err != global.ResponseSuccess {
		return err
	}
	// 传入 ctx 和 u.repo.Get_db()
	err = u.repo.DestroyPassword(ctx, u.repo.Get_db(), e.Id)
	if err != global.ResponseSuccess {
		return err
	}
	u.repo.UpdateActiveInRedisByUserId(e.Id, ctx)
	return global.ResponseSuccess
}

func (u *User_service_impl) Update_user(ctx context.Context, e *User_entity.User) global.ResponseStatusCode {
	if !(e.Password == "") {
		e.Password = u.hash_password(e.Password)
	}
	err := u.repo.Update(ctx, u.repo.Get_db(), e)
	return err
}

func (u *User_service_impl) CreateFriendships(ctx context.Context, id1 int, id2 int) global.ResponseStatusCode {
	return u.repo.SaveFriendships(ctx, u.repo.Get_db(), id1, id2)
}
func (u *User_service_impl) DeleteFriendships(ctx context.Context, id1 int, id2 int) global.ResponseStatusCode {
	return u.repo.DeleteFriendships(ctx, u.repo.Get_db(), id1, id2)
}
func (u *User_service_impl) FindFriendships(ctx context.Context, id int) (global.ResponseStatusCode, map[int]string) {
	err, idList := u.repo.FindFriendships(ctx, u.repo.Get_db(), id)
	res := make(map[int]string)
	if err != global.ResponseSuccess {
		return err, nil
	}
	for _, id := range idList {
		e, err := u.repo.Get_by_id(ctx, u.repo.Get_db(), id)
		if err != global.ResponseSuccess {
			return err, nil
		}
		res[id] = e.Name
	}
	return global.ResponseSuccess, res
}
func (u *User_service_impl) AddFriendshipsRequest(ctx context.Context, userId int, requestId int) global.ResponseStatusCode {
	tx, errDb := u.repo.Get_db().BeginTx(ctx, nil)
	if errDb != nil {
		return global.ResponseInternalServersError

	}
	defer tx.Rollback()

	e, err := u.repo.Get_by_id(ctx, u.repo.Get_db(), userId)
	if err != global.ResponseSuccess {
		return err
	}
	UserName := e.Name
	m, err := mail.NewMail(requestId, userId, fmt.Sprintf("%s向你发起好友请求", UserName), "FriendshipsRequest")
	if err != global.ResponseSuccess {
		return err
	}
	err = u.repo.SaveMail(ctx, tx, m)
	if err != global.ResponseSuccess {
		return err
	}
	err = u.repo.SaveFriendships(ctx, u.repo.Get_db(), userId, requestId)
	if err != global.ResponseSuccess {
		return err
	}
	tx.Commit()
	return global.ResponseSuccess
}
func (u *User_service_impl) ChangeFriendshipsRequest(ctx context.Context, RequestId int, IsFriend bool, mailId int) global.ResponseStatusCode {
	var SendId int
	mailList, err := u.repo.FindMails(ctx, u.repo.Get_db(), mail.Filter{Id: fmt.Sprintf("%d", mailId)}, 1)
	if err != global.ResponseSuccess {
		return err
	}
	m := mailList[0]
	if m.Category != "FriendshipsRequest" {
		return global.ResponseInvalidReqParams
	}
	if m.AcceptId != RequestId {
		return global.ResponseForbidden
	}
	SendId = m.SendId
	if IsFriend {
		return u.repo.ChangeFriendships(ctx, u.repo.Get_db(), RequestId, SendId, IsFriend)
	} else {
		return u.repo.DeleteFriendships(ctx, u.repo.Get_db(), RequestId, SendId)
	}
}

func (u *User_service_impl) GiveCardByCardId(ctx context.Context, UserId int, CardId int) global.ResponseStatusCode {
	return u.repo.AddCardInBags(ctx, u.repo.Get_db(), CardId, UserId)
}

func (u *User_service_impl) GiveInitCardBag(ctx context.Context, UserId int) global.ResponseStatusCode {
	CardIdList := make([]int, 0, 5)
	CardIdList = append(CardIdList, 0+Util.RandomRange(0, global.Lev1Category1Num-1))
	CardIdList = append(CardIdList, 0+Util.RandomRange(0, global.Lev1Category1Num-1))    //两张母战
	CardIdList = append(CardIdList, 1000+Util.RandomRange(0, global.Lev1Category2Num-1)) //一张母法
	CardIdList = append(CardIdList, 2000+Util.RandomRange(0, global.Lev1Category3Num-1)) //一张子战
	CardIdList = append(CardIdList, 3000+Util.RandomRange(0, global.Lev1Category4Num-1)) //一张子法
	fmt.Println(CardIdList)
	tx, errDb := u.repo.Get_db().BeginTx(ctx, nil)
	if errDb != nil {
		return global.ResponseInternalServersError

	}
	defer tx.Rollback()
	for _, cardId := range CardIdList {
		err := u.repo.AddCardInBags(ctx, tx, cardId, UserId)
		if err != global.ResponseSuccess {
			return err
		}
	}
	tx.Commit()
	return global.ResponseSuccess

}

func (u *User_service_impl) GetBags(ctx context.Context, UserId int) ([]BattleData.BagStuffDto, global.ResponseStatusCode) {
	res, err := u.repo.GetBagsByUserId(ctx, u.repo.Get_db(), UserId)
	return res, err
}

func (u *User_service_impl) GetUserGold(ctx context.Context, userId int) (int, global.ResponseStatusCode) {
	err, res := u.repo.GetAssetGold(ctx, u.repo.Get_db(), userId)
	return res, err
}

func (u *User_service_impl) SellCard(ctx context.Context, userId int, stuffIdList []int) global.ResponseStatusCode {
	tx, errDb := u.repo.Get_db().BeginTx(ctx, nil)
	if errDb != nil {
		return global.ResponseInternalServersError

	}
	defer tx.Rollback()
	for _, stuffId := range stuffIdList {
		err, temp := u.repo.GetStuffByStuffId(ctx, tx, userId, stuffId)
		if err != global.ResponseSuccess {
			return err
		}
		price := temp.Price
		err = u.repo.UpdateAssetGold(ctx, tx, userId, price)
		if err != global.ResponseSuccess {
			return err
		}
		err = u.repo.DeleteStuff(ctx, tx, userId, stuffId)
		if err != global.ResponseSuccess {
			return err
		}
	}

	tx.Commit()
	return global.ResponseSuccess
}

func (u *User_service_impl) CheckBtDataIsValid(ctx context.Context, userId int, data []BattleData.BagStuffDto, gold int) global.ResponseStatusCode {
	if len(data) != 5 {
		return global.BattleCardNumErr
	}
	for _, e := range data {
		err, _ := u.repo.GetStuffByStuffId(ctx, u.repo.Get_db(), userId, e.StuffId)
		if err != global.ResponseSuccess {
			return err
		}
	}

	child_num := 0
	character_num := 0
	for _, e := range data {
		err1, flag1 := u.repo.JudgeCardIsParent(ctx, u.repo.Get_db(), e.CardId)
		if err1 != global.ResponseSuccess {
			return err1
		}
		if !flag1 {
			child_num++
		}
		err2, flag2 := u.repo.JudgeCardIsCharacter(ctx, u.repo.Get_db(), e.CardId)
		if err2 != global.ResponseSuccess {
			return err2
		}
		if flag2 {
			character_num++
		}

	}
	//至少2子
	if child_num > 2 {
		return global.BattleEnterDataInvalid
	}
	//三战二法
	if character_num != 3 {
		return global.BattleEnterDataInvalid
	}

	//判定钱够不够
	err3, gold_num := u.repo.GetAssetGold(ctx, u.repo.Get_db(), userId)
	if err3 != global.ResponseSuccess {
		return err3
	}
	if gold_num-gold < 0 {
		return global.ResponseGoldNotEnough
	}

	return global.ResponseSuccess
}
