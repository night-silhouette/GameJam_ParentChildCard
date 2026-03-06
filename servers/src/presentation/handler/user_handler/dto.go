package user_handler

type UserSearchReqDto struct {
	ID   int    `form:"id"`
	Name string `form:"name"`
}
type UserPostDto struct {
	Name     string `form:"name" binding:"required,max=16"`
	Password string `form:"password" binding:"required,max=256"`
}
type UserDeleteReqDto struct {
	ID   int    `form:"id" binding:"required"`
	Name string `form:"name" binding:"required,max=16"`
}
