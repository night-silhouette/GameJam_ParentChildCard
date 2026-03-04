package repo

import (
	"database/sql"
	"pcc_card/application/entity"
)

//数据库表user字段名：
//id
//user_name
//hash_password
//create_at

type User_repo interface {
	Repo
	Create(e *entity.User) error
	Get_by_name(name string) (*entity.User, error)
	Get_by_id(id int) (*entity.User, error)
	Update(e *entity.User) error
	Delete(e *entity.User) error
}

type User_repo_impl struct { //repo的实现
	db *sql.DB
}

func (r *User_repo_impl) Set_db(db *sql.DB) {
	r.db = db
}

func (r *User_repo_impl) Create(e *entity.User) error {
	query := "INSERT INTO user(user_name, hash_password) VALUES ($1,$2)"
	_, err := r.db.Exec(query, e.Name, e.Password)
	return err
}

func (r *User_repo_impl) Get_by_name(name string) (*entity.User, error) {
	query := "SELECT id, user_name, hash_password FROM user WHERE user_name = $1"
	row := r.db.QueryRow(query, name)
	e := entity.User{}
	err := row.Scan(&e.Id, &e.Name, &e.Password)
	return &e, err
}

func (r *User_repo_impl) Get_by_id(id int) (*entity.User, error) {
	query := "SELECT id, user_name, hash_password FROM user WHERE user_name = $1"
	row := r.db.QueryRow(query, id)
	e := entity.User{}
	err := row.Scan(&e.Id, &e.Name, &e.Password)
	return &e, err

}

func (r *User_repo_impl) Update(e *entity.User) error {
	query := "update user set user_name = $1, hash_password = $2 where id = $3"
	_, err := r.db.Exec(query, e.Name, e.Password, e.Id)
	return err
}

func (r *User_repo_impl) Delete(e *entity.User) error {
	query := "delete from user where id = $1"
	_, err := r.db.Exec(query, e.Id)
	return err
}
