package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": data,
	})
}

func fail(c *gin.Context, httpStatus, code int, msg string) {
	c.JSON(httpStatus, gin.H{
		"code": code,
		"msg":  msg,
		"data": nil,
	})
}
