package ioc

import (
	"context"
	"github.com/asynccnu/bff/web/banner"
	"github.com/asynccnu/bff/web/calendar"
	"github.com/asynccnu/bff/web/card"
	"github.com/asynccnu/bff/web/class"
	"github.com/asynccnu/bff/web/department"
	"github.com/asynccnu/bff/web/elecprice"
	"github.com/asynccnu/bff/web/feed"
	"github.com/asynccnu/bff/web/feedback_help"
	"github.com/asynccnu/bff/web/grade"
	"github.com/asynccnu/bff/web/infoSum"
	"github.com/asynccnu/bff/web/middleware"
	"github.com/asynccnu/bff/web/static"
	"github.com/asynccnu/bff/web/tube"
	"github.com/asynccnu/bff/web/user"
	"github.com/asynccnu/bff/web/website"
	"github.com/gin-gonic/gin"
	"time"
)

// 逆天参数数量,依赖注入一堆服务
func InitGinServer(
	loggerMiddleware *middleware.LoggerMiddleware,
	loginMiddleware *middleware.LoginMiddleware,
	corsMiddleware *middleware.CorsMiddleware,
	tube *tube.TubeHandler,
	user *user.UserHandler,
	static *static.StaticHandler,
	banner *banner.BannerHandler,
	department *department.DepartmentHandler,
	website *website.WebsiteHandler,
	calendar *calendar.CalendarHandler,
	feed *feed.FeedHandler,
	elecprice *elecprice.ElecPriceHandler, //添加你的服务handler
	grade *grade.GradeHandler,
	class *class.ClassHandler,
	feedback *feedback_help.FeedbackHelpHandler,
	infoSum *infoSum.InfoSumHandler,
	card *card.CardHandler,
) *gin.Engine {
	//初始化一个gin引擎
	engine := gin.Default()
	//使用全局中间件
	engine.Use(
		corsMiddleware.MiddlewareFunc(),
		loggerMiddleware.MiddlewareFunc(),
	)

	//创建用户认证中间件
	authMiddleware := loginMiddleware.MiddlewareFunc()

	//注册一堆路由
	user.RegisterRoutes(engine, authMiddleware)
	static.RegisterRoutes(engine, authMiddleware)
	banner.RegisterRoutes(engine, authMiddleware)
	department.RegisterRoutes(engine, authMiddleware)
	website.RegisterRoutes(engine, authMiddleware)
	calendar.RegisterRoutes(engine, authMiddleware)
	feed.RegisterRoutes(engine, authMiddleware)
	elecprice.RegisterRoutes(engine, authMiddleware)
	class.RegisterRoutes(engine, authMiddleware)
	feedback.RegisterRoutes(engine, authMiddleware)
	infoSum.RegisterRoutes(engine, authMiddleware)
	grade.RegisterRoutes(engine, authMiddleware)
	card.RegisterRoute(engine, authMiddleware)
	tube.RegisterRoutes(engine, authMiddleware)

	//返回路由
	return engine
}

func timeout() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, ok := ctx.Request.Context().Deadline()
		if !ok {
			// 强制给一个超时，省得我前端调试等得不耐烦
			newCtx, cancel := context.WithTimeout(ctx.Request.Context(), time.Second*10)
			defer cancel()
			ctx.Request = ctx.Request.Clone(newCtx)
		}
		ctx.Next()
	}
}
