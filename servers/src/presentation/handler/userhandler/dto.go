package userhandler

type UserSearchReqDto struct {
	ID   int    `form:"id"`
	Name string `form:"name" binding:"max=32"`
}
type UserPostDto struct {
	Name     string `form:"name" binding:"required,max=16"`
	Password string `form:"password" binding:"required,max=256"`
}
type UserDeleteReqDto struct {
	ID int `form:"id" binding:"min=1"`
}

type UserPutReqDto struct {
	Id       int    `form:"id" binding:"required,min=1"`
	Name     string `form:"name" binding:"required,max=16"`
	Password string `form:"password" binding:"required,max=256"`
}

type UserPatchDto struct {
	Name     string `form:"name" binding:"max=16"`
	Password string `form:"password" binding:"max=256"`
}
