# 0.101版本号接口文档

    这个部署的版本流程全部通了，但是没有技能
    这markdown都给你链接好了，字段，定义点进去查看。action_code映射表给你当目录用的，点击就可以索引到action_code相关的解释
    只写了成功请求的情况，各种请求失败的报错没有详细写出
    StatusCode在最后

## 一.action_code映射表

| action_code                              | action_name | value | 注解                                                                                                                  |
|:----------------------------------------:|:-----------:|:-----:|:------------------------------------------------------------------------------------------------------------------- |
| Fault                                    | 错误          | 0     | 接口里没有用到，只是默认值为0，所以如果传默认值就认为是错误                                                                                      |
| [CancelMatch](#1cancelmatch)             | 取消匹配        | 1     |                                                                                                                     |
| [GetSelfCardInHard](#3getselfcardinhard) | 获取自己的卡牌信息   | 2     |                                                                                                                     |
| GetOpponentCardInHard                    | 获取对手的卡牌信息   | 3     | 这个也只是调试用的，正常不能给用户获取对手手牌信息                                                                                           |
| [GetBtCardInfo](#4getbtcardinfo)         | 获取场上的战斗信息   | 4     |                                                                                                                     |
| [OverBattle](#2overbattle)               | 结束战斗        | 5     |                                                                                                                     |
| StartBattle                              | 开始战斗        | 6     |                                                                                                                     |
| [DeployCard](#1deploycard)               | 部署一张牌       | 7     | 这个actioncode我说实话设计的不好，但也不会改他了，现在他的作用就两个，一个是在状态机为选择技能牌的阶段的时候用，这个链接点进去也是讲这个的，然后还有就是游戏刚开始的看牌阶段和combat阶段，拖动上牌用这个。就这两个用处 |
| [Judge](#1judge)                         | 战斗回合判断      | 8     |                                                                                                                     |
| [MatchSuccess](#1matchsuccess)           | 匹配成功        | 9     |                                                                                                                     |
| [Combat](#1combat)                       | 执行战斗行动      | 11    |                                                                                                                     |
| [CardCalc](#1cardcalc)                   | 卡牌效果结算      | 12    |                                                                                                                     |
| Debug                                    | 测试          | 13    | 这个是我自己要查看一些运行时数据用的，没有你要用的正式接口                                                                                       |
| Interrupt                                | 中断选牌        | 14    |                                                                                                                     |
| [GetDisCard](#5getdiscard)               | 查看弃牌堆       | 15    |                                                                                                                     |
| [SkillCardCalc](#1skillcardcalc)         | 法术牌计算       | 16    |                                                                                                                     |

---

## 二.具体接口说明

### <1> 常用请求

#### 1.CancelMatch

| action_code | 方向       | Predicates | 传参(action_data) | 注释                         |
|:-----------:|:--------:|:----------:|:---------------:|:--------------------------:|
| CancelMatch | 客户端->服务器 | Notify     | 空               | 还在匹配的时候传，字面意思，取消匹配，ws会直接断掉 |

---

#### 2.OverBattle

| action_code | 方向       | Predicates | 传参(action_data) | 注释                         |
|:-----------:|:--------:|:----------:|:---------------:|:--------------------------:|
| OverBattle  | 客户端->服务器 | Notify     | 空               | 传了就会直接结束游戏，双方都会退出来，ws会直接断掉 |

---

#### 3.GetSelfCardInHard

| action_code       | 方向       | Predicates | 传参(action_data)                     | 注释  |
|:-----------------:|:--------:|:----------:|:-----------------------------------:|:---:|
| GetSelfCardInHard | 客户端->服务器 | Query      | 空                                   |     |
| GetSelfCardInHard | 服务器->客户端 | Result     | [SelfCardDtoList](#selfcarddtolist) |     |

---

#### 4.GetBtCardInfo

| action_code   | 方向       | Predicates | 传参(action_data)           | 注释  |
|:-------------:|:--------:|:----------:|:-------------------------:|:---:|
| GetBtCardInfo | 客户端->服务器 | Query      | 空                         |     |
| GetBtCardInfo | 服务器->客户端 | Result     | [BtCardInfo](#btcardinfo) |     |

---

#### 5.GetDisCard

| action_code | 方向       | Predicates | 传参(action_data)                 | 注释  |
|:-----------:|:--------:|:----------:|:-------------------------------:|:---:|
| GetDisCard  | 客户端->服务器 | Query      | 空                               |     |
| GetDisCard  | 服务器->客户端 | Result     | [DisCardDtoList](#discardolist) |     |

---

### <2> 看牌阶段

#### 1.MatchSuccess

| action_code  | 方向       | Predicates | 传参(action_data)                 | 注释                                                                                  |
|:------------:|:--------:|:----------:|:-------------------------------:|:-----------------------------------------------------------------------------------:|
| MatchSuccess | 服务器->客户端 | Notify     | [StateWaitTime](#statewaittime) | 收到这个信息表明匹配成功了，现在是特殊的看牌阶段,这个阶段不会快速结束，即使大家都上完牌了                                       |
| DeployCard   | 客户端->服务器 | Result     | [SelectCard](#selectcard)       | 这个阶段和combat阶段一样，可以拖动上牌，这个就是上牌的接口                                                    |
| DeployCard   | 服务器->客户端 | Success    | 空                               | 收到这个表示你上牌成功了，值得一提的是，这个阶段时间到了会自动上牌,自动上牌也会提醒success有这条消息，至于自动上了什么，可以通过GetBtCardInfo查看 |
| StartBattle  | 服务器->客户端 | Notify     | 空                               | 直到收到这个，表示这个阶段结束了（也可以直接去监听下一阶段的开始信号）                                                 |

---

### <3> 选择技能牌阶段

#### 1.DeployCard

| action_code | 方向       | Predicates | 传参(action_data)             | 注释                                                                                      |
|:-----------:|:--------:|:----------:|:---------------------------:|:---------------------------------------------------------------------------------------:|
| DeployCard  | 服务器->客户端 | Query      | [SelectSkill](#SelectSkill) | 这是开始进入选技能牌阶段的标志                                                                         |
| DeployCard  | 客户端->服务器 | Result     | [SelectCard](#selectcard)   | 传where是skillcard的SelectCard,特殊的一个用法是，如果你传cardTempId=-1，那表示你明确不上技能,这样双方都确认了之后就可快速的finish |
| DeployCard  | 服务器->客户端 | Success    | 空                           | 选成功了会给你这个,但是要等Finish，否则就是对方还在选，技能牌可上可不上，所以即使时间到了也不会系统自动选择                               |
| DeployCard  | 服务器->客户端 | Finish     | 空                           | 选结束了，如果双方都在倒计时结束之前结束，那他会提前结束                                                            |

---

### <4> Judge阶段

#### 1.Judge

| action_code      | 方向       | Predicates | 传参(action_data)                 | 注释                                                                |
|:----------------:|:--------:|:----------:|:-------------------------------:|:-----------------------------------------------------------------:|
| Judge            | 服务器->客户端 | Query      | [StateWaitTime](#statewaittime) | 开始回合判断，向客户端询问                                                     |
| Judge            | 客户端->服务器 | Result     | [JudgeData](#JudgeData)         | 你要传的，剪刀石头布                                                        |
| Judge            | 服务器->客户端 | Success    | 空                               | 传完上一个之后，就会返回你条这个，但是你那边要继续显示对方还在选择，直到收到Finish                      |
| Judge            | 服务器->客户端 | Finish     | [JudgeRes](#judgeres)           | 时间到了或者双方都选完了.就会返回你这个,如果是平局，会重新进入Judge状态，又会发你query，有胜负就会进入Combat阶段 |
| AnimationPlayEnd | 客户端->服务器 | Notify     | 空                               | 每次收到finish之后，剪刀石头布的动画结束之后，传我一个这个，后端状态机才会恢复继续运作                    |

---

### <5> Combat阶段

#### 1.Combat

| action_code      | 方向       | Predicates   | 传参(action_data)           | 注释                                                                                                                                                                                                                                                     |
|:----------------:|:--------:|:------------:|:-------------------------:|:------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------:|
| Combat           | 服务器->客户端 | Query\Notify | 空                         | 如果是回合判断的赢者，那他是Query,反之是Notify，输还是赢你那边要打个标记，后续的卡牌效果结算也会发给输的人，毕竟要播放动画                                                                                                                                                                                    |
| Combat           | 客户端->服务器 | Finish       | 空                         | 确认结束回合                                                                                                                                                                                                                                                 |
| DeployCard       | 客户端->服务器 | Result       | [SelectCard](#selectcard) | 在这个阶段也可以拖动上牌，刚上的牌也有行动力                                                                                                                                                                                                                                 |
| Combat           | 客户端->服务器 | Result       | [CombatDto](#combatdto)   | 因为牌的槽位置有两个，玩家可以先后选择子牌或者母牌行动，发了这条之后，我后端的状态机就会进入结算状态，怎么结算见CardCalc，但结算完之后，你会收到CardCalc的Finish，即下面这条，你把动画播放完之后，传AnimationPlayEnd                                                                                                                          |
| CardCalc         | 服务器->客户端 | Finish       | 空                         | 一系列卡牌效果结算的结束标识                                                                                                                                                                                                                                         |
| AnimationPlayEnd | 客户端->服务器 | Notify       | 空                         | 传完这个之后，后端状态机会继续运转，你会再次收到Combat.Query,如果你还有牌没有行动过，那你可以再次行动，然后重复上述过程。如果时间到了，有牌上了却没有行动，那么系统会自动的让他发动攻击,默认优先攻击对方母牌。值得一提的是，即使你没上子牌（或者没上母牌）那在一次行动之后，也会回到Combat.Query,直到回合倒计时结束，或者用户手动结束回合，当然如果你想提升用户体验，就做一个检测，用户是否还有牌没行动过，没有的话，就你帮用户自动结束回合，也就是调用Combat.Finish |

---

### <6> 卡牌效果结算

#### 1.CardCalc

| action_code | 方向       | Predicates | 传参(action_data)               | 注释                                                 |
|:-----------:|:--------:|:----------:|:-----------------------------:|:--------------------------------------------------:|
| CardCalc    | 服务器->客户端 | Finish     | 空                             | 一系列卡牌效果结算的结束标识                                     |
| CardCalc    | 服务器->客户端 | Notify     | [AnimationDto](#animationdto) |                                                    |
| Interrupt   | 服务器->客户端 | Notify     | [InterruptDto](#interruptdto) | 发起中断(死亡，技能都有可能)，选完牌或者时间到了，就会恢复卡牌结算                 |
| Interrupt   | 服务器->客户端 | Result     | temp_id_list                  | 从上一条Dto里的数组的tempid里选                               |
| Interrupt   | 服务器->客户端 | Succeed    | 空                             | 不管是时间到了，系统选的还是自己选的，都会有这条,系统选的话，可以根据返回的数组，告知用户随机了什么 |

---

### <7> 法术牌结算阶段

    正常这个阶段也要等你发动画结束才会进行下去，但是法术牌具体还没写，只是留了下虚函数，所以就没要AnimationPlayEnd,后面肯定是要的，你把这个AnimationPlayEnd在框架里的位置留着就好了

#### 1.SkillCardCalc

| action_code   | 方向       | Predicates | 传参(action_data) | 注释                                                  |
|:-------------:|:--------:|:----------:|:---------------:|:---------------------------------------------------:|
| SkillCardCalc | 服务器->客户端 | Notify     | 空               | 这个actioncode就怎么一个用，主要就是通知你结算skillcard了，赢的人先结算       |
| CardCalc      | 服务器->客户端 | Finish     | 空               | 也是会出现一系列的效果结算，也是finish结尾，如果没有上法术卡，那什么也不会收到，包括finish |

---

## 三.具体字段说明

### StateWaitTime

| 字段名             | 类型    | 注解                 |
|:---------------:|:-----:|:------------------:|
| state_wait_time | int64 | 用上次封装好的函数解析出offset |

---

### JudgeData

| 字段名        | 类型         | 注解  |
|:----------:|:----------:|:---:|
| judge_data | 0,1,2(三选一) |     |

---

### JudgeRes

| 字段名      | 类型    | 注解                                      |
|:--------:|:-----:|:---------------------------------------:|
| self     | int   | 返回对方和自己选择了什么，有时因为超时，系统自动帮他选了，你要在画面上有所显示 |
| opponent | int   |                                         |
| is_win   | 0,1,2 | 0是平局，1是赢，-1是输                           |

---

### SelfCardDtoList

| 字段名 | 类型                     | 注解             |
|:---:|:----------------------:|:--------------:|
| 空   | [] [CardDto](#carddto) | 没有字段名，直接就是一个数组 |

---

### CardDto

| 字段名     | 类型      | 注解  |
|:-------:|:-------:|:---:|
| id      | int     |     |
| hp      | float64 |     |
| damage  | float64 |     |
| buff_id | int     |     |
| temp_id | int     |     |

---

### BtCardInfo

| 字段名      | 类型                                | 注解  |
|:--------:|:---------------------------------:|:---:|
| self     | [BtCardInfoMeta](#btcardinfometa) |     |
| opponent | [BtCardInfoMeta](#btcardinfometa) |     |

---

### BtCardInfoMeta

| 字段名            | 类型                  | 注解  |
|:--------------:|:-------------------:|:---:|
| skill_card_bt  | [CardDto](#carddto) |     |
| parent_card_bt | [CardDto](#carddto) |     |
| child_card_bt  | [CardDto](#carddto) |     |

---

### SelectCard

| 字段名          | 类型              | 注解                                   |
|:------------:|:---------------:|:------------------------------------:|
| where        | [Where](#where) | 选择的位置，ParentCard/ChildCard/SkillCard |
| card_id      | int             |                                      |
| card_temp_id | int             |                                      |

---

### Where

| 值   | 名称         | 注解    |
|:---:|:----------:|:-----:|
| 0   | ParentCard | 母卡位置  |
| 1   | ChildCard  | 子卡位置  |
| 2   | SkillCard  | 技能卡位置 |

---

### SelectSkill

| 字段名             | 类型                              | 注解  |
|:---------------:|:-------------------------------:|:---:|
| state_wait_time | [StateWaitTime](#statewaittime) |     |
| where           | [Where](#where)                 |     |

---

### DisCardDtoList

| 字段名 | 类型                     | 注解             |
|:---:|:----------------------:|:--------------:|
| 空   | [] [CardDto](#carddto) | 没有字段名，直接就是一个数组 |

---

### CombatDto

| 字段名            | 类型                    | 注解        |
|:--------------:|:---------------------:|:---------:|
| behavior       | [Behavior](#behavior) | 攻击行为      |
| self_where     | [Where](#where)       | 自己选择的卡牌位置 |
| opponent_where | [Where](#where)       | 对手选择的卡牌位置 |

---

### Behavior

| 值   | 名称     | 注解  |
|:---:|:------:|:---:|
| 0   | Attack | 攻击  |
| 1   | Skill  | 技能  |

---

### AnimationDto

| 字段名                | 类型                                      | 注解  |
|:------------------:|:---------------------------------------:|:---:|
| temp_id            | int                                     |     |
| animation_behavior | [AnimationBehavior](#animationbehavior) |     |
| bt_card_info       | [BtCardInfo](#btcardinfo)               |     |

---

### InterruptDto

| 字段名             | 类型    | 注解  |
|:---------------:|:-----:|:---:|
| state_wait_time | int64 |     |
| temp_id_list    | []int |     |
| select_num      | int   | 选几个 |

---

### AnimationBehavior

| 值   | 名称        | 注解                                                       |
|:---:|:---------:|:--------------------------------------------------------:|
| 0   | AnAttack  |                                                          |
| 1   | AnHurt    |                                                          |
| 2   | AnDeath   | 弃牌和死亡这两个通知是解藕开的，触发死亡效果函数的才会有死亡通知。而在手牌中没血了也会弃牌，但不会触发死亡的函数 |
| 3   | AnSkill   |                                                          |
| 4   | AnDisCard | 在手牌中也会触发的                                                |

---

## 四.StatusCode

    从0开始的枚举
    
    ResponseSuccess 
    ResponseDataNotFound
    ResponseInternalServersError
    ResponseInvalidReqParams
    ResponseInvalidToken
    ResponseTokenExpired
    ResponseIncorrectTokenFormat
    ResponseDuplicateDataEntry
    ResponseRequiredParamsMissing
    ResponseDependentRecordsExist
    ResponseNotImplemented
    ResponseIncorrectPassword
    ResponseTokenMissing
    ResponseForbidden
    ResponseRepeatRequest
    ResponseUnknownError
    ResponseTokenHasUpdate
    BattleInvalidTiming
    BattleEffectStackOverflow
    BattleCardCategoryError
    BattleCardNotFound
    BattleNotInYourRound
    BattleHasCard
    BattleCardNumErr-8
