package auth

import "github.com/gin-gonic/gin"

func Not_req_from_simulator(c *gin.Context) (statusCode int, errStr string) {
	auth := c.Request.Header.Get("Authorization")
	if auth != "Basic c2ltdWxhdG9yOnN1cGVyX3NhZmUh" {
		statusCode = 403
		errStr = "You are not authorized to use this resource!"
		return statusCode, errStr
	}
	return
}
