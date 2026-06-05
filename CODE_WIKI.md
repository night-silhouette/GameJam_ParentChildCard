# Napori Game - 卡牌游戏项目 Wiki

## 项目概述

Napori Game 是一款基于 Godot 4.5 引擎开发的卡牌对战游戏，采用 **子母牌** 系统设计理念。玩家在对局中体验奇思妙想的卡片效果与紧张刺激的对局比拼。

- **游戏引擎**: Godot 4.5 (GDScript)
- **后端框架**: Go (Gin Web框架)
- **数据库**: PostgreSQL + Redis
- **网络协议**: HTTP REST API + WebSocket 实时通信

---

## 项目架构总览

```
Napori_Game/
├── base_class/           # 核心基础类（Autoload 全局单例）
│   ├── card/            # 卡牌相关基类
│   ├── manage/          # 管理器（场景、UI、数据、库存）
│   ├── net/             # 网络层（HTTP、WebSocket）
│   ├── signal/          # 信号总线
│   └── global.gd        # 全局工具函数
├── game_scene/          # 游戏场景（场景文件 + 逻辑脚本）
│   ├── battle/          # 战斗场景
│   ├── bag/             # 背包场景
│   ├── menu/            # 菜单场景
│   └── main/            # 主入口场景
├── game_ui/             # 游戏 UI 层
│   ├── battle/          # 战斗 UI
│   ├── login/           # 登录 UI
│   └── match/           # 匹配 UI
├── game_data/           # 游戏数据（CSV 导入、卡牌资源）
├── servers/             # Go 后端服务
│   └── src/
│       ├── application/ # 业务逻辑层
│       ├── presentation/ # 控制器/处理器层
│       ├── infra/        # 基础设施层
│       └── cmd/          # 入口点
└── 素材/                # 游戏资源素材
```

---

## 核心模块详解

### 1. Autoload 全局单例（base_class/）

| 名称 | 文件路径 | 职责 |
|------|----------|------|
| **Global** | `base_class/global.gd` | 全局工具函数、区域定义(ZONE_CARD)、假死/复活节点接口 |
| **SignalBus** | `base_class/signal/signal_bus.gd` | 全局信号总线，所有模块间通信通过信号 |
| **ScenceManage** | `base_class/manage/scence_manage.gd` | 场景切换管理 |
| **UiManage** | `base_class/manage/ui_manage.gd` | UI 切换管理 |
| **NetworkClient** | `base_class/net/Http/Network_Client.gd` | HTTP 请求客户端 |
| **ApiManager** | `base_class/net/Http/Api_manager.gd` | API 请求封装 |
| **Packethandler** | `base_class/net/Http/Packethandler.gd` | HTTP 响应解析与路由 |
| **BattleWs** | `base_class/net/WS/BattleWS.gd` | WebSocket 战斗连接 |
| **BattleRequestManager** | `base_class/net/WS/Battle_RequestManager.gd` | 战斗请求发送 |
| **BattleResponseManager** | `base_class/net/WS/Battle_ResponseManager.gd` | 战斗响应处理 |
| **TokenManager** | `base_class/net/token_manager.gd` | 登录 Token 管理 |
| **DataManagerBt** | `base_class/manage/data_manager_BT.gd` | 战斗数据管理 |
| **InventoryManager** | `base_class/manage/inventory_manager.gd` | 背包/库存数据管理 |

---

### 2. 卡牌系统（Card System）

#### 2.1 卡牌类层次

```
card (base_class/card/card.gd)
├── battle_card (base_class/card/battle_card.gd)    # 战斗卡牌（可进入战斗区）
├── match_card (base_class/card/match_card.gd)      # 匹配卡牌
├── bag_card (base_class/card/bag_card.gd)          # 背包卡牌
└── choose_card (base_class/card/choose_card.gd)    # 选牌卡牌
```

#### 2.2 CardResource 资源类

**文件**: `base_class/card/card_data.gd`

```gdscript
class_name CardResource extends Resource

# 视觉资源
@export var card_texture: Texture2D
@export var texture_filename: String

# 基础数据
@export var name: String
@export var id: int
@export var level: int
@export var value: int

# 战斗属性
@export var damage: int
@export var initial_health: int
@export var max_health: int

# 技能设置
@export var skill_charge: int
@export var skill_card_use_num: int

# 卡牌类型
@export var is_combat_card: bool  # 是否为战斗牌
@export var is_sub_card: bool      # 是否为子牌

# 描述
@export_multiline var skill_description: String
@export_multiline var notes: String
@export_multiline var sub_card_trigger_effect: String
```

