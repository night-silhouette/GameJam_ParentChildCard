package userhandler

type PostFriendshipsReq struct {
	Id int `form:"id" json:"id"`
}
type DeleteFriendshipsReq struct {
	Id int `form:"id" json:"id"`
}
type GetFriendshipsRes struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
