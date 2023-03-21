package routes

import (
	"github.com/gin-gonic/gin"
	"xmarket_gin/controller"
	"xmarket_gin/middleware"
)

func CollectRoute(r *gin.Engine) *gin.Engine {
	// 中间件
	r.Use(middleware.CORSMiddleware(), middleware.RecoveryMiddleware())
	// 用户的注册登录业务
	r.POST("/register", controller.Register)
	r.POST("/login", controller.Login)
	// 用户登陆后读取用户信息
	r.GET("/info", middleware.AuthMiddleware(), controller.Info)

	// 管理系统的用户列表业务
	VIPRoutes := r.Group("/VIP")
	VIPController := controller.NewVIPController()
	// 管理系统新建用户
	VIPRoutes.POST("", VIPController.Create)
	// 管理系统修改用户信息
	VIPRoutes.PUT("/:Telephone", VIPController.Update)
	// 管理系统修改会员权限信息
	VIPRoutes.PUT("/set/:Telephone", VIPController.UpdateVIP)
	// 会员系统充值更新会员期限
	VIPRoutes.PUT("/user/:Telephone", VIPController.User2VIP)
	// 管理系统删除用户信息
	VIPRoutes.DELETE("/:Telephone", VIPController.Delete)
	// 管理系统查询指定用户信息 by Telephone
	VIPRoutes.GET("/:Telephone", VIPController.Show)
	// 管理系统查询所有用户信息
	VIPRoutes.GET("/show", VIPController.ShowALL)
	// 管理系统查询所有代金券信息
	VIPRoutes.GET("/voucher/show", VIPController.ReadVouchers)
	// 管理系统使用代金券
	VIPRoutes.PUT("voucher/:ID", VIPController.UpdateVouchers)
	// 会员系统兑换代金券
	VIPRoutes.PUT("/point", VIPController.UpdatePoint)
	// 会员系统查看代金券
	VIPRoutes.GET("/voucher/:Telephone", VIPController.ReadVouchersByTelephone)

	// 管理系统的设置商品业务
	GoodsRoutes := r.Group("/Goods")
	GoodsController := controller.NewGoodsController()
	GoodsRoutes.POST("", GoodsController.Create)
	GoodsRoutes.PUT("/:ID", GoodsController.Update)
	GoodsRoutes.DELETE("/:ID", GoodsController.Delete)
	GoodsRoutes.GET("/:ID", GoodsController.Show)
	GoodsRoutes.GET("/show", GoodsController.ShowAll)

	// 管理系统的订单列表业务
	OrderRoutes := r.Group("/Order")
	OrderController := controller.NewOrderController()
	OrderRoutes.POST("", OrderController.Create)
	OrderRoutes.PUT("/:OrderID", OrderController.Update)
	OrderRoutes.DELETE("/:OrderID", OrderController.Delete)
	OrderRoutes.GET("/:OrderID", OrderController.Show)
	OrderRoutes.GET("/show", OrderController.ShowAll)
	// 会员消费类型饼状图数据
	OrderRoutes.GET("/pie", OrderController.ReadPie)
	// 会员消费类型柱状图数据
	OrderRoutes.GET("/bar", OrderController.ReadBar)
	return r
}
