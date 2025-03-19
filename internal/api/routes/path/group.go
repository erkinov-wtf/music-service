package path

import (
	"github.com/gin-gonic/gin"
	"music-service/internal/api/handlers"
)

func RegisterGroupRoutes(r *gin.RouterGroup, handler *handlers.GroupHandler) {
	groups := r.Group("/groups")
	{
		groups.POST("", handler.CreateGroup)
		groups.GET("", handler.GetAllGroups)
		groups.GET("/:id", handler.GetGroup)
		groups.PUT("/:id", handler.UpdateGroup)
		groups.DELETE("/:id", handler.DeleteGroup)
	}
}
