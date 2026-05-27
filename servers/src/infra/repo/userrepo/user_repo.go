package userrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"pcc_card/Util"
	"pcc_card/application/entity/BattleData"
	"pcc_card/application/entity/User_entity"
	"pcc_card/application/entity/mail"
	"pcc_card/global"
	"pcc_card/infra/repo"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/redis/go-redis/v9"
)

// User_repo 接口定义
type User_repo interface {
	repo.Repo
	Create(ctx context.Context, db repo.SQLQueryer, e *User_entity.User) global.ResponseStatusCode
	Get_by_name(ctx context.Context, db repo.SQLQueryer, name string) (*User_entity.User, global.ResponseStatusCode)
	Get_by_id(ctx context.Context, db repo.SQLQueryer, id int) (*User_entity.User, global.ResponseStatusCode)
	Update(ctx context.Context, db repo.SQLQueryer, e *User_entity.User) global.ResponseStatusCode
	Delete(ctx context.Context, db repo.SQLQueryer, e *User_entity.User) global.ResponseStatusCode
	UpdateActiveInRedisByUserId(id int, ctx context.Context) int
	CheckActiveInRedisByUserId(id int, ctx context.Context) int
	ChangeUserNameByID(ctx context.Context, db repo.SQLQueryer, id int, name string) global.ResponseStatusCode
	DestroyPassword(ctx context.Context, db repo.SQLQueryer, id int) global.ResponseStatusCode
	UpdateMail(ctx context.Context, db repo.SQLQueryer, f *mail.Filter, data *mail.Mail) global.ResponseStatusCode
	SaveMail(ctx context.Context, db repo.SQLQueryer, m *mail.Mail) global.ResponseStatusCode
	DeleteMail(ctx context.Context, db repo.SQLQueryer, f *mail.Filter) global.ResponseStatusCode
	FindMails(ctx context.Context, db repo.SQLQueryer, f mail.Filter, page int) ([]*mail.Mail, global.ResponseStatusCode)
	CheckMailUnReadNumByUserId(ctx context.Context, db repo.SQLQueryer, userId int) (int, global.ResponseStatusCode)
	UserSearch(ctx context.Context, db repo.SQLQueryer, NameVague string) (global.ResponseStatusCode, []*User_entity.User)
	SaveFriendships(ctx context.Context, db repo.SQLQueryer, id1 int, id2 int) global.ResponseStatusCode
	FindFriendships(ctx context.Context, db repo.SQLQueryer, userId int) (global.ResponseStatusCode, []int)
	DeleteFriendships(ctx context.Context, db repo.SQLQueryer, id1 int, id2 int) global.ResponseStatusCode
	ChangeFriendships(ctx context.Context, db repo.SQLQueryer, id1 int, id2 int, request bool) global.ResponseStatusCode

	AddCardInBags(ctx context.Context, db repo.SQLQueryer, cardID int, userID int) global.ResponseStatusCode
	GetBagsByUserId(ctx context.Context, db repo.SQLQueryer, userID int) ([]BattleData.BagStuffDto, global.ResponseStatusCode)
	CreateAsset(ctx context.Context, db repo.SQLQueryer, userId int) global.ResponseStatusCode
	UpdateAssetGold(ctx context.Context, db repo.SQLQueryer, userId int, gold int) global.ResponseStatusCode
	GetAssetGold(ctx context.Context, db repo.SQLQueryer, userId int) (global.ResponseStatusCode, int)
	DeleteStuff(ctx context.Context, db repo.SQLQueryer, userId int, stuffId int) global.ResponseStatusCode
	GetStuffByStuffId(ctx context.Context, db repo.SQLQueryer, userId int, stuffId int) (global.ResponseStatusCode, BattleData.BagStuffDto)
}

type User_repo_impl struct {
	db *sql.DB
	rd *redis.Client
}

func (r *User_repo_impl) Get_db() *sql.DB {
	return r.db
}

func (r *User_repo_impl) Set_db(db *sql.DB, rd *redis.Client) {
	r.db = db
	r.rd = rd
}

// ---------------------------------------------------- User ----------------------------------------------------------

