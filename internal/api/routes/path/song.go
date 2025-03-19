package path

import (
	"github.com/gin-gonic/gin"
	"music-service/internal/api/handlers"
)

func RegisterSongRoutes(r *gin.RouterGroup, handler *handlers.SongHandler) {
	groups := r.Group("/songs")
	{
		groups.POST("", handler.CreateSong)
		groups.GET("", handler.GetAllSongs)
		groups.GET("/:id", handler.GetSong)
		groups.PUT("/:id", handler.UpdateSong)
		groups.DELETE("/:id", handler.DeleteSong)
	}
}
