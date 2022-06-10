package jwt

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var (
	ErrClaimsNotFound = errors.New("获取令牌身份信息出错")
)

type Jwt struct {
	Key []byte
	// 超时时长
	Timeout time.Duration
	// token对应的http header名称。
	headName string
	// token对应的http header名称中的方案别名，如Bearer。
	schemeName string

	// TokenFrom 调用其他服务生成的token
	TokenFrom     func(*gin.Context) (interface{}, error)
	TokenResponse func(*gin.Context, interface{})
	TokenError    func(*gin.Context, error)

	// Login 自定义jwt生成
	Login         func(*gin.Context) (map[string]interface{}, error)
	LoginResponse func(*gin.Context, string, string)
	// Error 登录业务错误，如密码不正确，参数不正确等。
	LoginError func(*gin.Context, error)

	// 获取签名
	Signature (*gin.Context)
	// 黑名单
	Blacklist func(*gin.Context, string) bool
	// AuthorizeError 授权错误回调。
	AuthorizeError func(*gin.Context, error)
}

func (j *Jwt) Header(headName, schemeName string) {
	j.headName, j.schemeName = headName, schemeName
}

// LoginHandle 通常作用域登录/注册，或其他相关验证合法用户的路由。
func (j *Jwt) LoginHandle(ctx *gin.Context) {
	// 如果g.Authorize等于nil，则任其崩溃。否则AuthorizeHandle函数就没有任何意义。
	maps, err := j.Login(ctx)
	if err != nil {
		j.LoginError(ctx, err)
		return
	}

	// 创建token
	claims := jwt.MapClaims{}
	t := time.Now()

	if j.Timeout != time.Duration(0) {
		// 令牌过期时间戳。
		claims["exp"] = t.Add(j.Timeout).Unix()
	}
	// 令牌颁发的时间戳。
	claims["iat"] = t.Unix()

	for k, v := range maps {
		claims[k] = v
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	var (
		tokenStr string
	)
	if len(j.Key) != 0 {
		tk, err := token.SignedString(j.Key)
		if err != nil {
			j.LoginError(ctx, err)
			return
		}
		tokenStr = tk
	} else {
		tk, err := token.SigningString()
		if err != nil {
			j.LoginError(ctx, err)
			return
		}
		tokenStr = tk
	}

	if j.LoginResponse != nil {
		for k, v := range maps {
			ctx.Set(k, v)
		}
		token, err := j.Parse(tokenStr)
		if err != nil {
			j.AuthorizeError(ctx, err)
			return
		}
		j.LoginResponse(ctx, tokenStr, token.Signature)
	}
}

var (
	ErrTokenKeyNil     = errors.New("无法获取Token请求头")
	ErrTokenRule       = errors.New("请求头无法分离scheme和值")
	ErrTokenHeadScheme = errors.New("请求头中的scheme不匹配")
	ErrTokenValueNil   = errors.New("Token值为空")

	// ErrSignatureIsInvalid 签名无效。
	ErrSignatureIsInvalid = errors.New("签名无效")
	ErrTokenIsExpired     = errors.New("Token已过期")
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
		if j.AuthorizeError != nil {
			j.AuthorizeError(ctx, ErrTokenKeyNil)
		}
		return
	}

	var auth string

	if j.schemeName != "" {
		strs := strings.Split(authorization, " ")
		if len(strs) != 2 {
			if j.AuthorizeError != nil {
				j.AuthorizeError(ctx, ErrTokenRule)
			}
			return
		}

		if strs[0] != j.schemeName {
			if j.AuthorizeError != nil {
				j.AuthorizeError(ctx, ErrTokenHeadScheme)
				return
			}
		}
		auth = strs[1]
	} else {
		auth = authorization
	}

	if auth == "" || auth == "null" {
		if j.AuthorizeError != nil {
			j.AuthorizeError(ctx, ErrTokenValueNil)
			return
		}
	}

	token, err := j.Parse(auth)
	if err != nil {
		switch err.Error() {
		case "signature is invalid":
			if j.AuthorizeError != nil {
				j.AuthorizeError(ctx, ErrSignatureIsInvalid)
			}
		case "Token is expired":
			if j.AuthorizeError != nil {
				j.AuthorizeError(ctx, ErrTokenIsExpired)
			}
		default:
			if j.AuthorizeError != nil {
				j.AuthorizeError(ctx, err)
			}
		}

	} else {
		if claims, ok := token.Claims.(MapClaims); ok && token.Valid {
			// 如果验证通过，则需要确定黑名单
			if j.Blacklist != nil {
				if j.Blacklist(ctx, token.Signature) {
					j.AuthorizeError(ctx, errors.New("非法请求"))
					return
				}
			}
			for k, v := range claims {
				ctx.Set(k, v)
			}
		}
	}
}

func (j *Jwt) Parse(auth string) (*jwt.Token, error) {
	// 解析并验证传入的token。
	token, err := jwt.Parse(auth, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return j.Key, nil
	})
	return token, err
}
