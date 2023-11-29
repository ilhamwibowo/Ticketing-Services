package apiserver

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (s *APIServer) getAllEvents(c *gin.Context) {
    events, err := s.storage.GetEvents()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, events)
}

func (s *APIServer) getEmptySeats(c *gin.Context) {
	eventIDStr := c.Param("event_id")
	eventID, err := strconv.ParseUint(eventIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	seats, err := s.storage.GetEmptySeatsByEventID(uint(eventID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, seats)
}