#### 2.3 卡牌区域 (Zone)

**文件**: `base_class/global.gd` - `ZONE_CARD` 字典

| Zone | 值 | 说明 |
|------|-----|------|
| DECK_ZONE | 1 | 手牌/牌堆 |
| DISCARD_ZONE | 2 | 弃牌堆 |
| PARENT_BATTLE_ZONE | 3 | 己方母牌战斗区 |
| CHILD_BATTLE_ZONE | 4 | 己方子牌战斗区 |
| SPELL_ZONE | 5 | 己方法术牌区 |
| FREE_ZONE | 6 | 自由区域（拖拽中） |
| ENEMY_PARENT_ZONE | 7 | 敌方母牌战斗区 |
| ENEMY_CHILD_ZONE | 8 | 敌方子牌区 |
| ENEMY_SPELL_ZONE | 9 | 敌方法术牌区 |
| BAG_ZONE | 10 | 背包区域 |
| SELL_ZONE | 11 | 卖出区域 |
| MATCH_ZONE | 12 | 出战区域 |
| CHILD_ACTIVE | 13 | 子牌已激活 |
| CHILD_NOT_ACTIVE | 14 | 子牌未激活 |
| CHILD_DIED | 15 | 子牌已死亡 |
| CHILD_HAS_CATCH | 16 | 子牌已被捕获 |

---

### 3. 信号系统 (SignalBus)

**文件**: `base_class/signal/signal_bus.gd`

信号采用 **region** 组织，按功能分为：

#### 3.1 游戏信号 (Game)
```gdscript
signal change_scence(path: String)      # 切换场景
signal change_ui(name: String)          # 切换UI
signal network_disconnected()            # 网络断开
signal battle_information               # 战斗信息
signal online_match                    # 在线匹配
```

#### 3.2 HTTP 信号
```gdscript
# 请求信号
signal request_login(username, password)
signal request_register_user(name, password)
signal request_bag_card
signal request_card_random
signal request_get_self_gold
signal request_sell_card(card_list)

# 响应信号
signal login_success()
signal login_failed(msg)
signal get_card_bag(card_list)
signal get_self_gold(gold)
```

#### 3.3 WebSocket 战斗信号
```gdscript
# 连接状态
signal ws_connected
signal ws_disconnected

# 战斗状态
signal battle_started
signal battle_over
signal match_success(t)

# 卡牌数据
signal self_inhand_updated(cards)
signal bt_selfinfo_updated(cards)
signal bt_oppinfo_updated(cards)
signal energy_updated(energy_list)
signal child_card_list_updated(child_cards)

# 阶段信号
signal magic_card_start(t)
signal combat_start_success(t, is_win)
signal judge_start(t)
signal judge_finish
```

---

### 4. 战斗系统 (Battle System)

#### 4.1 战斗状态机

**文件**: `game_scene/battle/state_machine.gd`

```gdscript
enum GameState {
    INIT_STATE,           # 看牌阶段（可拖拽上牌）
    CHOOSE_CHILD_CARD,    # 选择子卡牌阶段
    CHOOSE_WEATHER,       # 选择天气阶段
    USE_MAGIC_CARD,       # 选择技能牌阶段
    USE_COMBAT_CARD,      # 战斗行动阶段
    JUDGEMENT             # 判定阶段（剪刀石头布）
}
```

#### 4.2 战斗流程

1. **INIT_STATE** - 初始阶段，玩家可部署母牌和法术牌
2. **CHOOSE_CHILD_CARD** - 选择要激活的子卡
3. **CHOOSE_WEATHER** - 选择天气效果
4. **USE_MAGIC_CARD** - 使用技能/法术牌
5. **USE_COMBAT_CARD** - 执行战斗操作
6. **JUDGEMENT** - 回合判定（剪刀石头布）

---

### 5. 网络层 (Network)

#### 5.1 HTTP API 端点

**服务器基础 URL**: `http://120.26.145.68:10086`

