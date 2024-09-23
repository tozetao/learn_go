package web

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

type JWTHandler struct {
	RefreshTokenKey []byte
	ShortTokenKey   []byte
	redis           redis.Cmdable
	rcExpiration    time.Duration
}

func NewJWTHandler(redis redis.Cmdable) *JWTHandler {
	return &JWTHandler{
		RefreshTokenKey: []byte("Ahbwtj74udvDKmgE2gyQL8gBU2saxF"),
		ShortTokenKey:   []byte("ihbwtj3dGZvDKmgE2gyQL8gBU2saIE"),
		redis:           redis,
		rcExpiration:    time.Hour * 7,
	}
}

// CheckSession 检测token的有效性
func (j *JWTHandler) CheckSession(c *gin.Context, ssid string) error {
	key := fmt.Sprintf("users:ssid:%s", ssid)
	cnt, err := j.redis.Exists(c, key).Result()
	if err != nil {
		return err
	}
	if cnt > 0 {
		return errors.New("invalid ssid")
	}
	return nil
}

func (j *JWTHandler) ClearSession(c *gin.Context) error {
	c.Header("x-refresh-token", "")
	c.Header("x-jwt-token", "")

	// 从context中取出UserClaims
	user := c.MustGet("user").(UserClaims)
	key := fmt.Sprintf("users:ssid:%s", user.SSid)

	return j.redis.Set(c, key, "", j.rcExpiration).Err()
}

func (j *JWTHandler) SetJWTToken(c *gin.Context, uid int64, ssid string, userAgent string) error {
	claims := UserClaims{
		SSid:      ssid,
		Uid:       uid,
		UserAgent: userAgent,
		// token过期时间, 1分钟有效期。
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(j.ShortTokenKey)
	if err != nil {
		return err
	}
	c.Header("x-jwt-token", tokenStr)
	return nil
}

func (j *JWTHandler) SetRefreshToken(c *gin.Context, uid int64, ssid string, userAgent string) error {
	claims := RefreshClaims{
		UserClaims{
			Uid:       uid,
			SSid:      ssid,
			UserAgent: userAgent,
			// refresh token过期时间, 7天有效期。
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.rcExpiration)),
			},
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(j.RefreshTokenKey)
	if err != nil {
		return err
	}
	c.Header("x-refresh-token", tokenStr)
	return nil
}

func (j *JWTHandler) SetLoginToken(c *gin.Context, uid int64, userAgent string) error {
	ssid := uuid.New().String()

	err := j.SetJWTToken(c, uid, ssid, userAgent)
	if err != nil {
		return err
	}

	err = j.SetRefreshToken(c, uid, ssid, userAgent)
	if err != nil {
		return err
	}
	return nil
}

func (j *JWTHandler) ExtractToken(c *gin.Context) (string, error) {
	// 1. 获取长token，
	authToken := c.GetHeader("Authorization")
	if authToken == "" {
		return "", errors.New("error token")
	}
	// 提取出token字符串
	segs := strings.Split(" ", authToken)
	if len(segs) != 2 {
		return "", errors.New("error token")
	}
	return segs[1], nil
}

type RefreshClaims struct {
	UserClaims
}

type UserClaims struct {
	jwt.RegisteredClaims
	Uid       int64  `json:"uid"`
	SSid      string `json:"ssid"`
	UserAgent string `json:"user_agent"`
}

/*

让token携带一个ssid。
如果用户退出（标记该ssid失效），在验证这个ssid是否有效。

1. 生成token时也生成ssid
2. 退出时标时ssid失效
3. 需要检查ssid的有效性

*/
