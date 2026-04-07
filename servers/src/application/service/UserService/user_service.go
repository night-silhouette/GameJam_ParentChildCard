package UserService

import (
	"context"
	"fmt"
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
	user.Password = u.hash_password(user.Password)
	// 传入 ctx 和 u.repo.Get_db()
	return u.repo.Create(ctx, u.repo.Get_db(), user)
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
	// 传入 ctx 和 u.repo.Get_db()
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
	var res map[int]string
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
	mailList, err := u.repo.FindMails(ctx, u.repo.Get_db(), mail.Filter{Id: fmt.Sprintf("%s", mailId)}, 1)
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
	return u.repo.ChangeFriendships(ctx, u.repo.Get_db(), RequestId, SendId, IsFriend)
}
