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
	"strconv"
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
	JudgeCardIsParent(ctx context.Context, db repo.SQLQueryer, CardId int) (global.ResponseStatusCode, bool)
	JudgeCardIsCharacter(ctx context.Context, db repo.SQLQueryer, CardId int) (global.ResponseStatusCode, bool)
	CreateBattle(ctx context.Context, db repo.SQLQueryer, playerIdA int, playerIdB int) (int, global.ResponseStatusCode)
	CheckUserIdIsBattle(ctx context.Context, db repo.SQLQueryer, userId int) (int, global.ResponseStatusCode)
	DeleteBattle(ctx context.Context, db repo.SQLQueryer, BtId int) global.ResponseStatusCode
	CreateLoot(ctx context.Context, db repo.SQLQueryer, loot []int, UserId int) global.ResponseStatusCode
	GetLoot(ctx context.Context, db repo.SQLQueryer, UserId int) (global.ResponseStatusCode, []BattleData.LootDto)
	DeleteLoot(ctx context.Context, db repo.SQLQueryer, LootId int) global.ResponseStatusCode
	GetGoodsByUserId(ctx context.Context, db repo.SQLQueryer, UserId int) (global.ResponseStatusCode, []BattleData.GoodsDto)
	CreateGoods(ctx context.Context, db repo.SQLQueryer, UserId int, GoodsList []*BattleData.GoodsDto) global.ResponseStatusCode
	DeleteGoodsByUserId(ctx context.Context, db repo.SQLQueryer, UserId int) global.ResponseStatusCode
	GetCardPrice(ctx context.Context, db repo.SQLQueryer, cardID int) (global.ResponseStatusCode, int)
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

// 10% 随机浮动
func (r *User_repo_impl) GetCardPrice(ctx context.Context, db repo.SQLQueryer, cardID int) (global.ResponseStatusCode, int) {
	var basePrice int
	queryPrice := `select price from newcards where id = $1`
	err := db.QueryRowContext(ctx, queryPrice, cardID).Scan(&basePrice)
	if err != nil {
		log.Println("获取原始价格失败:", err)
		return global.ResponseDataNotFound, 0
	}
	// 2. 计算 10% 随机浮动 (整数)
	floatRange := int(float64(basePrice) * 0.1)
	finalPrice := basePrice
	if floatRange > 0 {
		// 使用你的 RandomRange 生成偏移
		offset := Util.RandomRange(-floatRange, floatRange)
		finalPrice = basePrice + offset
	}
	return global.ResponseSuccess, finalPrice
}

func (r *User_repo_impl) AddCardInBags(ctx context.Context, db repo.SQLQueryer, cardID int, userID int) global.ResponseStatusCode {
	// 价格
	err, price := r.GetCardPrice(ctx, db, cardID)
	if err != global.ResponseSuccess {
		return err
	}
	query := `insert into bags (user_id, card_id, price) values ($1, $2, $3)`
	_, err2 := db.ExecContext(ctx, query, userID, cardID, price)
	if err2 != nil {
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

func (r *User_repo_impl) JudgeCardIsParent(ctx context.Context, db repo.SQLQueryer, CardId int) (global.ResponseStatusCode, bool) {
	query := "select category from newcards where id = $1"

	var category int
	err := db.QueryRowContext(ctx, query, CardId).Scan(&category)
	if err != nil {
		fmt.Println(err)
		if err == sql.ErrNoRows {
			return global.ResponseDataNotFound, false
		}
		return global.ResponseInternalServersError, false
	}

	if category == 1 || category == 2 {
		return global.ResponseSuccess, true
	} else {
		return global.ResponseSuccess, false
	}

}
func (r *User_repo_impl) JudgeCardIsCharacter(ctx context.Context, db repo.SQLQueryer, CardId int) (global.ResponseStatusCode, bool) {
	query := "select category from newcards where id = $1"
	var category int
	err := db.QueryRowContext(ctx, query, CardId).Scan(&category)
	if err != nil {
		fmt.Println(err)
		if err == sql.ErrNoRows {
			return global.ResponseDataNotFound, false
		}
		return global.ResponseInternalServersError, false
	}
	if category == 1 || category == 3 {
		return global.ResponseSuccess, true
	} else {
		return global.ResponseSuccess, false
	}

}

// CreateBattle 创建战斗并返回数据库自动生成的 battle_id
func (r *User_repo_impl) CreateBattle(ctx context.Context, db repo.SQLQueryer, playerIdA int, playerIdB int) (int, global.ResponseStatusCode) {
	var battleId int
	query := `insert into battle (player_ida, player_idb) values ($1, $2) returning battle_id`
	err := db.QueryRowContext(ctx, query, playerIdA, playerIdB).Scan(&battleId)
	if err != nil {
		log.Printf("failed to create battle: %v", err)
		return 0, global.ResponseInternalServersError
	}
	return battleId, global.ResponseSuccess
}

// 检查并且返回 battleid
// 返回值：battleId (int), status (global.ResponseStatusCode)
func (r *User_repo_impl) CheckUserIdIsBattle(ctx context.Context, db repo.SQLQueryer, userId int) (int, global.ResponseStatusCode) {
	query := `select battle_id from battle where player_ida = $1 or player_idb = $1 limit 1`

	var battleId int
	err := db.QueryRowContext(ctx, query, userId).Scan(&battleId)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return 0, global.ResponseDataNotFound
		}
		fmt.Println("CheckUserIdIsBattle:", err)
		return 0, global.ResponseInternalServersError
	}

	return battleId, global.ResponseSuccess
}

