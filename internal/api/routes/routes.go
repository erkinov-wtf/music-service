package routes

import (
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"music-service/internal/api/handlers"
	"music-service/internal/api/routes/path"
)

func RegisterRoutes(router *Router,
	groupHandler *handlers.GroupHandler,
	songHandler *handlers.SongHandler,
) {
	// Swagger docs
	router.Engine().GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Engine().Group("/api/v1")
	{
		path.RegisterGroupRoutes(api, groupHandler)
		path.RegisterSongRoutes(api, songHandler)
	}
}
