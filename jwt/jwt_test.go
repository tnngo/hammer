package jwt

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestJwt_LoginHandle(t *testing.T) {
	j := &Jwt{
		Key: []byte("test"),
	}

	j.Login = func(ctx *gin.Context) (map[string]interface{}, error) {
		maps := make(map[string]interface{})
		maps["sub"] = 1
		return maps, nil
	}

	j.LoginHandle(nil)
}
