package tokenhandler

type TokenPostDto struct {
	Name     string `form:"name" json:"name" binding:"max=32,required"`
	Password string `form:"password" json:"password" binding:"max=256,required"`
}
