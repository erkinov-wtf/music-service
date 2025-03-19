package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"music-service/internal/api/services"
	"net/http"
	"strconv"
)

type GroupHandler struct {
	groupService *services.GroupService
}

// NewGroupHandler creates a new group handler
func NewGroupHandler(groupService *services.GroupService) *GroupHandler {
	return &GroupHandler{
		groupService: groupService,
	}
}

// CreateGroup godoc
// @Summary Create a new music group
// @Description Create a new music group with the provided name
// @Tags groups
// @Accept json
// @Produce json
// @Param group body object{name=string} true "Group Name"
// @Success 201 {object} object{id=string,name=string,created_at=string,updated_at=string} "Created group data"
// @Failure 400 {object} object{error=string} "Bad request"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /groups [post]
func (h *GroupHandler) CreateGroup(c *gin.Context) {
	var body struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdGroup, err := h.groupService.CreateGroup(c, body.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create group: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdGroup)
}

// GetGroup godoc
// @Summary Get a music group by ID
// @Description Retrieve a music group by its ID
// @Tags groups
// @Produce json
// @Param id path string true "Group ID" format(uuid)
// @Success 200 {object} object{id=string,name=string,created_at=string,updated_at=string}
// @Failure 400 {object} object{error=string} "Bad request"
// @Failure 404 {object} object{error=string} "Group not found"
// @Router /groups/{id} [get]
func (h *GroupHandler) GetGroup(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID format"})
		return
	}

	group, err := h.groupService.GetGroup(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	c.JSON(http.StatusOK, group)
}

// GetAllGroups godoc
// @Summary Get all music groups
// @Description Get a paginated list of music groups
// @Tags groups
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} object{data=array,page=int,limit=int,pages=int,total=int}
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /groups [get]
func (h *GroupHandler) GetAllGroups(c *gin.Context) {
	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	groups, err := h.groupService.GetGroupsWithPagination(c, int32(limit), int32(offset))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve groups: " + err.Error()})
		return
	}

	total, err := h.groupService.GetGroupsCount(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve groups count: " + err.Error()})
		return
	}

	totalPages := (int(total) + limit - 1) / limit

	response := gin.H{
		"data":  groups,
		"page":  page,
		"limit": limit,
		"pages": totalPages,
		"total": total,
	}

	c.JSON(http.StatusOK, response)
}

// UpdateGroup godoc
// @Summary Update a music group
// @Description Update a music group's information
// @Tags groups
// @Accept json
// @Produce json
// @Param id path string true "Group ID" format(uuid)
// @Param group body object{name=string} true "Group Info"
// @Success 200 {object} object{id=string,name=string,created_at=string,updated_at=string} "Group updated successfully"
// @Failure 400 {object} object{error=string} "Bad request"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /groups/{id} [put]
func (h *GroupHandler) UpdateGroup(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID format"})
		return
	}

	var body struct {
		Name string `json:"name" binding:"required"`
	}

	if err = c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	group, err := h.groupService.UpdateGroup(c, id, body.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update group: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": group})
}

// DeleteGroup godoc
// @Summary Delete a music group
// @Description Delete a music group by ID
// @Tags groups
// @Produce json
// @Param id path string true "Group ID" format(uuid)
// @Success 204 {object} object{message=string} "Group deleted successfully"
// @Failure 400 {object} object{error=string} "Bad request"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /groups/{id} [delete]
func (h *GroupHandler) DeleteGroup(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID format"})
		return
	}

	if err := h.groupService.DeleteGroup(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete group: " + err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "Group deleted successfully"})
}
