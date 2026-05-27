package battlerepo

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"pcc_card/Util"
	"pcc_card/application/entity/BattleData"
	"pcc_card/global"
	"pcc_card/infra/repo"

	"github.com/redis/go-redis/v9"
)

type BattleRepo interface {
	repo.Repo
	ReadCardByID(ctx context.Context, db repo.SQLQueryer, ID int) map[string]any
	AddCardInBags(ctx context.Context, db repo.SQLQueryer, cardID int, userID int) global.ResponseStatusCode
	GetBagsByUserId(ctx context.Context, db repo.SQLQueryer, userID int) ([]BattleData.BagStuffDto, global.ResponseStatusCode)
}

type BattleRepoImpl struct {
	db *sql.DB
	rd *redis.Client
}

func (r *BattleRepoImpl) Get_db() *sql.DB {
	return r.db
}

func (r *BattleRepoImpl) Set_db(db *sql.DB, rd *redis.Client) {
	r.db = db
	r.rd = rd
}

func (r *BattleRepoImpl) ReadCardByID(ctx context.Context, db repo.SQLQueryer, ID int) map[string]any {
	var info []byte
	var res map[string]any

	// 仅仅修改了 query，利用 json_build_object 把字段打包
	// 这样数据库返回的就是一个完整的 JSON 字符串，刚好能塞进你的 info []byte
	query := `select json_build_object(
        'damage', damage, 
        'initHp', "initHp", 
        'maxHp', "maxHp", 
        'price', price, 
        'skillCharge', "skillCharge", 
        'skillcardUseNum', "skillcardUseNum", 
        'category', category
    ) from newcards where id = $1`

	data := db.QueryRowContext(ctx, query, ID)
	err := data.Scan(&info)
	if err != nil {
		log.Println(err)
	}
	err = json.Unmarshal(info, &res)
	if err != nil {
		log.Println(err)
	}
	if res["category"] == 1 || res["category"] == 2 {
		res["is_parent"] = true
	}
	if res["category"] == 3 || res["category"] == 4 {
		res["is_parent"] = false
	}

	return res
}
func (r *BattleRepoImpl) AddCardInBags(ctx context.Context, db repo.SQLQueryer, cardID int, userID int) global.ResponseStatusCode {
	// 1. 获取基准价格
	var basePrice int
	queryPrice := `select price from newcards where id = $1`
	err := db.QueryRowContext(ctx, queryPrice, cardID).Scan(&basePrice)
	if err != nil {
		log.Println("获取原始价格失败:", err)
		return global.ResponseBagsUnknownError
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

func (r *BattleRepoImpl) GetBagsByUserId(ctx context.Context, db repo.SQLQueryer, userID int) ([]BattleData.BagStuffDto, global.ResponseStatusCode) {
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
