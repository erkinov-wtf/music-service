package path

import (
	"github.com/gin-gonic/gin"
	"music-service/internal/api/handlers"
)

func RegisterSongRoutes(r *gin.RouterGroup, handler *handlers.SongHandler) {
	songs := r.Group("/songs")
	{
		songs.POST("", handler.CreateSong)
		songs.GET("", handler.GetAllSongs)
		songs.GET("/:id", handler.GetSong)
		songs.GET("/:id/verses", handler.GetSongVerses)
		songs.PUT("/:id", handler.UpdateSong)
		songs.DELETE("/:id", handler.DeleteSong)
	}
}
