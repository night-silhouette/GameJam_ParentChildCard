package userrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"pcc_card/application/entity/User_entity"
	"pcc_card/application/entity/mail"
	"pcc_card/global"
	"pcc_card/infra/repo"
	"strings"

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
	UpdateActiveInRedisByUserId(id int, ctx context.Context) int
	CheckActiveInRedisByUserId(id int, ctx context.Context) int
	ChangeUserNameByID(id int, name string) global.ResponseStatusCode
	DestroyPassword(id int) global.ResponseStatusCode
	UpdateMail(f *mail.Filter, data *mail.Mail) global.ResponseStatusCode
	SaveMail(m *mail.Mail) global.ResponseStatusCode
	DeleteMail(f *mail.Filter) global.ResponseStatusCode
	FindMails(f mail.Filter) ([]*mail.Mail, global.ResponseStatusCode)
	CheckMailUnReadNumByUserId(userId int) (int, global.ResponseStatusCode)
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

func (r *User_repo_impl) ChangeUserNameByID(id int, name string) global.ResponseStatusCode {
	query := "update users set user_name = $1 where id = $2"
	res, err := r.db.Exec(query, name, id)
	if err != nil {
		fmt.Println(err)
		return global.ResponseInternalServersError
	}
	count, _ := res.RowsAffected()
	if count == 0 {
		return global.ResponseDataNotFound
	}

	return global.ResponseSuccess
}

func (r *User_repo_impl) DestroyPassword(id int) global.ResponseStatusCode {
	query := "update users set hash_password = 'DISABLED_' || hash_password where id = $1"
	_, err := r.db.Exec(query, id)
	if err != nil {
		fmt.Println(err)
		return global.ResponseInternalServersError
	}
	return global.ResponseSuccess
}

//----------------------------------------------------mails----------------------------------------------------------

//type Mail struct {
//	AcceptId int    `json:"accept_id"`
//	SendId int    `json:"send_id"`
//	Body     string `json:"body"`
//	Category string `json:"category"`
//	Status   int    `json:"status"`
//}

func (r *User_repo_impl) SaveMail(m *mail.Mail) global.ResponseStatusCode {
	if m.SendId == 0 || m.AcceptId == 0 || m.Category == "" {
		return global.ResponseRequiredParamsMissing
	}
	query := "INSERT INTO mails(send_id, accept_id, body, category,status) values ($1, $2, $3, $4,$5)"
	_, err := r.db.Exec(query, m.SendId, m.AcceptId, m.Body, m.Category, m.Status)
	if err != nil {
		log.Fatalln(err.Error())
		return global.ResponseInternalServersError
	}
	return global.ResponseSuccess
}

func (r *User_repo_impl) DeleteMail(f *mail.Filter) global.ResponseStatusCode {
	query := "delete from mails where 1=1"
	var args []interface{}
	argCount := 1
	if f.Id != "" {
		query += fmt.Sprintf(" and id = $%d", argCount)
		args = append(args, f.Id)
		argCount++
	}
	//1. 处理 AcceptId
	if f.AcceptId != "" {
		query += fmt.Sprintf(" and accept_id = $%d", argCount)
		args = append(args, f.AcceptId)
		argCount++
	}

	// 2. 处理 SendId
	if f.SendId != "" {
		query += fmt.Sprintf(" and send_id = $%d", argCount)
		args = append(args, f.SendId)
		argCount++
	}

	// 3. 处理 Category
	if f.Category != "" {
		query += fmt.Sprintf(" and category = $%d", argCount)
		args = append(args, f.Category)
		argCount++
	}

	// 4. 处理 Status
	if f.Status != "" {
		query += fmt.Sprintf(" and status = $%d", argCount)
		args = append(args, f.Status)
		argCount++
	}

	// 5. 安全检查：如果没有任何过滤条件，拒绝执行
	if argCount == 1 {
		return global.ResponseRequiredParamsMissing
	}

	// 6. 执行删除
	result, err := r.db.Exec(query, args...)
	if err != nil {
		return global.ResponseInternalServersError
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return global.ResponseDataNotFound
	}

	return global.ResponseSuccess
}

