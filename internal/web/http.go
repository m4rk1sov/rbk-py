package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func OK(c *gin.Context, v any) {
	c.JSON(http.StatusOK, v)
}

func Error(c *gin.Context, code int, err error) {
	c.JSON(code, gin.H{"error": err.Error()})
}

func BadRequest(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, gin.H{"error": msg})
}
