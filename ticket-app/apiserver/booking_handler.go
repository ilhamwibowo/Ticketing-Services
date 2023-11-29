//ticket-app/apiserver/handlers
package apiserver

import (
	"net/http"
	"strconv"
	"time"
  	"math/rand"
	"bytes"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
	"github.com/jung-kurt/gofpdf"
	"ticket/storage"
)

func (s *APIServer) defaultRoute(c *gin.Context) {
	c.String(http.StatusOK, "Hello World")
}

func (s *APIServer) holdSeat(c *gin.Context) {
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

	if seat == nil || seat.Status != "OPEN" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Seat is not available"})
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
// send pdf to client

// TODO: impl
// func generatePDF(invoiceID string) ([]byte, error) { 
// 	content := []byte(fmt.Sprintf("PDF Content for Invoice ID: %s", invoiceID))
// 	return content, nil
// }

// NOTE: DELETE LATER
func (s *APIServer) testGeneratePDF(c *gin.Context) {
	eventName := "Sample Event"
	seatNumber := "A1"
	invoiceID := "INV123"
	bookingID := "BK456"
	status := "Paid"

	pdfContent, err := generatePDF(eventName, seatNumber, invoiceID, bookingID, status, "Unexpected Failure")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate PDF"})
		return
	}

	c.Data(http.StatusOK, "application/pdf", pdfContent)
}


func generatePDF(eventName, seatNumber, invoiceID, bookingID, status, failureReason string) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)
	pdf.MultiCell(0, 10, "Invoice", "0", "C", false)

	pdf.SetFont("Arial", "", 12)
	pdf.Ln(20)
	pdf.Cell(0, 10, fmt.Sprintf("Event Name: %s", eventName))
	pdf.Ln(10)
	pdf.Cell(0, 10, fmt.Sprintf("Seat Number: %s", seatNumber))
	pdf.Ln(10)
	pdf.Cell(0, 10, fmt.Sprintf("Invoice ID: %s", invoiceID))
	pdf.Ln(10)
	pdf.Cell(0, 10, fmt.Sprintf("Booking ID: %s", bookingID))
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 10, fmt.Sprintf("Status: %s", strings.ToUpper(status)))

	if strings.ToLower(status) == "failed" {
		pdf.Ln(10)
		pdf.SetFont("Arial", "", 12)
		pdf.Cell(0, 10, fmt.Sprintf("Failure Reason: %s", failureReason))
	}

	qrCode, err := qrcode.New(bookingID, qrcode.Medium)
	if err != nil {
		return nil, err
	}

	qrCode.Content = bookingID

	qrCodeBytes, err := qrCode.PNG(256)
	if err != nil {
		return nil, err
	}

	pdf.RegisterImageOptionsReader("qrCode", gofpdf.ImageOptions{ImageType: "png"}, bytes.NewReader(qrCodeBytes))
	pdf.ImageOptions("qrCode", 10, 170, 40, 0, false, gofpdf.ImageOptions{}, 0, "")

	var buf bytes.Buffer

	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func simulateCall() bool {
	rand.Seed(time.Now().UnixNano())
	randomNum := rand.Intn(100)
	return randomNum > 20
}

// TODO: impl
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