func (r *User_repo_impl) FindMails(f mail.Filter) ([]*mail.Mail, global.ResponseStatusCode) {
	query := "select id,accept_id, send_id, body, category, status,created_at from mails where 1=1"
	var args []interface{}
	argCount := 1
	if f.Id != "" {
		query += fmt.Sprintf(" and id = $%d", argCount)
		args = append(args, f.Id)
		argCount++
	}

	if f.AcceptId != "" {
		query += fmt.Sprintf(" and accept_id = $%d", argCount)
		args = append(args, f.AcceptId)
		argCount++
	}
	if f.SendId != "" {
		query += fmt.Sprintf(" and send_id = $%d", argCount)
		args = append(args, f.SendId)
		argCount++
	}
	if f.Category != "" {
		query += fmt.Sprintf(" and category = $%d", argCount)
		args = append(args, f.Category)
		argCount++
	}
	if f.Status != "" {
		query += fmt.Sprintf(" and status = $%d", argCount)
		args = append(args, f.Status)
		argCount++
	}

	query += " order by created_at desc"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, global.ResponseInternalServersError
	}
	defer rows.Close()

	var mails []*mail.Mail
	for rows.Next() {
		var m mail.Mail
		if err := rows.Scan(&m.Id, &m.AcceptId, &m.SendId, &m.Body, &m.Category, &m.Status, &m.CreateAt); err != nil {
			return nil, global.ResponseInternalServersError
		}
		mails = append(mails, &m)
	}

	if len(mails) == 0 {
		return nil, global.ResponseDataNotFound
	}

	return mails, global.ResponseSuccess
}

func (r *User_repo_impl) UpdateMail(f *mail.Filter, data *mail.Mail) global.ResponseStatusCode {
	// 1. 基础语句
	query := "update mails set "
	var args []interface{}
	argCount := 1

	// --- 第一步：拼接 SET 部分 (要改什么) ---
	// 注意：这里我们只更新有意义的字段，比如 Body, Category, Status
	var setClauses []string

	// 如果 Body 不为空，说明要改正文
	if data.Body != "" {
		setClauses = append(setClauses, fmt.Sprintf("body = $%d", argCount))
		args = append(args, data.Body)
		argCount++
	}

	// 如果 Category 不为空，说明要改分类
	if data.Category != "" {
		setClauses = append(setClauses, fmt.Sprintf("category = $%d", argCount))
		args = append(args, data.Category)
		argCount++
	}

	// 更新状态 (比如从 0 未读变成 1 已读)
	// 这里假设 Status 是传进来的新状态
	setClauses = append(setClauses, fmt.Sprintf("status = $%d", argCount))
	args = append(args, data.Status)
	argCount++

	// 如果没有任何要改的内容，直接返回
	if len(setClauses) == 0 {
		return global.ResponseRequiredParamsMissing
	}
	query += strings.Join(setClauses, ", ")

	// --- 第二步：拼接 WHERE 部分 (你的 Filter 条件) ---
	query += " where 1=1"
	if f.Id != "" {
		query += fmt.Sprintf(" and id = $%d", argCount)
		args = append(args, f.Id)
		argCount++
	}

	if f.AcceptId != "" {
		query += fmt.Sprintf(" and accept_id = $%d", argCount)
		args = append(args, f.AcceptId)
		argCount++
	}

	if f.SendId != "" {
		query += fmt.Sprintf(" and send_id = $%d", argCount)
		args = append(args, f.SendId)
		argCount++
	}

	if f.Category != "" {
		query += fmt.Sprintf(" and category = $%d", argCount)
		args = append(args, f.Category)
		argCount++
	}

	if f.Status != "" {
		query += fmt.Sprintf(" and status = $%d", argCount)
		args = append(args, f.Status)
		argCount++
	}

	// 安全检查：Update 必须带条件，否则拒绝执行，防止全表更新
	if !strings.Contains(query, "and") {
		return global.ResponseRequiredParamsMissing
	}

	// 2. 执行更新
	result, err := r.db.Exec(query, args...)
	if err != nil {
		return global.ResponseInternalServersError
	}

	// 3. 检查是否真的更新到了数据
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return global.ResponseDataNotFound
	}

	return global.ResponseSuccess
}

func (r *User_repo_impl) CheckMailUnReadNumByUserId(userId int) (int, global.ResponseStatusCode) {
	query := "select unread_count from users where id=$1"
	var unreadCount int

	err := r.db.QueryRow(query, userId).Scan(&unreadCount)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, global.ResponseDataNotFound
		}
		log.Println(err)
		return 0, global.ResponseInternalServersError
	}
	return unreadCount, global.ResponseSuccess

}
