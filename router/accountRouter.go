package router

import (
	"github.com/gin-gonic/gin"
	account "sharingunittest/account/handler"
)

func AccountRoutes(router *gin.RouterGroup, handler *account.AccountHandler) *gin.RouterGroup {
	r := router.Group("")
	r.POST("/account", handler.Insert)
	r.GET("/accounts", handler.GetById)
	return r
}
