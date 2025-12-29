//go:build wireinject

package main

import (
	"webook/internal/repository"
	"webook/internal/repository/cache"
	"webook/internal/repository/dao"
	"webook/internal/service"
	"webook/internal/web"
	"webook/ioc"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		// 初始化第三方依赖
		ioc.InitDB, ioc.InitRedis,

		// 初始化DAO
		dao.NewUserDAO,

		// 初始化缓存
		cache.NewUserCache,
		cache.NewCodeCache,

		// 初始化Repository
		repository.NewUserRepository,
		repository.NewCodeRepository,

		// 初始化Service
		service.NewUserService,
		service.NewCodeService,

		// 初始化Handler
		ioc.InitSMSService,
		web.NewUserHandler,

		ioc.InitWebServer,
		ioc.InitMiddlewares,
	)

	return new(gin.Engine)
}
