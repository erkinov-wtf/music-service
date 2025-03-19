package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"music-service/internal/api/services"
	"music-service/internal/pkg/utils/constants"
	"music-service/internal/pkg/utils/parser"
	"music-service/internal/storage/database"
	"music-service/internal/storage/database/repository"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type SongHandler struct {
	songService  *services.SongService
	groupService *services.GroupService
}

func NewSongHandler(songService *services.SongService, groupService *services.GroupService) *SongHandler {
	return &SongHandler{
		songService:  songService,
		groupService: groupService,
	}
}

// GroupData represents group information to be included in song responses
type GroupData struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SongResponse is the formatted song response for the API
type SongResponse struct {
	ID          string    `json:"id"`
	Group       GroupData `json:"group"` // Changed from GroupID to Group
	Title       string    `json:"title"`
	Runtime     int32     `json:"runtime"`
	Lyrics      string    `json:"lyrics"`
	ReleaseDate time.Time `json:"release_date"`
	Link        string    `json:"link"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateSong godoc
// @Summary Create a new song
// @Description Create a new song with the provided details and return the created song data
// @Tags songs
// @Accept json
// @Produce json
// @Param song body object{group_id=string,title=string,runtime=integer,lyrics=string,release_date=string,link=string} true "Song Information"
// @Success 201 {object} object{data=object{id=string,group=object{id=string,name=string,created_at=string,updated_at=string},title=string,runtime=integer,lyrics=string,release_date=string,link=string,created_at=string,updated_at=string}} "Created song data"
// @Failure 400 {object} object{error=string} "Bad request - Invalid input data"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /songs [post]
func (h *SongHandler) CreateSong(c *gin.Context) {
	var body struct {
		GroupID     string `json:"group_id" binding:"required"`
		Title       string `json:"title" binding:"required"`
		Runtime     int32  `json:"runtime" binding:"required"`
		Lyrics      string `json:"lyrics"`
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

	releaseDate, err := time.Parse(constants.DateFormat, body.ReleaseDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid release date format"})
		return
	}

	lyricsJSON, err := parser.ParseLyrics(body.Lyrics)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process lyrics"})
		return
	}

	params := repository.SongCreateParams{
		GroupID:     groupID,
		Title:       body.Title,
		Runtime:     body.Runtime,
		Lyrics:      lyricsJSON,
		ReleaseDate: releaseDate,
		Link:        body.Link,
	}

	song, err := h.songService.CreateSong(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create song: " + err.Error()})
		return
	}

	response, err := h.formatSong(c, song)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve song: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": response})
}

// GetSong godoc
// @Summary Get a song by ID
// @Description Retrieve a song by its ID
// @Tags songs
// @Produce json
// @Param id path string true "Song ID" format(uuid)
// @Success 200 {object} object{id=string,group=object{id=string,name=string,created_at=string,updated_at=string},title=string,runtime=integer,lyrics=string,release_date=string,link=string,created_at=string,updated_at=string}
// @Failure 400 {object} object{error=string} "Bad request"
// @Failure 404 {object} object{error=string} "Song not found"
// @Router /songs/{id} [get]
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

	response, err := h.formatSong(c, song)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve song: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetAllSongs godoc
// @Summary Get all songs with pagination and filtering
// @Description Get a paginated list of songs with optional filtering by group name and song title
// @Tags songs
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param group query string false "Filter by group name"
// @Param song query string false "Filter by song title"
// @Success 200 {object} object{data=array,page=int,limit=int,pages=int,total=int}
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /songs [get]
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

	bulkSongs, err := h.formatBulkSongs(c, songs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to formatting songs: " + err.Error()})
		return
	}

	response := gin.H{
		"data":  bulkSongs,
		"page":  page,
		"limit": limit,
		"pages": totalPages,
		"total": total,
	}

	c.JSON(http.StatusOK, response)
}

// GetSongVerses godoc
// @Summary Get song verses with pagination
// @Description Get a song's lyrics split by verses (lines) with pagination
// @Tags songs
// @Produce json
// @Param id path string true "Song ID" format(uuid)
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} object{song_id=string,page=int,limit=int,pages=int,total=int,verses=array} "Paginated verses"
// @Failure 400 {object} object{error=string} "Bad request"
// @Failure 404 {object} object{error=string} "Song not found"
// @Router /songs/{id}/verses [get]
func (h *SongHandler) GetSongVerses(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID format"})
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	song, err := h.songService.GetSong(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
		return
	}

	var lyricsData struct {
		Text   string   `json:"text"`
		Verses []string `json:"verses"`
	}

	if err = json.Unmarshal(song.Lyrics, &lyricsData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse lyrics"})
		return
	}

	total := len(lyricsData.Verses)
	pages := (total + limit - 1) / limit

	startIndex := (page - 1) * limit
	endIndex := startIndex + limit

	if startIndex >= total {
		c.JSON(http.StatusOK, gin.H{
			"song_id": song.ID.String(),
			"page":    page,
			"limit":   limit,
			"pages":   pages,
			"total":   total,
			"verses":  []string{},
		})
		return
	}

	if endIndex > total {
		endIndex = total
	}

	verses := lyricsData.Verses[startIndex:endIndex]

	c.JSON(http.StatusOK, gin.H{
		"song_id": song.ID.String(),
		"page":    page,
		"limit":   limit,
		"pages":   pages,
		"total":   total,
		"verses":  verses,
	})
}

// UpdateSong godoc
// @Summary Update a song
// @Description Update an existing song's information by ID and return the updated song data
// @Tags songs
// @Accept json
// @Produce json
// @Param id path string true "Song ID" format(uuid)
// @Param song body object{group_id=string,title=string,runtime=integer,lyrics=string,release_date=string,link=string} true "Song Information"
// @Success 200 {object} object{message=object{id=string,group=object{id=string,name=string,created_at=string,updated_at=string},title=string,runtime=integer,lyrics=string,release_date=string,link=string,created_at=string,updated_at=string}} "Updated song data"
// @Failure 400 {object} object{error=string} "Bad request - Invalid input or ID"
// @Failure 404 {object} object{error=string} "Song not found"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /songs/{id} [put]
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
		Lyrics      string `json:"lyrics"`
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

	// If lyrics were provided, parse them
	var lyricsJSON []byte
	if body.Lyrics != "" {
		lyricsJSON, err = parser.ParseLyrics(body.Lyrics)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process lyrics"})
			return
		}
	} else {
		// If no lyrics were provided, get the existing lyrics
		song, err := h.songService.GetSong(c, id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
			return
		}
		lyricsJSON = song.Lyrics
	}

	params := repository.SongUpdateParams{
		ID:          id,
		GroupID:     groupID,
		Title:       body.Title,
		Runtime:     body.Runtime,
		Lyrics:      lyricsJSON,
		ReleaseDate: releaseDate,
		Link:        body.Link,
	}

	song, err := h.songService.UpdateSong(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update song: " + err.Error()})
		return
	}

	response, err := h.formatSong(c, song)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve song: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": response})
}

// DeleteSong godoc
// @Summary Delete a song
// @Description Delete a song by ID
// @Tags songs
// @Produce json
// @Param id path string true "Song ID" format(uuid)
// @Success 204 {object} object{message=string} "Song deleted successfully"
// @Failure 400 {object} object{error=string} "Bad request"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /songs/{id} [delete]
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

// Format a single song with group data
func (h *SongHandler) formatSong(c *gin.Context, song database.Song) (SongResponse, error) {
	var lyricsData struct {
		Text   string   `json:"text"`
		Verses []string `json:"verses"`
	}

	if err := json.Unmarshal(song.Lyrics, &lyricsData); err != nil {
		return SongResponse{}, err
	}

	lyrics := lyricsData.Text
	if lyrics == "" && len(lyricsData.Verses) > 0 {
		lyrics = strings.Join(lyricsData.Verses, "\n")
	}

	groupId, err := uuid.Parse(song.GroupID.String())
	if err != nil {
		return SongResponse{}, err
	}

	group, err := h.groupService.GetGroup(c, groupId)
	if err != nil {
		return SongResponse{}, err
	}

	return SongResponse{
		ID: song.ID.String(),
		Group: GroupData{
			ID:        group.ID.String(),
			Name:      group.Name,
			CreatedAt: group.CreatedAt.Time,
			UpdatedAt: group.UpdatedAt.Time,
		},
		Title:       song.Title,
		Runtime:     song.Runtime,
		Lyrics:      lyrics,
		ReleaseDate: song.ReleaseDate.Time,
		Link:        song.Link,
		CreatedAt:   song.CreatedAt.Time,
		UpdatedAt:   song.UpdatedAt.Time,
	}, nil
}

// Format multiple songs with group data
func (h *SongHandler) formatBulkSongs(c *gin.Context, songs []database.GetSongsWithPaginationRow) ([]SongResponse, error) {
	var formattedSongs []SongResponse

	groupCache := make(map[string]database.Group)

	for _, song := range songs {
		var lyricsData struct {
			Text   string   `json:"text"`
			Verses []string `json:"verses"`
		}

		if err := json.Unmarshal(song.Lyrics, &lyricsData); err != nil {
			return nil, err
		}

		lyrics := lyricsData.Text
		if lyrics == "" && len(lyricsData.Verses) > 0 {
			lyrics = strings.Join(lyricsData.Verses, "\n")
		}

		groupID := song.GroupID.String()
		var group database.Group
		var ok bool

		if group, ok = groupCache[groupID]; !ok {
			groupId, err := uuid.Parse(song.GroupID.String())
			if err != nil {
				return nil, err
			}
			group, err = h.groupService.GetGroup(c, groupId)
			if err != nil {
				return nil, err
			}
			groupCache[groupID] = group
		}

		formattedSong := SongResponse{
			ID: song.ID.String(),
			Group: GroupData{
				ID:        group.ID.String(),
				Name:      group.Name,
				CreatedAt: group.CreatedAt.Time,
				UpdatedAt: group.UpdatedAt.Time,
			},
			Title:       song.Title,
			Runtime:     song.Runtime,
			Lyrics:      lyrics,
			ReleaseDate: song.ReleaseDate.Time,
			Link:        song.Link,
			CreatedAt:   song.CreatedAt.Time,
			UpdatedAt:   song.UpdatedAt.Time,
		}

		formattedSongs = append(formattedSongs, formattedSong)
	}

	return formattedSongs, nil
}
