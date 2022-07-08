package jwt

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/tnngo/hammer/sts/herr"
	"github.com/tnngo/lad"
)

var (
	whitelist sync.Map
)

func SetWhitelist(userId int, userKey string) {
	whitelist.Store(userId, userKey)
}

// JWT内部用户结构
type User struct {
	Id        int
	Key       string
	ExpiredAt time.Time
}

type Jwt struct {
	// token对应的http header名称。
	headName string
	// token对应的http header名称中的方案别名，如Bearer。
	schemeName string

	// TokenFrom 调用其他服务生成的token
	TokenFrom     func(*gin.Context) (interface{}, error)
	TokenResponse func(*gin.Context, interface{})
	TokenError    func(*gin.Context, error)

	// Login 自定义jwt生成
	Login         func(*gin.Context) (*User, error)
	LoginResponse func(*gin.Context, string)
	// Error 登录业务错误，如密码不正确，参数不正确等。
	LoginError func(*gin.Context, error)

	// 获取签名
	Signature (*gin.Context)
	// AuthorizeError 授权错误回调。
	AuthorizeError func(*gin.Context, error)
}

func (j *Jwt) Header(headName, schemeName string) {
	j.headName, j.schemeName = headName, schemeName
}

// LoginHandle 通常作用域登录/注册，或其他相关验证合法用户的路由。
func (j *Jwt) LoginHandle(ctx *gin.Context) {
	// key := ctx.
	// 如果g.Authorize等于nil，则任其崩溃。否则AuthorizeHandle函数就没有任何意义。
	u, err := j.Login(ctx)
	if err != nil {
		j.LoginError(ctx, err)
		return
	}

	// 创建token
	claims := jwt.MapClaims{}
	t := time.Now()

	if !u.ExpiredAt.IsZero() {
		claims["exp"] = int(u.ExpiredAt.Unix())
	}
	// 令牌颁发的时间戳。
	claims["iat"] = t.Unix()
	claims["sub"] = u.Id
	claims["key"] = u.Key

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	var (
		tokenStr string
	)
	tk, err := token.SignedString([]byte(u.Key))
	if err != nil {
		lad.L().Error(err.Error())
		j.LoginError(ctx, herr.Unauthenticated())
		return
	}
	tokenStr = tk

	j.LoginResponse(ctx, tokenStr)
}

var (
	ErrTokenIllegal = errors.New("非法请求")

	// ErrSignatureIsInvalid 签名无效。
	ErrTokenInvalid   = errors.New("签名无效")
	ErrTokenIsExpired = errors.New("账号已过期")
)

type MapClaims = jwt.MapClaims

// AuthorizedHandle 授权处理，指定路由参与签名验证。
// 比如，gin.Get("/user/list", func)
// 或者 group := gin.Group("/api", j.AuthorizeHandle)，
//
func (j *Jwt) AuthorizeHandle(ctx *gin.Context) {
	if j.headName == "" {
		j.headName = "Authorization"
	}
	if j.schemeName == "" {
		j.schemeName = "Bearer"
	}
	authorization := ctx.Request.Header.Get(j.headName)
	if authorization == "" {
		lad.L().Error("未传请求头参数Authorization")
		if j.AuthorizeError != nil {
			j.AuthorizeError(ctx, ErrTokenIllegal)
		}
		return
	}

	var auth string

	if j.schemeName != "" {
		strs := strings.Split(authorization, " ")
		if len(strs) != 2 {
			lad.L().Error("请求头Authorization参数值错误: " + authorization)
			if j.AuthorizeError != nil {
				j.AuthorizeError(ctx, ErrTokenIllegal)
			}
			return
		}

		if strs[0] != j.schemeName {
			lad.L().Error("请求头Authorization参数中的scheme错误: " + strs[0])
			if j.AuthorizeError != nil {
				j.AuthorizeError(ctx, ErrTokenIllegal)
				return
			}
		}
		auth = strs[1]
	} else {
		auth = authorization
	}

	if auth == "" || auth == "null" {
		lad.L().Error("请求头Authorization参数中的token为空")
		if j.AuthorizeError != nil {
			j.AuthorizeError(ctx, ErrTokenIllegal)
			return
		}
	}

	token, err := j.Parse(auth, ctx.GetString("key"))
	if err != nil {
		switch err.Error() {
		case "signature is invalid":
			if j.AuthorizeError != nil {
				j.AuthorizeError(ctx, ErrTokenIllegal)
			}
		case "Token is expired":
			if j.AuthorizeError != nil {
				j.AuthorizeError(ctx, ErrTokenIsExpired)
			}
		default:
			if j.AuthorizeError != nil {
				lad.L().Error("jwt解析错误: " + err.Error())
				j.AuthorizeError(ctx, herr.Unauthenticated())
			}
		}

	} else {
		if claims, ok := token.Claims.(MapClaims); ok && token.Valid {
			for k, v := range claims {
				ctx.Set(k, v)
			}
		}
	}
}

func (j *Jwt) Parse(auth, key string) (*jwt.Token, error) {
	// 解析并验证传入的token。
	token, err := jwt.Parse(auth, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		var (
			sub int
			key string
		)
		if claims, ok := token.Claims.(MapClaims); ok {
			if v, ok := claims["key"]; ok {
				key = v.(string)
			}
			if v, ok := claims["sub"]; ok {
				sub = int(v.(float64))
			}
		} else {
			return nil, ErrTokenIllegal
		}

		if v, ok := whitelist.Load(sub); ok {
			if v.(string) != key {
				return nil, ErrTokenIllegal
			}
		}

		return []byte(key), nil

	})
	return token, err
}
