package userrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"pcc_card/application/entity/User_entity"
	"pcc_card/global"
	"pcc_card/infra/repo"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/redis/go-redis/v9"
)

//数据库表user字段名：
//id
//user_name
//hash_password
//create_at

type User_repo interface {
	repo.Repo
	Create(e *User_entity.User) global.ResponseStatusCode
	Get_by_name(name string) (*User_entity.User, global.ResponseStatusCode)
	Get_by_id(id int) (*User_entity.User, global.ResponseStatusCode)
	Update(e *User_entity.User) global.ResponseStatusCode
	Delete(e *User_entity.User) global.ResponseStatusCode
}

type User_repo_impl struct { //repo的实现
	db *sql.DB
	rd *redis.Client
}

func (r *User_repo_impl) Set_db(db *sql.DB, rd *redis.Client) {
	r.db = db
	r.rd = rd
}

func (r *User_repo_impl) Create(e *User_entity.User) global.ResponseStatusCode {
	query := "INSERT INTO users (user_name, hash_password,is_admin) VALUES ($1,$2,$3)"
	_, err := r.db.Exec(query, e.Name, e.Password, e.Is_admin)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				return global.ResponseDuplicateDataEntry
			case "23502":
				return global.ResponseRequiredParamsMissing
			}
		}
	}
	return global.ResponseSuccess
}

func (r *User_repo_impl) Get_by_name(name string) (*User_entity.User, global.ResponseStatusCode) {
	query := "SELECT id, user_name, hash_password, is_admin FROM users WHERE user_name = $1"
	row := r.db.QueryRow(query, name)
	e := User_entity.User{}
	err := row.Scan(&e.Id, &e.Name, &e.Password, &e.Is_admin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &User_entity.User{}, global.ResponseDataNotFound
		} else {
			fmt.Println("sql err", err)
			return &User_entity.User{}, global.ResponseInternalServersError
		}
	}
	return &e, global.ResponseSuccess
}

func (r *User_repo_impl) Get_by_id(id int) (*User_entity.User, global.ResponseStatusCode) {
	query := "SELECT id, user_name, hash_password,is_admin FROM users WHERE id = $1"
	row := r.db.QueryRow(query, id)
	e := User_entity.User{}
	err := row.Scan(&e.Id, &e.Name, &e.Password, &e.Is_admin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &User_entity.User{}, global.ResponseDataNotFound
		} else {
			fmt.Println("sql err", err)
			return &User_entity.User{}, global.ResponseInternalServersError
		}
	}
	return &e, global.ResponseSuccess

}

func (r *User_repo_impl) Update(e *User_entity.User) global.ResponseStatusCode {
	var res sql.Result
	var err error
	if e.Password != "" {
		query := "update users set user_name = $1, hash_password = $2 where id = $3"
		res, err = r.db.Exec(query, e.Name, e.Password, e.Id)
	} else {
		query := "update users set user_name = $1 where id = $2"
		res, err = r.db.Exec(query, e.Name, e.Id)
	}

	if err != nil {
		return global.ResponseInternalServersError
	}
	count, _ := res.RowsAffected()
	if count == 0 {
		return global.ResponseDataNotFound
	}
	return global.ResponseSuccess
}

func (r *User_repo_impl) Delete(e *User_entity.User) global.ResponseStatusCode {
	query := "delete from users where id = $1"
	res, err := r.db.Exec(query, e.Id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23503" {
				return global.ResponseDependentRecordsExist
			}
		}
	}
	count, _ := res.RowsAffected()
	if count == 0 {
		return global.ResponseDataNotFound
	}
	return global.ResponseSuccess
}

func (r *User_repo_impl) UpdateActiveInRedisByUserId(id int, ctx context.Context) int { //添加或者自增
	var res int
	param := fmt.Sprintf("user:active:%d", id)
	count, _ := r.rd.Exists(ctx, param).Result()
	if count == 0 {
		r.rd.Set(ctx, param, 0, 0)
		res = 0
	} else {
		value, _ := r.rd.Get(ctx, param).Int()
		value += 1
		r.rd.Set(ctx, param, value, 0)
		res = value
	}
	return res
}

func (r *User_repo_impl) CheckActiveInRedisByUserId(id int, ctx context.Context) int {
	param := fmt.Sprintf("user:active:%d", id)
	count, _ := r.rd.Exists(ctx, param).Result()
	if count == 0 {
		return -1
	} else {
		value, _ := r.rd.Get(ctx, param).Int()
		return value
	}
}
