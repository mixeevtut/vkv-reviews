package v1

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "vkvideo-reviews/docs"
	"vkvideo-reviews/internal/api/http/v1/handlers"
)

func StartRouter(
	reviews *handlers.HTTPReviewHandler,
) error {

	g := gin.Default()

	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r := g.Group("/api/v1/reviews")
	{
		r.GET("/:id", reviews.GetById)
		r.POST("", reviews.Create)
		r.PATCH("/:id", reviews.Update)
		r.DELETE("/:id", reviews.Delete)
	}

	return g.Run(":8080")
}
