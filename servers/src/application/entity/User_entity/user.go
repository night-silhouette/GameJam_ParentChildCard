package User_entity

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"` //hash过的
	Is_admin bool   `json:"is_admin"`
}