| 端点 | 方法 | 功能 |
|------|------|------|
| `/v1/token/` | POST | 登录 |
| `/v1/token/` | GET | 验证 Token |
| `/v1/user/` | GET | 获取用户信息 |
| `/v1/user/` | POST | 注册用户 |
| `/v1/bags/` | GET | 获取背包卡牌 |
| `/v1/bags/sell/` | POST | 出售卡牌 |
| `/v1/start_pack/` | GET | 开随机包 |
| `/v1/user/gold/` | GET | 获取金币 |
| `/v1/mail/` | GET/POST | 邮件操作 |

#### 5.2 WebSocket 战斗协议

**文件**: `base_class/net/WS/Net_def.gd`

```gdscript
enum Action {
    FAULT = 0,
    CANCEL_MATCH = 1,
    GET_SELF_CARDS = 2,
    GET_OPPONENT_CARDS = 3,
    GET_BT_INFO = 4,
    OVER_BATTLE = 5,
    START_BATTLE = 6,
    DEPLOY_CARD = 7,
    JUDGE = 8,
    MATCH_SUCCESS = 9,
    ANIMATION_END = 10,
    COMBAT = 11,
    CardCalc = 12,
    Interrupt = 14,
    GetEnergy = 17,
    SelectWeather = 18,
    GetChildCardList = 19,
    ActiveChildCard = 20
}

enum Predicate {
    EMPTY = 0,
    NOTIFY = 1,
    QUERY = 2,
    RESULT = 3,
    FINISH = 4,
    SUCCEED = 5
}
```

---

### 6. 后端服务 (Go Server)

#### 6.1 项目结构

```
servers/src/
├── cmd/main/main.go           # 入口函数
├── application/
│   ├── entity/               # 实体定义
│   │   ├── Card/            # 卡牌实体
│   │   │   ├── CardAbstract/ # 卡牌接口
│   │   │   ├── CardImpl/     # 卡牌实现 (00-44)
│   │   │   └── character_card/
│   │   ├── BattleData/       # 战斗数据传输对象
│   │   └── User_entity/      # 用户实体
│   └── service/              # 业务服务
│       ├── UserService/      # 用户服务
│       └── battleservice/     # 战斗服务
├── presentation/
│   ├── handler/              # HTTP 处理器
│   ├── route/                # 路由定义
│   └── response/             # 响应封装
├── infra/                     # 基础设施
│   ├── config/               # 配置
│   ├── db/                   # 数据库连接
│   ├── redisConnect/          # Redis 连接
│   └── repo/                 # 数据仓库
└── global/                    # 全局定义
```

#### 6.2 核心服务

**BattleService** (`application/service/battleservice/`)

- `Battle.go` - 战斗房间管理
- `Match.go` - 匹配系统（带权重的等待时间匹配算法）
- `BattleService.go` - 战斗业务接口
- `state.go` - 战斗状态机
- `EffectStack.go` - 效果结算栈

#### 6.3 依赖项

```go
require (
    github.com/gin-gonic/gin v1.11.0      // Web 框架
    github.com/golang-jwt/jwt/v5 v5.3.1   // JWT 认证
    github.com/gorilla/websocket v1.5.3   // WebSocket
    github.com/jackc/pgx/v5 v5.8.0        // PostgreSQL
    github.com/redis/go-redis/v9 v9.18.0  // Redis
    golang.org/x/crypto v0.41.0           // 加密
)
```

---

### 7. 关键类与函数

#### 7.1 Global (base_class/global.gd)

```gdscript
# 节点假死/复活
func fake_death(target_node: Node) -> void
func revive(target_node: Node) -> void

# 获取区域优先级
func _get_zone_priority(zone: int) -> int
func _should_override_zone(old_zone, new_zone) -> bool
```

#### 7.2 DataManagerBt (base_class/manage/data_manager_BT.gd)

```gdscript
# 卡牌数据管理
var card_list: Array                          # 所有卡牌数据
var self_energy: int                          # 己方能量
var opponent_energy: int                     # 敌方能量
var weather_data: Dictionary                 # 天气数据

# 信号
signal UI_date_update                        # UI数据更新
signal change_card_zone(temp_id, new_zone)   # 区域变更
signal energy_changed(self_energy, opponent_energy)

# 方法
func querry_resoure_by_id(card_id: int) -> Resource
func get_cards_by_zone(zone) -> Array
func select_card_by_key(value, key_to_match)
func toggle_selection(match_code, temp_id) -> bool
```

