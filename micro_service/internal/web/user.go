package web

import (
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"learn_go/webook/internal/domain"
	"learn_go/webook/internal/service"
	"log"
	"net/http"
	"time"
	"unicode/utf8"
)

type UserHandler struct {
	passwordExp *regexp.Regexp
	emailExp    *regexp.Regexp
	birthdayExp *regexp.Regexp
	svc         service.UserService
	codeSvc     service.CodeService
	jwtHandler  *JWTHandler
}

func NewUserHandler(svc service.UserService, codeSvc service.CodeService, jwtHandler *JWTHandler) *UserHandler {
	const (
		passwordRegexpPattern = `^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d]{5,}$`
		emailRegexpPattern    = `^[A-Za-z0-9]+([_\.][A-Za-z0-9]+)*@([A-Za-z0-9\-]+\.)+[A-Za-z]{2,6}$`
		birthdayRegexpPattern = `^\d{4}-\d{2}-\d{2}$`
	)

	passwordExp := regexp.MustCompile(passwordRegexpPattern, regexp.None)
	emailExp := regexp.MustCompile(emailRegexpPattern, regexp.None)

	birthdayExp := regexp.MustCompile(birthdayRegexpPattern, regexp.None)

	return &UserHandler{
		passwordExp: passwordExp,
		emailExp:    emailExp,
		birthdayExp: birthdayExp,
		svc:         svc,
		codeSvc:     codeSvc,
		jwtHandler:  jwtHandler,
	}
}

func (u *UserHandler) SignUp(c *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirm_password"`
	}

	var req SignUpReq
	if err := c.Bind(&req); err != nil {
		c.String(http.StatusOK, "系统错误")
		return
	}

	ok, err := u.passwordExp.MatchString(req.Password)
	if err != nil {
		c.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		c.String(http.StatusOK, "你的密码强度不足.")
		return
	}

	// 取消email的验证
	ok, err = u.emailExp.MatchString(req.Email)
	if err != nil {
		log.Println(3, err)
		c.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		log.Println(4)
		c.String(http.StatusOK, "你的邮箱格式不对")
		return
	}

	err = u.svc.SignUp(c, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err == service.ErrDuplicateUser {
		c.String(http.StatusOK, "邮箱重复了")
		return
	}
	if err != nil {
		c.String(http.StatusOK, "系统错误")
		return
	}

	c.String(http.StatusOK, "注册成功")
}

func (u *UserHandler) Login(c *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := c.Bind(&req); err != nil {
		return
	}
	user, err := u.svc.Login(c, req.Email, req.Password)
	if err == service.ErrInvalidEmailOrPassword {
		c.String(http.StatusOK, "错误的用户名或密码")
		return
	}
	if err != nil {
		c.String(http.StatusOK, "系统错误")
		return
	}

	sess := sessions.Default(c)
	sess.Set("id", user.ID)
	sess.Options(sessions.Options{
		MaxAge: 60,
		Path:   "/",
	})
	sess.Save()

	c.String(http.StatusOK, "success")
}

