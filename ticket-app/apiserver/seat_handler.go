//ticket-app/apiserver/seat-handler

package apiserver

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"ticket/storage"
)

func (s *APIServer) createSeat(c *gin.Context) {
	var seat storage.Seat

	if err := c.ShouldBindJSON(&seat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if seat.EventID == 0 || seat.SeatNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incomplete data: seat number and event ID are required"})
		return
	}

	event, err := s.storage.GetEventByID(seat.EventID)
	if err != nil || event == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	if err := s.storage.CreateSeat(&seat); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, seat)
}


func (s *APIServer) listSeat(c *gin.Context) {
	seats, err := s.storage.GetSeats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, seats)
}

func (s *APIServer) getSeatStatus(c *gin.Context) {
	eventIDStr := c.Param("event_id")
	seatNumber := c.Param("seat_number")

	eventID, err := strconv.ParseUint(eventIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	seat, err := s.storage.GetSeatByEventIDAndNumber(uint(eventID), seatNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if seat == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Seat not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": seat.Status})
}

