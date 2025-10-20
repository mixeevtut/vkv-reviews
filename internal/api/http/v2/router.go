package v2

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	oapigin "github.com/oapi-codegen/gin-middleware"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	apigen "vkvideo-reviews/internal/api/gen"
	"vkvideo-reviews/internal/api/http/v2/handlers"
)

type Router struct {
	*handlers.HTTPReviewHandler
	*handlers.HTTPMoviesHandler
	*handlers.HTTPUserHandler
}

func NewRouter(r *handlers.HTTPReviewHandler, m *handlers.HTTPMoviesHandler, u *handlers.HTTPUserHandler) *Router {
	return &Router{r, m, u}
}

func (router *Router) StartRouter() error {

	g := gin.Default()

	swagger, err := apigen.GetSwagger()
	if err != nil {
		return err
	}
	g.GET("/swagger.json", func(c *gin.Context) {
		b, err := json.Marshal(swagger)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.Data(http.StatusOK, "application/json", b)
	})

	g.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerfiles.Handler,
		ginSwagger.URL("/swagger.json"),
		ginSwagger.DefaultModelsExpandDepth(-1),
	))

	v2 := g.Group("/api/v2")
	v2.Use(oapigin.OapiRequestValidator(swagger))

	apigen.RegisterHandlers(v2, router)

	return g.Run(":8080")
}
