package web

import (
	"fmt"
	"net/http"
	"time"
	"webook/internal/domain"
	"webook/internal/service"

	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	jwt "github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	jwt.RegisteredClaims
	Uid       int64 // 需要放入token中的数据
	UserAgent string
}

// 定义所有与用户相关的路由(Handler)
type UserHandler struct {
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
	svc         *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	// 正则表达式校验请求用户注册信息
	const (
		emailRegexPattern = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
		// 和上面比起来，用 ` 看起来就比较清爽
		passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,72}$`
	)

	return &UserHandler{
		emailExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		svc:         svc,
	}
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	// 分组路由
	ug := server.Group("/users")
	ug.POST("/signup", u.SignUp)
	// ug.POST("/login", u.Login)
	ug.POST("/login", u.LoginJWT)
	ug.POST("/edit", u.Edit)
	// ug.GET("/profile", u.Profile)
	ug.GET("/profile", u.ProfileJWT)
}

// 注册路由处理逻辑
func (u *UserHandler) SignUp(ctx *gin.Context) {
	// 注册请求体
	type SignUpReq struct {
		Email           string `json:"email"`
		ConfirmPassword string `json:"confirmPassword"`
		Password        string `json:"password"`
	}

	var req SignUpReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	// 邮箱格式校验
	isMatch, err := u.emailExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误......")
		return
	}

	if !isMatch {
		ctx.String(http.StatusOK, "邮箱格式不对......")
		return
	}

	// 密码校验
	if req.ConfirmPassword != req.Password {
		ctx.String(http.StatusOK, "两次输入的密码不一致......")
		return
	}

	ok, err := u.passwordExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误......")
		return
	}

	if !ok {
		ctx.String(http.StatusOK, "密码必须大于8位，包含特殊数字与字符......")
		return
	}

	err = u.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})

	if err == service.ErrUserDuplicateEmail {
		ctx.String(http.StatusOK, "邮箱冲突......")
		return
	}

	if err != nil {
		ctx.String(http.StatusOK, "系统异常......")
		return
	}

	ctx.String(http.StatusOK, "注册成功......")

	fmt.Println("%v", req)
}

// 使用JWT实现用户登录
func (u *UserHandler) LoginJWT(ctx *gin.Context) {
	type loginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req loginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	user, err := u.svc.Login(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})

	if err == service.ErrInvalidUserOrPassword {
		ctx.String(http.StatusOK, "用户名或密码错误......")
		return
	}

	if err != nil {
		ctx.String(http.StatusOK, "系统错误......")
		return
	}

	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		},
		Uid:       user.Id,
		UserAgent: ctx.Request.UserAgent(),
	}

	// 使用JWT设置登录态
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	key := []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0")
	tokenStr, err := token.SignedString(key)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误......")
		return
	}

	// 将生成的token写入到响应头
	ctx.Header("x-jwt-token", tokenStr)
	fmt.Println(user)

	ctx.String(http.StatusOK, "登录成功......")
}

// 用户登录处理逻辑
func (u *UserHandler) Login(ctx *gin.Context) {
	type loginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req loginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	user, err := u.svc.Login(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})

	if err == service.ErrInvalidUserOrPassword {
		ctx.String(http.StatusOK, "用户名或密码错误......")
		return
	}

	if err != nil {
		ctx.String(http.StatusOK, "系统错误......")
		return
	}

	// 获取session
	sess := sessions.Default(ctx)
	// 存储数据
	sess.Set("userId", user.Id)
	sess.Options(sessions.Options{
		//Secure: true,
		//HttpOnly: true,
		MaxAge: 60,
	})
	sess.Save()

	ctx.String(http.StatusOK, "登录成功......")
}

func (u *UserHandler) Logout(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	// todo
	// sess.Clear()
	sess.Options(sessions.Options{
		MaxAge: -1, // 立即过期
	})
	err := sess.Save()
	if err != nil {
		ctx.String(http.StatusOK, "退出登录失败......")
		return
	}

	ctx.String(http.StatusOK, "退出登录成功......")
}

// 用户信息编辑处理逻辑
func (u *UserHandler) Edit(ctx *gin.Context) {

}

// 获取用户配置处理逻辑
func (u *UserHandler) Profile(ctx *gin.Context) {
	ctx.String(http.StatusOK, "用户配置页面......")
}

func (u *UserHandler) ProfileJWT(ctx *gin.Context) {
	c, ok := ctx.Get("claims")
	if !ok {
		ctx.String(http.StatusOK, "系统错误......")
		return
	}

	claims, ok := c.(*UserClaims)
	if !ok {
		ctx.String(http.StatusOK, "系统错误......")
		return
	}

	fmt.Println("claims.Uid:", claims.Uid)
	ctx.String(http.StatusOK, "profile......")
}
