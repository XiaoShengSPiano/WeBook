package middleware

import (
	"encoding/gob"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// loginMiddlewareBuilder 登录中间件构造器
type LoginMiddlewareBuilder struct {
	ignorePaths []string
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

func (l *LoginMiddlewareBuilder) IgnorePaths(path string) *LoginMiddlewareBuilder {
	l.ignorePaths = append(l.ignorePaths, path)
	return l
}

func (l *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	gob.Register(time.Now().UnixMilli())
	return func(ctx *gin.Context) {
		// 不需要登录校验的路由
		for _, path := range l.ignorePaths {
			if path == ctx.Request.URL.Path {
				return
			}
		}

		//path := ctx.Request.URL.Path
		//
		//if path == "/users/login" ||
		//	path == "/users/signup" {
		//	return
		//}

		sess := sessions.Default(ctx)
		id := sess.Get("userId")
		if id == nil {
			// 未登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		updateTime := sess.Get("update_time")
		sess.Set("userId", id)
		sess.Options(sessions.Options{
			MaxAge: 60,
		})
		now := time.Now().UnixMilli()
		// 首次登录，未设置刷新key
		if updateTime == nil {
			sess.Set("update_time", now)
			sess.Save()
			return
		}

		updateTimeVal, ok := updateTime.(int64)
		if !ok {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// 60 s
		if now-updateTimeVal > 60*1000 {
			sess.Set("update_time", now)
			sess.Save()
			return
		}
	}
}
