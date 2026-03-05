package entity

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"` //hash过的
}
