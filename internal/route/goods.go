package route

import (
	"goods-api/internal/controllers"

	"github.com/gin-gonic/gin"
)

func RouteGoods(router *gin.Engine, controller controllers.GoodsController) *gin.Engine {

	group := router.Group("/goods")
	{
		group.GET("", controller.GetAll)
		group.POST("", controller.Create)
		group.DELETE("", controller.Delete)
		group.PATCH("", controller.Update)
		group.PATCH("/reprioritize", controller.Reprioritize)
	}

	return router
}
