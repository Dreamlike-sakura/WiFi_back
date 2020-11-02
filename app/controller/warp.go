package controller

import "github.com/gin-gonic/gin"

type warp struct {
}

func (*warp) SuccessWarp(data interface{}) *gin.H {
	return &gin.H{
		"status":  "success",
		"message":  nil,
		"data":     data,
	}
}

func (*warp) FailWarp(errMsg interface{}) *gin.H {
	return &gin.H{
		"status":  "error",
		"message":  errMsg,
		"data":     nil,
	}
}

