package UserService

import (
	"context"
	"fmt"
	"pcc_card/application/entity/User_entity"
	"pcc_card/application/entity/mail"
	"pcc_card/global"
)

func (u *User_service_impl) GetAllOnePage(ctx context.Context, AcceptId int, page int) ([]*mail.Mail, global.ResponseStatusCode) {
	// 传入 ctx 和 u.repo.Get_db()
	res, err := u.repo.FindMails(ctx, u.repo.Get_db(), mail.Filter{AcceptId: fmt.Sprintf("%d", AcceptId)}, page)
	return res, err
}

func (u *User_service_impl) GetMailStatus(ctx context.Context, id int) (int, global.ResponseStatusCode) {
	res, err := u.repo.CheckMailUnReadNumByUserId(ctx, u.repo.Get_db(), id)
	return res, err
}

func (u *User_service_impl) SendMail(ctx context.Context, SendId int, body string, AcceptId int, Category string) global.ResponseStatusCode {
	m, err := mail.NewMail(AcceptId, SendId, body, Category)
	if err != global.ResponseSuccess {
		return err
	}
	err = u.repo.SaveMail(ctx, u.repo.Get_db(), m)
	return err
}

func (u *User_service_impl) ChangeMailStatus(ctx context.Context, AcceptId int, MailId []int, status int) global.ResponseStatusCode {
	tx, errDb := u.repo.Get_db().BeginTx(ctx, nil)
	if errDb != nil {
		return global.ResponseInternalServersError
	}
	defer tx.Rollback()

	for _, id := range MailId {
		err := u.repo.UpdateMail(ctx, tx, &mail.Filter{Id: fmt.Sprintf("%d", id), AcceptId: fmt.Sprintf("%d", AcceptId)}, &mail.Mail{Status: status})
		if err != global.ResponseSuccess {
			return err
		}
	}
	tx.Commit()
	return global.ResponseSuccess
}

func (u *User_service_impl) DeleteMailByMailId(ctx context.Context, MailId []int, AcceptId int) global.ResponseStatusCode {
	tx, err := u.repo.Get_db().BeginTx(ctx, nil)
	if err != nil {
		return global.ResponseInternalServersError
	}
	defer tx.Rollback()
	for _, id := range MailId {
		f := &mail.Filter{Id: fmt.Sprintf("%d", id), AcceptId: fmt.Sprintf("%d", AcceptId)}
		err2 := u.repo.DeleteMail(ctx, tx, f)
		if err2 != global.ResponseSuccess {
			return err2
		}
	}
	err = tx.Commit()
	if err != nil {
		return global.ResponseInternalServersError
	}
	return global.ResponseSuccess
}

func (u *User_service_impl) DeleteMailAll(ctx context.Context, AcceptId int) global.ResponseStatusCode {

	err := u.repo.DeleteMail(ctx, u.repo.Get_db(), &mail.Filter{AcceptId: fmt.Sprintf("%d", AcceptId)})
	if err != global.ResponseSuccess {
		return err
	}
	return global.ResponseSuccess
}

func (u *User_service_impl) UserSearch(ctx context.Context, NameVague string) (global.ResponseStatusCode, []*User_entity.User) {
	return u.repo.UserSearch(ctx, u.repo.Get_db(), NameVague)
}