func (u *UserHandler) LoginJWT(c *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := c.Bind(&req); err != nil {
		return
	}
	user, err := u.svc.Login(c, req.Email, req.Password)
	if err == service.ErrInvalidEmailOrPassword {
		c.String(http.StatusOK, "错误的用户名或密码")
		return
	}
	if err != nil {
		c.String(http.StatusOK, "系统错误")
		return
	}

	// 设置token
	userAgent := c.GetHeader("User-Agent")
	err = u.jwtHandler.SetLoginToken(c, user.ID, userAgent)
	if err != nil {
		c.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	c.String(http.StatusOK, "success")
}

func (u *UserHandler) LoginSMS(c *gin.Context) {
	// 手机号、验证码
	type LoginReq struct {
		Code  string `json:"code"`
		Phone string `json:"phone"`
	}
	req := &LoginReq{}
	if err := c.Bind(req); err != nil {
		return
	}
	// 手机号、验证码验证
	if req.Phone == "" {
		c.String(http.StatusOK, "请输入手机号")
		return
	}

	ok, err := u.codeSvc.Verify(c, "login", req.Phone, req.Code)
	if err != nil {
		log.Printf("code error: %v", err)
		c.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		c.String(http.StatusOK, "验证码错误")
		return
	}

	user, err := u.svc.FindOrCreate(c, req.Phone)
	if err != nil {
		log.Printf("find or create failed: %v", err)
		c.String(http.StatusOK, "内部错误")
		return
	}

	userAgent := c.GetHeader("User-Agent")
	err = u.jwtHandler.SetLoginToken(c, user.ID, userAgent)
	if err != nil {
		c.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	c.String(http.StatusOK, "success")
}

func (u *UserHandler) Edit(c *gin.Context) {
	type EditReq struct {
		Birthday string `json:"birthday"`
		AboutMe  string `json:"profile"`
		Nickname string `json:"nickname"`
	}
	var req EditReq
	if err := c.Bind(&req); err != nil {
		return
	}

	birthday, err := time.Parse(time.DateOnly, req.Birthday)
	if err != nil {
		c.String(http.StatusOK, "生日格式不正确")
		return
	}

	// nickname, 长度
	l := utf8.RuneCountInString(req.Nickname)
	if l <= 0 || l > 30 {
		c.String(http.StatusOK, "昵称的长度要大于0或者小于30个字符")
		return
	}
	//l = utf8.RuneCountInString(req.AboutMe)
	//if l <= 0 || l > 255 {
	//	c.String(http.StatusOK, "简介的长度要大于0或者小于255个字符")
	//	return
	//}

	sess := sessions.Default(c)
	idInterface := sess.Get("id")
	id, ok := idInterface.(int64)
	if !ok {
		c.String(http.StatusOK, "系统错误1")
		return
	}

	err = u.svc.UpdateNonSensitiveInfo(c, domain.User{
		ID:       id,
		Birthday: birthday,
		Nickname: req.Nickname,
		AboutMe:  req.AboutMe,
	})
	if err != nil {
		c.String(http.StatusOK, "更新失败")
		return
	}
	c.String(http.StatusOK, "success")
}

func (u *UserHandler) Profile(c *gin.Context) {
	// session实现
	//sess := sessions.Default(c)
	//idInterface := sess.Get("id")
	//id, ok := idInterface.(int64)
	//if !ok {
	//	c.String(http.StatusOK, "系统错误")
	//	return
	//}

	claimsVal, _ := c.Get("claims")
	claims, ok := claimsVal.(*UserClaims)
	if !ok {
		// 必然存在claims的情况出错了，这种异常情况需要记录日志。
		c.String(http.StatusOK, "系统错误")
		return
	}

	user, err := u.svc.Profile(c, claims.Uid)
	if err == service.ErrUserNotFound {
		c.String(http.StatusOK, "找不到该用户")
		return
	}
	if err != nil {
		c.String(http.StatusOK, "系统错误")
		return
	}
	response := map[string]any{
		"id":         user.ID,
		"nickname":   user.Nickname,
		"birthday":   user.Birthday.Format(time.DateOnly),
		"about_me":   user.AboutMe,
		"created_at": user.Ctime.Format(time.DateOnly),
	}
	c.JSON(http.StatusOK, response)
}

func (u *UserHandler) RefreshToken(c *gin.Context) {
	authToken, err := u.jwtHandler.ExtractToken(c)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var claims RefreshClaims
	token, err := jwt.ParseWithClaims(authToken, &claims, func(token *jwt.Token) (interface{}, error) {
		return u.jwtHandler.RefreshTokenKey, nil
	})
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if token == nil || !token.Valid {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	userAgent := c.GetHeader("user-agent")
	if userAgent != "" && userAgent != claims.UserAgent {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	err = u.jwtHandler.CheckSession(c, claims.SSid)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// 返回新的jwt token
	err = u.jwtHandler.SetJWTToken(c, claims.Uid, claims.SSid, userAgent)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.JSON(http.StatusOK, Result{Code: 0, Msg: "success"})
}

func (u *UserHandler) Logout(c *gin.Context) {
	err := u.jwtHandler.ClearSession(c)
	if err != nil {
		c.JSON(http.StatusOK, Result{Code: 5, Msg: "internal server error."})
		return
	}
	c.JSON(http.StatusOK, Result{Code: 0, Msg: "success"})
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	rg := server.Group("/users")

	rg.POST("/signup", u.SignUp)
	rg.POST("/login", u.LoginJWT)
	rg.POST("/login_sms", u.LoginSMS)
	rg.POST("/edit", u.Edit)
	rg.GET("/profile", u.Profile)
	rg.POST("/logout", u.Logout)
}
