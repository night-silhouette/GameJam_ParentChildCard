package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var VU int = 9

var BaseUrl string = "http://127.0.0.1:10086/"
var globalClient = &http.Client{
	Timeout: 10 * time.Second,
	Transport: &http.Transport{
		MaxIdleConnsPerHost: 500, // 允许对同一个后端服务器保留多少个空闲连接
	},
}

func NewPayLoad(param map[string]string) []byte {
	payload, _ := json.Marshal(param)
	return payload
}

type Result struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func NewResult(data interface{}) *Result {
	res := &Result{}
	res.Data = data
	return res
}

func Resquest(url string, method string, payload []byte, result any) {
	apiUrl := BaseUrl + url
	req, _ := http.NewRequest(method, apiUrl, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	res, err := globalClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		log.Printf("服务器返回错误码: %d", res.StatusCode)
		return
	}

	json.NewDecoder(res.Body).Decode(result)

}

func Resgister(userID int, Res chan<- any) {
	param := NewPayLoad(map[string]string{
		"name":     fmt.Sprintf("go_tester_%d", userID),
		"password": "123456",
	})
	result := NewResult(struct{}{})
	Resquest("v1/user/", http.MethodPost, []byte(param), result)
	Res <- result
}

func getToken(userID int, Res chan<- any) {
	param := NewPayLoad(map[string]string{
		"name":     fmt.Sprintf("go_tester_%d", userID),
		"password": "123456",
	})

	result := NewResult(struct{}{})
	Resquest("v1/token/", http.MethodPost, param, &result)
	Res <- result.Data.(map[string]interface{})["token"]
}

func connectWS(token string) {
	// 1. 设置连接地址 (注意协议是 ws 或 wss)
	u := url.URL{Scheme: "ws", Host: "127.0.0.1:10086", Path: "/v1/ws/"}

	// 2. 如果需要通过 URL 传 token，可以加上参数
	q := u.Query()
	q.Set("token", token)
	u.RawQuery = q.Encode()

	// 3. 拨号连接
	// Dialer 可以设置超时、握手参数等
	dialer := websocket.DefaultDialer
	// 修改客户端 Dial 部分
	conn, resp, err := dialer.Dial(u.String(), nil)
	if err != nil {
		if resp != nil {
			log.Printf("握手失败，状态码: %d", resp.StatusCode)
		}
		log.Fatal("拨号失败:", err)
	}

	// 4. 保持连接并处理消息 (开启一个 Goroutine 读，主流程写)
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("读取失败:", err)
				return
			}

			fmt.Println(string(message))

			//Act := BattleDto.GetActionByWsResByte(message)
			//if Act.ActionCode == BattleDto.StartBattle {
			//	BattleDto.Send(conn, BattleDto.GetOpponentCardInHard, "")
			//}
		}
	}()
	<-done
}
func main() {
	var wg sync.WaitGroup
	PlayerChan := make(chan any)
	for i := 0; i < VU; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			getToken(i, PlayerChan)
		}()
	}
	go func() {
		wg.Wait()
		close(PlayerChan)
	}()

	var Wswg sync.WaitGroup
	for v := range PlayerChan {
		Wswg.Add(1)
		go func() {
			defer Wswg.Done()
			connectWS(v.(string))
		}()

	}
	Wswg.Wait()
}
