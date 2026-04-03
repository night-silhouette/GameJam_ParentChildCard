package UserService

import (
	"fmt"
	"pcc_card/application/entity/User_entity"
	"pcc_card/application/entity/mail"
	"pcc_card/global"
)

func (u *User_service_impl) GetAllOnePage(AcceptId int, page int) ([]*mail.Mail, global.ResponseStatusCode) {
	res, err := u.repo.FindMails(mail.Filter{AcceptId: fmt.Sprintf("%d", AcceptId)}, page)
	return res, err
}
func (u *User_service_impl) GetMailStatus(id int) (int, global.ResponseStatusCode) {
	res, err := u.repo.CheckMailUnReadNumByUserId(id)
	return res, err
}

func (u *User_service_impl) SendMail(SendId int, body string, AcceptId int) global.ResponseStatusCode {
	m := mail.NewMail(AcceptId, SendId, body, "UserMail")
	err := u.repo.SaveMail(m)
	return err
}
func (u *User_service_impl) ChangeMailStatus(AcceptId int, MailId int, status int) global.ResponseStatusCode {
	err := u.repo.UpdateMail(&mail.Filter{Id: fmt.Sprintf("%d", MailId), AcceptId: fmt.Sprintf("%d", AcceptId)}, &mail.Mail{Status: status})
	if err != global.ResponseSuccess {
		return err
	}
	return global.ResponseSuccess
}

func (u *User_service_impl) DeleteMailByMailId(MailId int, AcceptId int) global.ResponseStatusCode {
	err := u.repo.DeleteMail(&mail.Filter{Id: fmt.Sprintf("%d", MailId), AcceptId: fmt.Sprintf("%d", AcceptId)})
	if err != global.ResponseSuccess {
		return err
	}
	return global.ResponseSuccess
}
func (u *User_service_impl) DeleteMailAll(AcceptId int) global.ResponseStatusCode {
	err := u.repo.DeleteMail(&mail.Filter{AcceptId: fmt.Sprintf("%d", AcceptId)})
	if err != global.ResponseSuccess {
		return err
	}
	return global.ResponseSuccess
}

func (u *User_service_impl) UserSearch(NameVague string) (global.ResponseStatusCode, []*User_entity.User) {
	return u.repo.UserSearch(NameVague)
}