// DeleteBattle 根据 battle_id 删除战斗记录
func (r *User_repo_impl) DeleteBattle(ctx context.Context, db repo.SQLQueryer, BtId int) global.ResponseStatusCode {
	// 根据 battle_id 删除对应记录
	query := `delete from battle where battle_id = $1`

	result, err := db.ExecContext(ctx, query, BtId)
	if err != nil {
		log.Printf("failed to delete battle %d: %v", BtId, err)
		return global.ResponseInternalServersError
	}

	// 检查是否有行被删除
	rows, err := result.RowsAffected()
	if err != nil {
		fmt.Println(err)
		return global.ResponseInternalServersError
	}

	if rows == 0 {
		fmt.Println("battle_id 不存在")
		return global.ResponseDataNotFound
	}

	return global.ResponseSuccess
}

// CreateLoot 创建 loot 记录
func (r *User_repo_impl) CreateLoot(ctx context.Context, db repo.SQLQueryer, loot []int, UserId int) global.ResponseStatusCode {
	query := `insert into loot (userId, data) values ($1, $2)`

	_, err := db.ExecContext(ctx, query, UserId, loot)
	if err != nil {
		log.Printf("failed to create loot for user %d: %v", UserId, err)
		return global.ResponseInternalServersError
	}

	return global.ResponseSuccess
}

type IntSlice []int

// 用来把psql的数组(json字符串)转化成goland的数组
func (s *IntSlice) Scan(value interface{}) error {
	if value == nil {
		*s = nil
		return nil
	}

	var str string
	switch v := value.(type) {
	case string:
		str = v
	case []byte:
		str = string(v)
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}

	// 去除首尾的 postgres 数组大括号 { 和 }
	str = strings.TrimSpace(str)
	str = strings.TrimPrefix(str, "{")
	str = strings.TrimSuffix(str, "}")

	if str == "" {
		*s = []int{}
		return nil
	}

	// 按逗号分割并转成 int
	parts := strings.Split(str, ",")
	res := make([]int, len(parts))
	for i, p := range parts {
		val, err := strconv.Atoi(strings.TrimSpace(p))
		if err != nil {
			return fmt.Errorf("failed to convert %q to int: %v", p, err)
		}
		res[i] = val
	}
	*s = res
	return nil
}

