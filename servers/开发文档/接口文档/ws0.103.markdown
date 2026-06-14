
# 0.103版本号接口文档
	更新日志:
    1.增加了查看天气的公共接口,actionCode为GetWeather
    2.子牌获取的通知(cardcal阶段的),ChildCardCatchDto
    3.公共接口(GetUserId)双方的UserId,初始状态触发之后获取一次之后记录下来就好了
    4.cardcalc大改
    5.buffdto改了
    6.where字段变过了

## 一.action_code映射表
|action_code|action_name|value|注解|
| :---: | :---: | :---: | :---  |
|Fault|错误|0|接口里没有用到，只是默认值为0，所以如果传默认值就认为是错误|
| [CancelMatch](#1cancelmatch)|取消匹配|1||
| [GetSelfCardInHard](#3getselfcardinhard)|获取自己的卡牌信息|2||
|[GetOpponentCardInHard](#6getopponentcardinhard)|获取对手的卡牌信息|3|这个也只是调试用的，正常不能给用户获取对手手牌信息|
| [GetBtCardInfo](#4getbtcardinfo)|获取场上的战斗信息|4||
| [OverBattle](#2overbattle)|结束战斗|5||
|StartBattle|开始战斗|6||
| [DeployCard](#1deploycard)|部署一张牌|7|这个actioncode我说实话设计的不好，但也不会改他了，现在他的作用就两个，一个是在状态机为选择技能牌的阶段的时候用，这个链接点进去也是讲这个的，然后还有就是游戏刚开始的看牌阶段和combat阶段，拖动上牌用这个。就这两个用处|
| [Judge](#1judge)|战斗回合判断|8||
| [MatchSuccess](#1matchsuccess)|匹配成功|9||
|AnimationPlayEnd|动画结束|10||
| [Combat](#1combat)|执行战斗行动|11||
| [CardCalc](#1cardcalc)|卡牌效果结算|12||
|Debug|测试|13|这个是我自己要查看一些运行时数据用的，没有你要用的正式接口|
|Interrupt|中断选牌|14||
| [GetDisCard](#5getdiscard)|查看弃牌堆|15||
|SkillCardCalc|法术牌计算|16||
|[GetEnergy](#7getenergy)|查看能量值|17||
| [SelectWeather](#4-selectweather-阶段)|选择天气|18|新增：钱多者选择天气的阶段（看牌后）|
|[GetChildCardList](#8getchildcardlist)|查看子牌堆|19||
| [ActiveChildCard](#3-activechildcard-阶段)|激活子卡牌|20|新增：双方选择并激活子卡牌的阶段（看牌后）|
---

## 二.具体接口说明

### <1> 常用请求
#### 1.CancelMatch
|action_code|方向|Predicates|传参(action_data)|注释|
| :---: | :---: | :---: | :---: | :---: |
|CancelMatch|客户端->服务器|Notify|空|还在匹配的时候传，字面意思，取消匹配，ws会直接断掉|
---

#### 2.OverBattle
|action_code|方向|Predicates|传参(action_data)|注释|
| :---: | :---: | :---: | :---: | :---: |
|OverBattle|客户端->服务器|Notify|空|传了就会直接结束游戏，双方都会退出来，ws会直接断掉|
---
#### 3.GetSelfCardInHard
|action_code|方向|Predicates|传参(action_data)|注释|
| :---: | :---: | :---: | :---: | :---: |
|GetSelfCardInHard|客户端->服务器|Query|空||
|GetSelfCardInHard|服务器->客户端|Result|[SelfCardDtoList](#selfcarddtolist)||
---

#### 4.GetBtCardInfo
|action_code|方向|Predicates|传参(action_data)|注释|
| :---: | :---: | :---: | :---: | :---: |
|GetBtCardInfo|客户端->服务器|Query|空||
|GetBtCardInfo|服务器->客户端|Result|[BtCardInfo](#btcardinfo)||
---

#### 5.GetDisCard
|action_code|方向|Predicates|传参(action_data)|注释|
| :---: | :---: | :---: | :---: | :---: |
|GetDisCard|客户端->服务器|Query|空||
|GetDisCard|服务器->客户端|Result|[DisCardDtoList](#discardolist)||
---

#### 6.GetOpponentCardInHard
|action_code|方向|Predicates|传参(action_data)|注释|
| :---: | :---: | :---: | :---: | :---: |
|GetOpponentCardInHard|客户端->服务器|Query|空||
|GetOpponentCardInHard|服务器->客户端|Result|[][CardDto](#carddto)|获取对手的手牌信息（调试用）|
---

#### 7.GetEnergy
|action_code|方向|Predicates|传参(action_data)|注释|
| :---: | :---: | :---: | :---: | :---: |
|GetEnergy|客户端->服务器|Query|空||
|GetEnergy|服务器->客户端|Result|map{"self":int, "opponent":int}|返回自己和对手的当前能量值|
---

#### 8.GetChildCardList
|action_code|方向|Predicates|传参(action_data)|注释|
| :---: | :---: | :---: | :---: | :---: |
|GetChildCardList|客户端->服务器|Query|空||
|GetChildCardList|服务器->客户端|Result|[][ChildCardDto](#childcarddto)|获取所有可用的子卡牌列表|
---

### <2> 看牌阶段
#### 1.MatchSuccess
|action_code|方向|Predicates|传参(action_data)|注释|
| :---: | :---: | :---: | :---: | :---: |
|MatchSuccess|服务器->客户端|Notify|[StateWaitTime](#statewaittime)|收到这个信息表明匹配成功了，现在是特殊的看牌阶段,这个阶段不会快速结束，即使大家都上完牌了|
|DeployCard|客户端->服务器|Result|[SelectCard](#selectcard)|这个阶段和combat阶段一样，可以拖动上牌，这个就是上牌的接口|
|DeployCard|服务器->客户端|Success|空|收到这个表示你上牌成功了，值得一提的是，这个阶段时间到了会自动上牌,自动上牌也会提醒success有这条消息，至于自动上了什么，可以通过GetBtCardInfo查看|
|StartBattle|服务器->客户端|Notify|空|直到收到这个，表示这个阶段结束了（也可以直接去监听下一阶段的开始信号）|

---

### <3> ActiveChildCard 阶段
|action_code|方向|Predicates|传参(action_data)|注释|
| :---: | :---: | :---: | :---: | :---: |
|ActiveChildCard|服务器->客户端|Query|{"state_wait_time":StateWaitTime, "child_list":[]CardDto}|通知双方进入子卡激活选择阶段，包含倒计时与可选子卡列表|
|ActiveChildCard|客户端->服务器|Result|[ActiveChildCardDto](#activechildcarddto)|客户端返回想要激活的子卡 `temp_id_list`|
|ActiveChildCard|服务器->客户端|Succeed|空|服务器确认接收选择|
|ActiveChildCard|服务器->客户端|Finish|{"selected_temp_id_list": []int}|选择结束，返回最终被激活的子卡列表（双方交集或补齐后结果）|

在超时情况下，服务器会根据双方选择的交集补齐或随机补足到固定数量，并返回 Finish 后进入 `SelectWeather` 状态。

### <4> SelectWeather 阶段
|action_code|方向|Predicates|                     传参(action_data)                     |注释|
| :---: | :---: | :---: |:-------------------------------------------------------:| :---: |
|SelectWeather|服务器->客户端|Notify|                    {"is_more": bool}                    |告知双方谁在本回合钱更多（有优先选择权）|
|SelectWeather|服务器->客户端|Query| {"state_wait_time":StateWaitTime, "weather_list":[]int} |向拥有优先权的玩家请求选择天气，包含可选天气列表与倒计时|
|SelectWeather|客户端->服务器|Result|          [SelectWeatherDto](#selectweatherdto)          |玩家返回所选天气|
|SelectWeather|服务器->客户端|Succeed|                map{"weather":weather枚举}                 |服务器确认选择成功并广播结果给对手（若超时由系统随机选择）|

服务器在倒计时结束仍未选择时会随机选择一个天气并通知双方，然后进入下一状态 `SelectSkillCard`。


### <5> 选择技能牌阶段

#### 1.DeployCard
|action_code|方向|Predicates|传参(action_data)|注释|
| :---: | :---: | :---: | :---: | :---: |
|DeployCard|服务器->客户端|Query|[SelectSkill](#SelectSkill)|这是开始进入选技能牌阶段的标志|
|DeployCard|客户端->服务器|Result|[SelectCard](#selectcard)|传where是skillcard的SelectCard,特殊的一个用法是，如果你传cardTempId=-1，那表示你明确不上技能,这样双方都确认了之后就可快速的finish|
|DeployCard|服务器->客户端|Success|空|选成功了会给你这个,但是要等Finish，否则就是对方还在选，技能牌可上可不上，所以即使时间到了也不会系统自动选择|
|DeployCard|服务器->客户端|Finish|空|选结束了，如果双方都在倒计时结束之前结束，那他会提前结束|



---


### <6> Judge阶段
#### 1.Judge
赢的人能量会变多

|action_code|方向|Predicates|传参(action_data)|注释|
| :---: | :---: | :---: | :---: | :---: |
|Judge|服务器->客户端|Query|[StateWaitTime](#statewaittime)|开始回合判断，向客户端询问|
|Judge|客户端->服务器|Result|[JudgeData](#JudgeData)|你要传的，剪刀石头布|
|Judge|服务器->客户端|Success|空|传完上一个之后，就会返回你条这个，但是你那边要继续显示对方还在选择，直到收到Finish|
|Judge|服务器->客户端|Finish|[JudgeRes](#judgeres)|时间到了或者双方都选完了.就会返回你这个,如果是平局，会重新进入Judge状态，又会发你query，有胜负就会进入Combat阶段|
|AnimationPlayEnd|客户端->服务器|Notify|空|每次收到finish之后，剪刀石头布的动画结束之后，传我一个这个，后端状态机才会恢复继续运作|
---
### <7> Combat阶段

说明：`Combat.Result` 应当为一个 `[]CombatDto`（数组），服务器会分别收集胜者（Query 发起方）和败者（Notify 接收方）提交的数组并保存在内部的 `CombatMap`（键为 "Winner" 与 "Loser"）中。每次收到一方有效提交后，服务器会返回 `Succeed` 确认并防止重复提交；当双方都提交完毕或计时器到期时，服务器会把 `CombatMap` 通过内部通道发送到 `CardCalc`（文档中的 `CardCalc` 阶段）以开始结算。超时未提交的一方将视为空提交，服务器仍会继续结算并进入下一步。

#### 1.Combat
|action_code|方向|Predicates|传参(action_data)|注释|
| :---: | :---: | :---: | :---: | :---: |
|Combat|服务器->客户端|Query|[StateWaitTime](#statewaittime)|发给回合胜者（Winner），请求提交本回合的行动数组|
|Combat|服务器->客户端|Notify|[StateWaitTime](#statewaittime)|发给回合失败者（Loser）|
|Combat|客户端->服务器|Result|[][CombatDto](#combatdto)|客户端提交本回合的行动数组（多个 `CombatDto`）。当 `CombatDto.Behavior == SwitchCard` 时，`select_card` 必须包含 `card_temp_id` 用于换牌。|
|Combat|服务器->客户端|Succeed|空|服务器在成功接收并存储一方的 `Result` 后返回，防止重复提交(但不是这阶段结束了,可能还要等对方)|




---

### <8> 卡牌效果结算
#### 1.CardCalc
执行阶段会按照
换牌结算-->行动结算-->法术牌结算
这样的顺序;具体的我建议用postman尝试请求一下看看流程,我现在带入战斗的卡是不会被删掉的,所以postman ws请求参数写死之后(就是带哪些牌什么的),一直可以请求调试的

|action_code|方向|Predicates|传参(action_data)|注释|
| :---: | :---: | :---: | :---: | :---: |
|CardCalc|服务器->客户端|Finish|空||
|DeployCard|服务器->客户端|Finish|map{"self":[][SelectCard](#selectcard),"opponent":[][SelectCard](#selectcard)}|换牌结算阶段结束之后,通知你双方执行了哪些换牌操作,用来做换牌演示动画|
|CardCalc|服务器->客户端|Notify|[AnimationDto](#animationdto)||
|Interrupt|服务器->客户端|Notify|[InterruptDto](#interruptdto)|发起中断(死亡，技能都有可能)，选完牌或者时间到了，就会恢复卡牌结算|
|Interrupt|服务器->客户端|Result|temp_id_list|从上一条Dto里的数组的tempid里选|
|Interrupt|服务器->客户端|Succeed|空|不管是时间到了，系统选的还是自己选的，都会有这条,系统选的话，可以根据返回的数组，告知用户随机了什么|



说明：`CardCalc` 是整个回合效果结算的核心阶段，主要流程参考服务端实现：

1. 接收来自 `Combat` 的 `CombatMap`（包含 `Winner` / `Loser` 两方的 `[]CombatDto`），按顺序处理每一条 `CombatDto`：
	- 如果 `CombatDto.select_card` 存在且 `select_card.where != SkillCard`，优先执行换牌逻辑（将手牌换入到对应槽）。
	- 否则按 `behavior` 执行攻击或技能：使用 `self_where` / `opponent_where` 定位场上卡并调用 `Attack` 或 `Skill`。
2. 法术牌结算：先处理胜者（Winner）的法术牌（`SkillCardBT`），再处理败者（Loser）的法术牌。
3. 中断处理：在结算过程中可能会收到/产生 `InterruptDto`（例如死亡触发需要选择），服务器会发出 `Interrupt.Notify` 要求客户端选择，客户端通过 `Interrupt.Result` 返回选择，服务器返回 `Interrupt.Succeed` 确认并继续结算。超时会由服务器随机/默认选择。
4. 完成标志：当所有 `CombatDto`、法术与中断均处理完毕，服务器发送 `CardCalc.Finish` 表示该轮卡牌效果结算结束(只会发一次finish,这个finish表示所有执行全完成了)，随后可能进入下一阶段（再次进入`SelectSkillCard`）。



---


## 三.具体字段说明

### StateWaitTime
|字段名|类型|注解|
| :---: | :---: |:---: |
|state_wait_time|int64|用上次封装好的函数解析出offset|
---

### JudgeData
|字段名|类型|注解|
| :---: | :---: |:---: |
|judge_data|0,1,2(三选一)||
---

### JudgeRes
|字段名|类型|注解|
| :---: | :---: |:---: |
|self|int|返回对方和自己选择了什么，有时因为超时，系统自动帮他选了，你要在画面上有所显示|
|opponent|int|
|is_win|0,1,2|0是平局，1是赢，-1是输|
---

### SelfCardDtoList
|字段名|类型|注解|
| :---: | :---: |:---: |
|空|[] [CardDto](#carddto)|没有字段名，直接就是一个数组|
---

### CardDto
|字段名|类型|注解|
| :---: | :---: |:---: |
|id|int||
|hp|float64||
|damage|float64||
|buff_list|[][BuffDto](#buffdto)|卡牌的BUFF列表|
|temp_id|int||
---

### BuffDto
|字段名|类型|注解|
| :---: | :---: |:---: |
|buff_id|int|BUFF的ID|
|buff_stacks|int|BUFF的层数|
---

### BtCardInfo
|字段名|类型|注解|
| :---: | :---: |:---: |
|self|[BtCardInfoMeta](#btcardinfometa)||
|opponent|[BtCardInfoMeta](#btcardinfometa)||
---

### BtCardInfoMeta
|字段名|类型|注解|
| :---: | :---: |:---: |
|skill_card_bt|[CardDto](#carddto)||
|parent_card_bt|[CardDto](#carddto)||
|child_card_bt|[CardDto](#carddto)||
---

### SelectCard
|字段名|类型|注解|
| :---: | :---: | :---: |
|where|[Where](#where)|选择的位置，ParentCard/ChildCard/SkillCard|
|card_id|int|
|card_temp_id|int|
---

### Where
|值|名称|注解|
| :---: | :---: | :---: |
|0|ParentCard|母卡位置|
|1|ChildCard|子卡位置|
|2|SkillCard|技能卡位置|
---


### SelectSkill
|字段名|类型|注解|
| :---: | :---: | :---: |
|state_wait_time|[StateWaitTime](#statewaittime)||
|where|[Where](#where)||
---

### DisCardDtoList
|字段名|类型|注解|
| :---: | :---: | :---: |
|空|[] [CardDto](#carddto)|没有字段名，直接就是一个数组|
---

### CombatDto
|字段名|类型|注解|
| :---: | :---: | :---: |
|behavior|[Behavior](#behavior)|行为|
|self_where|[Where](#where)|发动攻击或者技能的是自己的哪个|
|opponent_where|[Where](#where)|目标|
|select_card|[SelectCard](#selectcard)|可选：当需要从手牌出牌或换牌时,也就是behavior==2填写对应换牌信息,如果不是换牌,这个字段就穿一个空对象{}即可(如果behavior!=2传都不重要)|
---

### Behavior
|值|名称|注解|
| :---: | :---: | :---: |
|0|Attack|攻击|
|1|Skill|技能|
|2|SwitchCard|换牌：从手牌换入场上的卡（需要在对应 `CombatDto.SelectCard` 中填写 `card_temp_id`）|

---

### AnimationDto
|字段名|类型|注解|
| :---: | :---: | :---: |
|temp_id|int||
|animation_behavior|[AnimationBehavior](#animationbehavior)||
|bt_card_info|[BtCardInfo](#btcardinfo)||
---

### InterruptDto
|字段名|类型|注解|
| :---: | :---: | :---: |
|state_wait_time|int64||
|temp_id_list|[]int||
|select_num|int|选几个|
---

### ActiveChildCardDto
|字段名|类型|注解|
| :---: | :---: | :---: |
|temp_id_list|[]int|所选子卡的 `temp_id` 列表|
---
### ChildCardDto
|字段名|类型|注解|
| :---: | :---: |:---: |
|id|int||
|hp|float64||
|damage|float64||
|buff_list|[][BuffDto](#buffdto)|卡牌的BUFF列表|
|temp_id|int||
|child_state|[ChildState](#childstate)|子卡的状态|
---

### ChildState
|值|名称|注解|
| :---: | :---: | :---: |
|0|Active|已激活|
|1|NotActive|未激活|
|2|Died|已死亡|
|3|HasCatch|已被捕获|
---
### SelectWeatherDto
|字段名|类型|注解|
| :---: | :---: | :---: |
|weather|int|天气的枚举值（由服务器定义）|
|weather_list|[]int|服务器在 Query 时返回的候选天气列表（可选）|
---

### AnimationBehavior
|值|名称|注解|
| :---: | :---: | :---: |
|0|AnAttack||
|1|AnHurt||
|2|AnDeath|弃牌和死亡这两个通知是解藕开的，触发死亡效果函数的才会有死亡通知。而在手牌中没血了也会弃牌，但不会触发死亡的函数|
|3|AnSkill||
|4|AnDisCard|在手牌中也会触发的|
---


## 四.StatusCode

|StatusCode|名称|值|注解|
| :---: | :---: | :---: | :---: |
|ResponseSuccess|成功|0||
|ResponseDataNotFound|数据未找到|1||
|ResponseInternalServersError|内部服务器错误|2||
|ResponseInvalidReqParams|无效请求参数|3||
|ResponseInvalidToken|无效令牌|4||
|ResponseTokenExpired|令牌过期|5||
|ResponseIncorrectTokenFormat|令牌格式不正确|6||
|ResponseDuplicateDataEntry|重复数据条目|7||
|ResponseRequiredParamsMissing|缺少必需参数|8||
|ResponseDependentRecordsExist|存在相关记录|9||
|ResponseNotImplemented|未实现|10||
|ResponseIncorrectPassword|密码不正确|11||
|ResponseTokenMissing|令牌缺失|12||
|ResponseForbidden|禁止访问|13||
|ResponseRepeatRequest|重复请求|14||
|ResponseUnknownError|未知错误|15||
|ResponseTokenHasUpdate|令牌已更新|16||
|BattleInvalidTiming|战斗无效时序|17||
|BattleEffectStackOverflow|效果堆栈溢出|18||
|BattleCardCategoryError|战斗卡牌类别错误|19||
|BattleCardNotFound|战斗卡牌未找到|20||
|BattleNotInYourRound|不在你的回合|21||
|BattleHasCard|位置已有卡牌|22||
|BattleCardNumErr|卡牌数量错误|23||
