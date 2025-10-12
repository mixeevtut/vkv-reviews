package v1

import (
	_ "vkvideo-reviews/docs"
	"vkvideo-reviews/internal/api/http/v1/handlers"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func StartRouter(
	reviews *handlers.HTTPReviewHandler,
	movies *handlers.HTTPMoviesHandler,
	user *handlers.HTTPUserHandler,
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

	// m := g.Group("/api/v1/movies")
	// {
	// 	m.GET()
	// }

	u := g.Group("/api/v1/user")
	{
		u.POST("/create", user.CreateUser)
	}

	return g.Run(":8080")
}
