package main

import (
	"strings"
	"time"
	"webook/internal/repository"
	"webook/internal/repository/dao"
	"webook/internal/service"
	"webook/internal/web"
	"webook/internal/web/middleware"
	"webook/pkg/ginx/middlewares/ratelimit"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions/memstore"

	// "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"

	"github.com/gin-contrib/sessions"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/redis/go-redis/v9"
)

func main() {
	db := initDB()
	server := initWebServer()
	u := initUserHandler(db)
	u.RegisterRoutes(server)
	server.Run(":8080")
}

func initDB() *gorm.DB {
	dns := "root:root@tcp(localhost:13316)/webook"
	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}

	return db
}

func initWebServer() *gin.Engine {
	server := gin.Default()

	// 初始化redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	// 注册ratelimit中间件，限制每分钟10次请求
	server.Use(ratelimit.NewBuilder(redisClient, time.Minute, 100).Build())

	// 使用解决cors的gin中间件
	server.Use(cors.New(cors.Config{
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"x-jwt-token"}, // 允许客户端获取的响应头
		AllowCredentials: true,
		// 自定义origin
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "company.com")
		},
		MaxAge: 12 * time.Hour,
	}))

	/* 使用Gin的session中间件管理session */
	// 这是基于内存的实现，第一个参数是 authentication key，最好是 32 或者 64 位
	//// 第二个参数是 encryption key
	store := memstore.NewStore([]byte("moyn8y9abnd7q4zkq2m73yw8tu9j5ixm"),
		[]byte("o6jdlo2cb9f9pb6h46fjmllw481ldebj"))

	// 创建一个基于 Redis 的 Session 存储
	//store, err := redis.NewStore(10, "tcp", "localhost:6379", "", "",
	//	[]byte("mpATDD5scnkhTuudmYs2y8HsbwrcvCnD"),
	//	[]byte("te7SQNMfKh54ZKv2GZVX3UPGQZ4WpASX"))

	//if err != nil {
	//	panic(err)
	//}

	// 存储在浏览器 Cookie中的 Session ID 的名称（默认叫 ssid）
	server.Use(sessions.Sessions("mysession", store))
	//// 接入登录校验的中间件 (使用链式调用)
	//server.Use(middleware.NewLoginMiddlewareBuilder().
	//	IgnorePaths("/users/login").
	//	IgnorePaths("/users/signup").Build())

	// 使用jwt中间件进行登录校验
	server.Use(middleware.NewLoginJWTMiddlewareBuilder().
		IgnorePaths("/users/login").
		IgnorePaths("/users/signup").
		Build())

	return server
}

func initUserHandler(db *gorm.DB) *web.UserHandler {
	ud := dao.NewUserDAO(db)
	ur := repository.NewUserRepository(ud)
	us := service.NewUserService(ur)
	uhdl := web.NewUserHandler(us)
	return uhdl
}
