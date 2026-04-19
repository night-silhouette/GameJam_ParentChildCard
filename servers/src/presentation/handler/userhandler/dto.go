package userhandler

type UserSearchReqDto struct {
	Id   int    `form:"id" json:"id"`
	Name string `form:"name" json:"name" binding:"max=32"`
}
type UserPostDto struct {
	Name     string `form:"name" json:"name" binding:"required,max=16"`
	Password string `form:"password" json:"password" binding:"required,max=256"`
}
type UserDeleteReqDto struct {
	Id int `form:"id" json:"id"`
}

type UserPutReqDto struct {
	Id       int    `form:"id" json:"id" binding:"required,min=1"`
	Name     string `form:"name" json:"name" binding:"required,max=16"`
	Password string `form:"password" json:"password" binding:"required,max=256"`
}

type UserPatchDto struct {
	Name     string `form:"name" json:"name" binding:"max=16"`
	Password string `form:"password" json:"password" binding:"max=256"`
}

type UserVagueSearchReq struct {
	VagueName string `form:"vague_name" json:"vague_name" binding:"required,max=16"`
}
