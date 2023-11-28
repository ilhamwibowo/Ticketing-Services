package apiserver

import (
	"net/http"
	"strconv"
	"time"
  	"math/rand"
	"bytes"
	"fmt"
	
	"github.com/gin-gonic/gin"
	"ticket/storage"
)

func (s *APIServer) defaultRoute(c *gin.Context) {
	c.String(http.StatusOK, "Hello World")
}

func (s *APIServer) createSeat(c *gin.Context) {
	var seat storage.Seat

	if err := c.ShouldBindJSON(&seat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

func (s *APIServer) holdSeat(c *gin.Context) {
	seatIDStr := c.Param("id")

  seatID, err := strconv.ParseUint(seatIDStr, 10, 64)
	seat, err := s.storage.GetSeatByID(uint(seatID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if seat.Status != "OPEN" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Seat already booked"})
		return
	}

	success := simulateCall()

	if !success {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Booking failed"})
		return
	}

	invoiceID, paymentURL, paymentSuccess := callPaymentAPI()

	if !paymentSuccess {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Payment failed"})
		return
	}

	// Create a new booking
	booking := &storage.Booking{
		SeatID:     seat.ID,
		InvoiceID:  invoiceID,
		PaymentURL: paymentURL,
		Status:     "ON GOING",
	}

	if err := s.storage.CreateBooking(booking); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update seat status to 'BOOKED'
	seat.Status = "ON GOING"
	if err := s.storage.UpdateSeat(seat); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return booking ongoing
	c.JSON(http.StatusOK, gin.H{
		"message":    "Booking ongoing",
		"invoice_id": invoiceID,
		"payment_url": paymentURL,
	})
}

func (s *APIServer) paymentWebhook(c *gin.Context) {
	// Parse the incoming JSON payload from the payment app
	var webhookData struct {
		InvoiceID string `json:"invoice_id"`
		Status    string `json:"status"`
	}

	if err := c.BindJSON(&webhookData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	// Update the status of the invoiceID in the database
	err := s.storage.UpdateBookingStatusByInvoiceID(webhookData.InvoiceID, webhookData.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update status"})
		return
	}

  // send pdf
  // Generate PDF content
  // pdfContent, err := generatePDF(webhookData.InvoiceID)
  // if err != nil {
  //   c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate PDF"})
  //   return
  // }

  
  // apiURL := "https://your-ticket-app.com/send-pdf"

  // resp, err := sendPDFToTicketApp(apiURL, pdfContent)
  // if err != nil {
  //   c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send PDF"})
  //   return
  // }
  // defer resp.Body.Close()

  // if resp.StatusCode != http.StatusOK {
  //   c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send PDF to ticket app"})
  //   return
  // }

  c.JSON(http.StatusOK, gin.H{"message": "PDF sent to ticket app"})
}

// TODO: impl
func sendPDFToTicketApp(apiURL string, pdfContent []byte) (*http.Response, error) {
	// Create a POST request to the specified API endpoint
	req, err := http.NewRequest(http.MethodPost, apiURL, bytes.NewBuffer(pdfContent))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/pdf")

	// Perform the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// TODO: impl
func generatePDF(invoiceID string) ([]byte, error) { 
	content := []byte(fmt.Sprintf("PDF Content for Invoice ID: %s", invoiceID))
	return content, nil
}

func simulateCall() bool {
	rand.Seed(time.Now().UnixNano())
	randomNum := rand.Intn(100)
	return randomNum > 20
}

func callPaymentAPI() (string, string, bool) {
	rand.Seed(time.Now().UnixNano())
	randomNum := rand.Intn(100)
	if randomNum < 10 {
		return "", "", false // Simulate failure
	}

	mockInvoiceID := "INV12345" // TODO: change
	mockPaymentURL := "https://payment.app.com/pay/INV12345"
	return mockInvoiceID, mockPaymentURL, true
}
