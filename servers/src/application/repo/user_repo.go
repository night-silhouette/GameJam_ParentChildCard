package repo

import "pcc_card/application/entity"

type User_repo interface {
	Create(e entity.User) error
	Get_by_name(name string) (entity.User, error)
	Get_by_id(id int) (entity.User, error)
	Update(e entity.User) error
	Delete(e entity.User) error
}
