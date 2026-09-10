package httpapi

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hyanhsing/hyandrop/backend/internal/room"
)

type RoomHandler struct {
	service *room.Service
}

func NewRoomHandler(service *room.Service) *RoomHandler {
	return &RoomHandler{service: service}
}

func (h *RoomHandler) Register(api *gin.RouterGroup) {
	api.POST("/rooms", h.Create)
}

type createRoomRequest struct {
	DeviceID string `json:"device_id"`
}

type createRoomResponse struct {
	RoomID      string    `json:"room_id"`
	ExpiresAt   time.Time `json:"expires_at"`
	DeleteToken string    `json:"delete_token"`
}

func (h *RoomHandler) Create(c *gin.Context) {
	var request createRoomRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "request body must contain valid JSON"})
		return
	}

	result, err := h.service.Create(c.Request.Context(), room.CreateInput{DeviceID: request.DeviceID})
	if err != nil {
		if errors.Is(err, room.ErrInvalidDeviceID) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, room.ErrInvalidRoomTTL) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "room creation is misconfigured"})
			return
		}

		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to create room"})
		return
	}

	c.Header("Cache-Control", "no-store")
	c.JSON(http.StatusCreated, createRoomResponse{
		RoomID:      result.RoomID,
		ExpiresAt:   result.ExpiresAt,
		DeleteToken: result.DeleteToken,
	})
}
