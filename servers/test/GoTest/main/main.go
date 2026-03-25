package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

var VU int = 10

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
	for v := range PlayerChan {
		fmt.Println(v)
	}
}