func (r *User_repo_impl) Create(ctx context.Context, db repo.SQLQueryer, e *User_entity.User) global.ResponseStatusCode {
	query := "insert into users (user_name, hash_password, is_admin) values ($1, $2, $3)"
	_, err := db.ExecContext(ctx, query, e.Name, e.Password, e.Is_admin)
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

func (r *User_repo_impl) Get_by_name(ctx context.Context, db repo.SQLQueryer, name string) (*User_entity.User, global.ResponseStatusCode) {
	query := "select id, user_name, hash_password, is_admin from users where user_name = $1"
	row := db.QueryRowContext(ctx, query, name)
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

func (r *User_repo_impl) Get_by_id(ctx context.Context, db repo.SQLQueryer, id int) (*User_entity.User, global.ResponseStatusCode) {

	query := "select id, user_name, hash_password, is_admin from users where id = $1"
	row := db.QueryRowContext(ctx, query, id)
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

func (r *User_repo_impl) Update(ctx context.Context, db repo.SQLQueryer, e *User_entity.User) global.ResponseStatusCode {
	var res sql.Result
	var err error
	if e.Password != "" {
		query := "update users set user_name = $1, hash_password = $2 where id = $3"
		res, err = db.ExecContext(ctx, query, e.Name, e.Password, e.Id)
	} else {
		query := "update users set user_name = $1 where id = $2"
		res, err = db.ExecContext(ctx, query, e.Name, e.Id)
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

func (r *User_repo_impl) Delete(ctx context.Context, db repo.SQLQueryer, e *User_entity.User) global.ResponseStatusCode {
	query := "delete from users where id = $1"
	res, err := db.ExecContext(ctx, query, e.Id)
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

func (r *User_repo_impl) ChangeUserNameByID(ctx context.Context, db repo.SQLQueryer, id int, name string) global.ResponseStatusCode {
	query := "update users set user_name = $1 where id = $2"
	res, err := db.ExecContext(ctx, query, name, id)
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

func (r *User_repo_impl) DestroyPassword(ctx context.Context, db repo.SQLQueryer, id int) global.ResponseStatusCode {
	query := "update users set hash_password = 'DISABLED_' || hash_password where id = $1"
	_, err := db.ExecContext(ctx, query, id)
	if err != nil {
		fmt.Println(err)
		return global.ResponseInternalServersError
	}
	return global.ResponseSuccess
}

// ---------------------------------------------------- Redis ----------------------------------------------------------

func (r *User_repo_impl) UpdateActiveInRedisByUserId(id int, ctx context.Context) int {
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

// ---------------------------------------------------- Mails ----------------------------------------------------------

func (r *User_repo_impl) SaveMail(ctx context.Context, db repo.SQLQueryer, m *mail.Mail) global.ResponseStatusCode {
	if m.SendId == 0 || m.AcceptId == 0 || m.Category == "" {

		return global.ResponseRequiredParamsMissing
	}
	query := "insert into mails (send_id, accept_id, body, category, status) values ($1, $2, $3, $4, $5)"
	_, err := db.ExecContext(ctx, query, m.SendId, m.AcceptId, m.Body, m.Category, m.Status)
	if err != nil {
		log.Println(err.Error())
		return global.ResponseInternalServersError
	}
	return global.ResponseSuccess
}

func (r *User_repo_impl) DeleteMail(ctx context.Context, db repo.SQLQueryer, f *mail.Filter) global.ResponseStatusCode {
	query := "delete from mails where 1=1"
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

	if argCount == 1 {
		return global.ResponseRequiredParamsMissing
	}

	result, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return global.ResponseInternalServersError
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return global.ResponseDataNotFound
	}
	return global.ResponseSuccess
}

func (r *User_repo_impl) FindMails(ctx context.Context, db repo.SQLQueryer, f mail.Filter, page int) ([]*mail.Mail, global.ResponseStatusCode) {
	if page <= 0 {
		page = 1
	}
	query := "select id, accept_id, send_id, body, category, status, created_at from mails where 1=1"
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
	offset := (page - 1) * 8
	query += " order by created_at desc"
	query += fmt.Sprintf(" limit $%d", argCount)
	args = append(args, 8)
	argCount++
	query += fmt.Sprintf(" offset $%d", argCount)
	args = append(args, offset)
	argCount++

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Println(err.Error())
		return nil, global.ResponseInternalServersError
	}
	defer rows.Close()

	var mails []*mail.Mail
	for rows.Next() {
		var m mail.Mail
		if err := rows.Scan(&m.MailId, &m.AcceptId, &m.SendId, &m.Body, &m.Category, &m.Status, &m.CreateAt); err != nil {
			log.Println(err.Error())
			return nil, global.ResponseInternalServersError
		}
		mails = append(mails, &m)
	}

	if len(mails) == 0 {
		return nil, global.ResponseDataNotFound
	}
	return mails, global.ResponseSuccess
}

func (r *User_repo_impl) UpdateMail(ctx context.Context, db repo.SQLQueryer, f *mail.Filter, data *mail.Mail) global.ResponseStatusCode {
	query := "update mails set "
	var args []interface{}
	argCount := 1
	var setClauses []string

	if data.Body != "" {
		setClauses = append(setClauses, fmt.Sprintf("body = $%d", argCount))
		args = append(args, data.Body)
		argCount++
	}
	if data.Category != "" {
		setClauses = append(setClauses, fmt.Sprintf("category = $%d", argCount))
		args = append(args, data.Category)
		argCount++
	}
	setClauses = append(setClauses, fmt.Sprintf("status = $%d", argCount))
	args = append(args, data.Status)
	argCount++

	if len(setClauses) == 0 {
		return global.ResponseRequiredParamsMissing
	}
	query += strings.Join(setClauses, ", ")
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

	if !strings.Contains(query, "and") {
		return global.ResponseRequiredParamsMissing
	}

	result, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return global.ResponseInternalServersError
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return global.ResponseDataNotFound
	}
	return global.ResponseSuccess
}

func (r *User_repo_impl) CheckMailUnReadNumByUserId(ctx context.Context, db repo.SQLQueryer, userId int) (int, global.ResponseStatusCode) {
	query := "select unread_count from users where id = $1"
	var unreadCount int
	err := db.QueryRowContext(ctx, query, userId).Scan(&unreadCount)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, global.ResponseDataNotFound
		}
		log.Println(err)
		return 0, global.ResponseInternalServersError
	}
	return unreadCount, global.ResponseSuccess
}

func (r *User_repo_impl) UserSearch(ctx context.Context, db repo.SQLQueryer, NameVague string) (global.ResponseStatusCode, []*User_entity.User) {
	query := `
        select id, user_name 
        from users 
        where user_name ilike $1 
        order by similarity(user_name, $2) desc 
        limit 8`
	pattern := "%" + NameVague + "%"
	var uList []*User_entity.User
	rows, err := db.QueryContext(ctx, query, pattern, NameVague)

	if err != nil {
		log.Println(err)
		return global.ResponseDataNotFound, nil
	}
	defer rows.Close()
	for rows.Next() {
		var u User_entity.User
		if err := rows.Scan(&u.Id, &u.Name); err != nil {
			log.Println(err)
			return global.ResponseInternalServersError, nil
		}
		uList = append(uList, &u)
	}
	if len(uList) == 0 {
		return global.ResponseDataNotFound, nil
	}
	return global.ResponseSuccess, uList
}

func (r *User_repo_impl) ChangeFriendships(ctx context.Context, db repo.SQLQueryer, id1 int, id2 int, request bool) global.ResponseStatusCode {
	if id1 > id2 {
		id1, id2 = id2, id1
	}
	query := "update friendships set request=$1 where user_id_1=$2 and user_id_2=$3"
	_, err := db.ExecContext(ctx, query, request, id1, id2)
	if err != nil {
		return global.ResponseInternalServersError
	}
	return global.ResponseSuccess
}

func (r *User_repo_impl) SaveFriendships(ctx context.Context, db repo.SQLQueryer, id1 int, id2 int) global.ResponseStatusCode {
	if id1 > id2 {
		id1, id2 = id2, id1
	}
	query := "INSERT INTO friendships (user_id_1, user_id_2,request) VALUES ($1, $2,false)"
	_, err := db.ExecContext(ctx, query, id1, id2)
	if err != nil {
		return global.ResponseInternalServersError
	}
	return global.ResponseSuccess
}
func (r *User_repo_impl) DeleteFriendships(ctx context.Context, db repo.SQLQueryer, id1 int, id2 int) global.ResponseStatusCode {
	if id1 > id2 {
		id1, id2 = id2, id1
	}
	query := "DELETE FROM friendships WHERE user_id_1 = $1 AND user_id_2 = $2"
	_, err := db.ExecContext(ctx, query, id1, id2)
	if err != nil {
		return global.ResponseInternalServersError
	}
	return global.ResponseSuccess
}

func (r *User_repo_impl) FindFriendships(ctx context.Context, db repo.SQLQueryer, userId int) (global.ResponseStatusCode, []int) {
	query := "select user_id_1, user_id_2 from friendships where $1 in (user_id_1, user_id_2) and request=true;"
	rows, err := db.QueryContext(ctx, query, userId)
	if err != nil {
		return global.ResponseInternalServersError, nil
	}
	defer rows.Close()

	var res []int
	for rows.Next() {
		var id1, id2 int
		if err := rows.Scan(&id1, &id2); err != nil {
			return global.ResponseInternalServersError, nil
		}
		if id1 == userId {
			res = append(res, id2)
		} else {
			res = append(res, id1)
		}
	}
	if err = rows.Err(); err != nil {
		return global.ResponseInternalServersError, nil
	}

	return global.ResponseSuccess, res
}

func (r *User_repo_impl) AddCardInBags(ctx context.Context, db repo.SQLQueryer, cardID int, userID int) global.ResponseStatusCode {
	// 1. 获取基准价格
	var basePrice int
	queryPrice := `select price from newcards where id = $1`
	err := db.QueryRowContext(ctx, queryPrice, cardID).Scan(&basePrice)
	if err != nil {
		log.Println("获取原始价格失败:", err)
		return global.ResponseDataNotFound
	}
	// 2. 计算 10% 随机浮动 (整数)
	floatRange := int(float64(basePrice) * 0.1)
	finalPrice := basePrice
	if floatRange > 0 {
		// 使用你的 RandomRange 生成偏移
		offset := Util.RandomRange(-floatRange, floatRange)
		finalPrice = basePrice + offset
	}

	query := `insert into bags (user_id, card_id, price) values ($1, $2, $3)`
	_, err = db.ExecContext(ctx, query, userID, cardID, finalPrice)
	if err != nil {
		return global.ResponseBagsUnknownError
	}
	return global.ResponseSuccess
}

func (r *User_repo_impl) GetBagsByUserId(ctx context.Context, db repo.SQLQueryer, userID int) ([]BattleData.BagStuffDto, global.ResponseStatusCode) {
	// 1. 初始化切片，防止返回 nil
	res := make([]BattleData.BagStuffDto, 0)

	// 2. 编写 SQL 语句
	// 直接从 bags 表中按 user_id 查询所有属于该玩家的卡片实例
	query := `select stuff_id, card_id, price from bags where user_id = $1`

	// 3. 执行查询
	rows, err := db.QueryContext(ctx, query, userID)
	if err != nil {
		log.Printf("查询玩家背包失败: %v", err)
		return nil, global.ResponseBagsUnknownError
	}
	defer rows.Close()

	// 4. 遍历结果集
	for rows.Next() {
		var dto BattleData.BagStuffDto
		// 按照 select 的顺序 Scan：stuff_id -> card_id -> price
		err := rows.Scan(&dto.StuffId, &dto.CardId, &dto.Price)
		if err != nil {
			log.Printf("读取背包数据行失败: %v", err)
			continue
		}
		res = append(res, dto)
	}

	// 5. 检查遍历过程中的错误
	if err = rows.Err(); err != nil {
		log.Printf("遍历背包数据过程中出错: %v", err)
		return nil, global.ResponseBagsUnknownError
	}

	return res, global.ResponseSuccess
}

func (r *User_repo_impl) CreateAsset(ctx context.Context, db repo.SQLQueryer, userId int) global.ResponseStatusCode {
	// 使用小写 SQL 语句，符合你的习惯
	// gold 初始化为 0，你也可以根据需要设置初始金币（比如新手送 1000）
	query := `
		insert into assets (user_id, gold) 
		values ($1, 0) 
		on conflict (user_id) do nothing
	`

	_, err := db.ExecContext(ctx, query, userId)
	if err != nil {
		fmt.Println(err)
		return global.ResponseInternalServersError
	}

	return global.ResponseSuccess
}

// UpdateAssetGold 传入的是增减值
func (r *User_repo_impl) UpdateAssetGold(ctx context.Context, db repo.SQLQueryer, userId int, gold int) global.ResponseStatusCode {
	// 使用小写 SQL
	// 注意：gold 参数可以是正数（增加金币），也可以是负数（消耗金币）

	// 如果是消耗金币（gold < 0），我们需要额外判断余额是否足够
	var query string
	var err error
	var result sql.Result

	if gold >= 0 {
		// 增加金币：直接加
		query = `update assets set gold = gold + $1 where user_id = $2`
		result, err = db.ExecContext(ctx, query, gold, userId)
	} else {
		// 消耗金币：必须确保扣除后余额 >= 0
		// abs_gold 为扣除金额的绝对值
		absGold := -gold
		query = `update assets set gold = gold - $1 where user_id = $2 and gold >= $1`
		result, err = db.ExecContext(ctx, query, absGold, userId)
	}

	if err != nil {
		// 数据库执行出错（连接断开、语法错误等）
		return global.ResponseInternalServersError
	}

	// 检查是否有行被更新
	rows, _ := result.RowsAffected()
	if rows == 0 {
		// 如果 gold < 0 且 rows == 0，通常意味着“金币不足”
		// 如果 gold >= 0 且 rows == 0，意味着“用户不存在”
		if gold < 0 {
			return global.ResponseGoldNotEnough
		}
		return global.ResponseDataNotFound
	}

	return global.ResponseSuccess
}

func (r *User_repo_impl) GetAssetGold(ctx context.Context, db repo.SQLQueryer, userId int) (global.ResponseStatusCode, int) {
	// 使用小写 SQL
	query := `select gold from assets where user_id = $1`

	var gold int64
	// 执行查询并扫描结果到 gold 变量
	err := db.QueryRowContext(ctx, query, userId).Scan(&gold)

	if err != nil {
		// 1. 如果没找到记录
		if err == sql.ErrNoRows {
			return global.ResponseInternalServersError, 0
		}
		return global.ResponseInternalServersError, 0
	}

	// 成功获取，将 int64 转为 int 返回
	return global.ResponseSuccess, int(gold)
}

func (r *User_repo_impl) DeleteStuff(ctx context.Context, db repo.SQLQueryer, userId int, stuffId int) global.ResponseStatusCode {

	query := `delete from bags where user_id = $1 and stuff_id = $2`

	result, err := db.ExecContext(ctx, query, userId, stuffId)
	if err != nil {
		return global.ResponseInternalServersError
	}

	// 检查是否有行被删除
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return global.ResponseDataNotFound
	}

	return global.ResponseSuccess
}

func (r *User_repo_impl) GetStuffByStuffId(ctx context.Context, db repo.SQLQueryer, userId int, stuffId int) (global.ResponseStatusCode, BattleData.BagStuffDto) {
	// 定义接收数据的结构体
	var dto BattleData.BagStuffDto

	// 使用小写 SQL 语句
	// 务必带上 user_id 条件，确保玩家只能查到自己的东西（防止越权）
	query := `select stuff_id, card_id, price from bags where user_id = $1 and stuff_id = $2`

	// 执行查询并将结果扫描进 dto
	err := db.QueryRowContext(ctx, query, userId, stuffId).Scan(
		&dto.StuffId,
		&dto.CardId,
		&dto.Price,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			// 没找到该物品
			return global.ResponseDataNotFound, dto
		}
		// 数据库其他错误
		return global.ResponseInternalServersError, dto
	}

	return global.ResponseSuccess, dto
}
