package token_handler

type TokenPostDto struct {
	Name     string `form:"name" binding:"max=32,required"`
	Password string `form:"password" binding:"max=256,required"`
}