#### 7.3 InventoryManager (base_class/manage/inventory_manager.gd)

```gdscript
# 背包数据管理
var card_list: Array     # 背包卡牌列表
var gold: int           # 金币数量

# 方法
func move_to_sell_zone(stuff_id: int) -> bool
func move_to_bag_zone(stuff_id: int) -> bool
func move_to_combat_zone(stuff_id: int) -> bool
func _find_card_resource_by_id(target_id: int) -> CardResource
```

#### 7.4 BattleWs (base_class/net/WS/BattleWS.gd)

```gdscript
var ws: WebSocketPeer
var is_connected: bool

func _connect_ws(body)              # 建立 WS 连接
func send_action(action_code, action_data, predicates)  # 发送动作
```

---

### 8. 数据流

#### 8.1 登录流程

```
用户输入 → request_login 信号
         → NetworkClient.call_api("/v1/token/", POST)
         → 服务器返回 Token
         → Packethandler._handle_raw_api_data()
         → TokenManager.save_token()
         → SignalBus.login_success.emit()
```

#### 8.2 战斗流程

```
匹配成功 → SignalBus.match_success.emit(t)
        → state_machine.change_state(INIT_STATE)
        → 玩家拖拽卡牌到战斗区
        → SignalBus.enter_freecard.emit(temp_id, zone)
        → DataManagerBt._sync_card_data()
        → 战斗回合判定
        → SignalBus.judge_start.emit(t)
        → 玩家选择手势
        → SignalBus.request_judge.emit(judge_data)
        → 服务器计算胜负
        → SignalBus.combat_start_success.emit()
```

---

### 9. 项目运行

#### 9.1 前端 (Godot)

1. 使用 Godot 4.5 打开项目
2. 项目配置文件: `project.godot`
3. 运行主场景: `game_scene/main/main.gd`

#### 9.2 后端 (Go)

```bash
cd servers/src/cmd/main
go run main.go
```

或使用部署脚本:
```bash
bash deploy.sh
```

#### 9.3 配置

后端配置目录: `servers/src/infra/config/info/`

- `DB_jnfo.json.example` - 数据库配置示例
- `route_info.json` - 路由配置
- `SecretKey.go` - 密钥配置

---

### 10. 开发指南

#### 10.1 添加新卡牌

1. 在 `game_data/card/` 目录创建 `card_XX.tres` 文件
2. 或使用 `CardataLoad` 工具从 CSV 导入
3. 设置卡牌属性（id, name, damage, is_combat_card 等）

#### 10.2 添加新技能

后端 Go 实现:
1. 在 `application/entity/Card/CardImpl/` 创建新文件
2. 实现 `Card` 接口
3. 在 `application/service/battleservice/cardlist.go` 注册

#### 10.3 添加新 UI

1. 在 `game_ui/` 创建场景 `.tscn` 和脚本 `.gd`
2. 在 `UiManage.gd` 的 `fui_change()` 添加路由

#### 10.4 场景切换

```gdscript
# 切换场景
SignalBus.change_scence.emit("tomenu")  # 菜单
SignalBus.change_scence.emit("tobattle") # 战斗
SignalBus.change_scence.emit("bag")      # 背包

# 切换 UI
SignalBus.change_ui.emit("tologin")
SignalBus.change_ui.emit("tobattle")
```

---

### 11. 状态码 (StatusCode)

**文件**: `servers/src/global/StatusCode.go`

| 状态码 | 说明 |
|--------|------|
| ResponseSuccess (0) | 成功 |
| ResponseDataNotFound | 数据未找到 |
| ResponseInvalidToken | 无效 Token |
| ResponseTokenExpired | Token 过期 |
| BattleInvalidTiming | 不在正确的战斗时机 |
| BattleCardNotFound | 卡牌未找到 |
| BattleNotInYourRound | 不在己方回合 |

---

### 12. 文件命名约定

| 类型 | 前缀 | 示例 |
|------|------|------|
| 卡牌资源 | card_ | `card_01.tres`, `card_44.tres` |
| 卡牌实现 (Go) | 数字 | `01.go`, `02.go`, ..., `44.go` |
| 场景文件 | 无 | `battle_scence.tscn` |
| 脚本文件 | 无 | `battle_scence.gd` |

---

*最后更新: 2026-06-05*