// GetLoot 获取用户的 loot 记录列表
func (r *User_repo_impl) GetLoot(ctx context.Context, db repo.SQLQueryer, UserId int) (global.ResponseStatusCode, []BattleData.LootDto) {

	query := `select lootId, data from loot where userid = $1`

	rows, err := db.QueryContext(ctx, query, UserId)
	if err != nil {
		log.Printf("failed to get loot for user %d: %v", UserId, err)
		return global.ResponseInternalServersError, nil
	}
	defer rows.Close()

	loots := make([]BattleData.LootDto, 0)
	for rows.Next() {
		var item BattleData.LootDto
		// 使用之前定义好的 IntSlice 来接收 data 字段
		var dataSlice IntSlice

		if err := rows.Scan(&item.LootID, &dataSlice); err != nil {
			log.Printf("failed to scan loot data for user %d: %v", UserId, err)
			return global.ResponseInternalServersError, nil
		}

		item.Data = dataSlice
		loots = append(loots, item)
	}

	if err := rows.Err(); err != nil {
		log.Printf("rows error getting loot for user %d: %v", UserId, err)
		return global.ResponseInternalServersError, nil
	}
	return global.ResponseSuccess, loots
}

func (r *User_repo_impl) DeleteLoot(ctx context.Context, db repo.SQLQueryer, LootId int) global.ResponseStatusCode {
	query := `delete from loot where lootId = $1`
	result, err := db.ExecContext(ctx, query, LootId)
	if err != nil {
		log.Printf("failed to delete loot for user %d: %v", LootId, err)
		return global.ResponseInternalServersError
	}
	rows, err := result.RowsAffected()
	if err != nil {
		fmt.Println(err)
		return global.ResponseInternalServersError
	}
	if rows == 0 {
		return global.ResponseDataNotFound
	}
	return global.ResponseSuccess
}
func (r *User_repo_impl) GetGoodsByUserId(ctx context.Context, db repo.SQLQueryer, UserId int) (global.ResponseStatusCode, []BattleData.GoodsDto) {
	query := `select goods_id, card_id, price, discount from goods where user_id = $1`

	rows, err := db.QueryContext(ctx, query, UserId)
	if err != nil {
		log.Printf("failed to get goods for user %d: %v", UserId, err)
		return global.ResponseInternalServersError, nil
	}
	defer rows.Close()

	goodsList := make([]BattleData.GoodsDto, 0)
	for rows.Next() {
		var item BattleData.GoodsDto
		if err := rows.Scan(&item.GoodsId, &item.CardId, &item.Price, &item.Discount); err != nil {
			log.Printf("failed to scan goods data for user %d: %v", UserId, err)
			return global.ResponseInternalServersError, nil
		}
		goodsList = append(goodsList, item)
	}

	if err := rows.Err(); err != nil {
		log.Printf("rows error getting goods for user %d: %v", UserId, err)
		return global.ResponseInternalServersError, nil
	}
	return global.ResponseSuccess, goodsList
}

func (r *User_repo_impl) DeleteGoodsByUserId(ctx context.Context, db repo.SQLQueryer, UserId int) global.ResponseStatusCode {
	query := `delete from goods where user_id = $1`
	result, err := db.ExecContext(ctx, query, UserId)
	if err != nil {
		log.Printf("failed to delete goods for user %d: %v", UserId, err)
		return global.ResponseInternalServersError
	}
	num, err2 := result.RowsAffected()
	if err2 != nil {
		fmt.Println(err)
		return global.ResponseInternalServersError
	}
	if num == 0 {
		return global.ResponseDataNotFound
	}
	return global.ResponseSuccess
}

func (r *User_repo_impl) CreateGoods(ctx context.Context, db repo.SQLQueryer, UserId int, GoodsList []*BattleData.GoodsDto) global.ResponseStatusCode {
	if len(GoodsList) == 0 {
		return global.ResponseSuccess
	}

	query := `insert into goods (user_id, card_id, price, discount) values ($1, $2, $3, $4)`

	for _, goods := range GoodsList {
		_, err := db.ExecContext(ctx, query, UserId, goods.CardId, goods.Price, goods.Discount)
		if err != nil {
			log.Printf("failed to create goods for user %d: %v", UserId, err)
			return global.ResponseInternalServersError
		}
	}
	return global.ResponseSuccess
}
