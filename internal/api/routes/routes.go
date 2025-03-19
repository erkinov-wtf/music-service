package routes

import (
	"music-service/internal/api/handlers"
	"music-service/internal/api/routes/path"
)

func RegisterRoutes(router *Router,
	groupHandler *handlers.GroupHandler,
	songHandler *handlers.SongHandler,
) {
	api := router.Engine().Group("/api/v1")
	{
		path.RegisterGroupRoutes(api, groupHandler)
		path.RegisterSongRoutes(api, songHandler)
	}
}
