package web

import "github.com/gin-gonic/gin"

type handler interface {
	RegisterRoutes(s *gin.Engine, authMiddleware gin.HandlerFunc)
}

type Response struct {
	Msg  string      `json:"msg"`
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}
