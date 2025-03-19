package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"music-service/internal/api/services"
	"music-service/internal/storage/database"
	"music-service/internal/storage/database/repository"
	"net/http"
	"strconv"
	"time"
)

type SongHandler struct {
	songService *services.SongService
}

func NewSongHandler(songService *services.SongService) *SongHandler {
	return &SongHandler{
		songService: songService,
	}
}

func (h *SongHandler) CreateSong(c *gin.Context) {
	var body struct {
		GroupID     string `json:"group_id" binding:"required"`
		Title       string `json:"title" binding:"required"`
		Runtime     int32  `json:"runtime" binding:"required"`
		Lyrics      []byte `json:"lyrics"`
		ReleaseDate string `json:"release_date" binding:"required"`
		Link        string `json:"link" binding:"required"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	groupID, err := uuid.Parse(body.GroupID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID format"})
		return
	}

	releaseDate, err := time.Parse(time.RFC3339, body.ReleaseDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid release date format"})
		return
	}

	params := repository.SongCreateParams{
		GroupID:     groupID,
		Title:       body.Title,
		Runtime:     body.Runtime,
		Lyrics:      body.Lyrics,
		ReleaseDate: releaseDate,
		Link:        body.Link,
	}

	if err := h.songService.CreateSong(c, params); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create song: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Song created successfully"})
}

func (h *SongHandler) GetSong(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID format"})
		return
	}

	song, err := h.songService.GetSong(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
		return
	}

	c.JSON(http.StatusOK, song)
}

func (h *SongHandler) GetAllSongs(c *gin.Context) {
	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	groupName := c.Query("group")
	songTitle := c.Query("song")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	var songs []database.GetSongsWithPaginationRow
	var total int64

	if groupName != "" || songTitle != "" {
		params := repository.SongFilterParams{
			Limit:     int32(limit),
			Offset:    int32(offset),
			GroupName: groupName,
			SongTitle: songTitle,
		}

		songs, err = h.songService.GetSongsWithFilters(c, params)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve songs: " + err.Error()})
			return
		}

		total, err = h.songService.GetSongsCountWithFilters(c, groupName, songTitle)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve songs count: " + err.Error()})
			return
		}
	} else {
		songs, err = h.songService.GetSongsWithPagination(c, int32(limit), int32(offset))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve songs: " + err.Error()})
			return
		}

		total, err = h.songService.GetSongsCount(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve songs count: " + err.Error()})
			return
		}
	}

	totalPages := (int(total) + limit - 1) / limit

	response := gin.H{
		"data":  songs,
		"page":  page,
		"limit": limit,
		"pages": totalPages,
		"total": total,
	}

	c.JSON(http.StatusOK, response)
}

func (h *SongHandler) UpdateSong(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID format"})
		return
	}

	var body struct {
		GroupID     string `json:"group_id" binding:"required"`
		Title       string `json:"title" binding:"required"`
		Runtime     int32  `json:"runtime" binding:"required"`
		Lyrics      []byte `json:"lyrics"`
		ReleaseDate string `json:"release_date" binding:"required"`
		Link        string `json:"link" binding:"required"`
	}

	if err = c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	groupID, err := uuid.Parse(body.GroupID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID format"})
		return
	}

	releaseDate, err := time.Parse(time.RFC3339, body.ReleaseDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid release date format"})
		return
	}

	params := repository.SongUpdateParams{
		ID:          id,
		GroupID:     groupID,
		Title:       body.Title,
		Runtime:     body.Runtime,
		Lyrics:      body.Lyrics,
		ReleaseDate: releaseDate,
		Link:        body.Link,
	}

	if err := h.songService.UpdateSong(c, params); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update song: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Song updated successfully"})
}

func (h *SongHandler) DeleteSong(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID format"})
		return
	}

	if err := h.songService.DeleteSong(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete song: " + err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "Song deleted successfully"})
}
