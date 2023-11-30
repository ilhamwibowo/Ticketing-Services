// ticket-app/apiserver/booking_handler.go
package apiserver

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"time"

	"ticket/storage"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
	"github.com/skip2/go-qrcode"
)

func (s *APIServer) defaultRoute(c *gin.Context) {
	c.String(http.StatusOK, "Hello World")
}

func (s *APIServer) checkClientHealth(c *gin.Context) {
	resp, err := http.Get("http://web:8000/payment/hello-world")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to the service"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		c.JSON(http.StatusOK, gin.H{"status": "Service is healthy"})
	} else {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "Service is not healthy", "code": resp.StatusCode})
	}
}

func (s *APIServer) holdSeat(c *gin.Context) {
	eventIDStr := c.Param("event_id")
	seatNumber := c.Param("seat_number")

	eventID, err := strconv.ParseUint(eventIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	event, err := s.storage.GetEventByID(uint(eventID))
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
	fmt.Println("External Call")
	if !success {
		paymentURL := "https://web:8000/payment/process-payment"

		pdfContent, err := generatePDF(event.EventName, seat.SeatNumber, "-1", "0", "False", "Payment Failed")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate PDF"})
			return
		}

		// Encode PDF content to base64 for inclusion in the response
		pdfBase64 := base64.StdEncoding.EncodeToString(pdfContent)

		c.JSON(http.StatusBadRequest, gin.H{
			"message":     "Booking failed",
			"invoice_id":  "-1",
			"payment_url": paymentURL,
			"pdf":         pdfBase64,
			"error":       "Payment failed",
		})
		return
	}

	invoiceID, paymentURL, err := callPaymentAPI()
	fmt.Println("Payment Done")

	if err != nil {
		pdfContent, err := generatePDF(event.EventName, seat.SeatNumber, invoiceID, "0", "Failed", "Payment Failed")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate PDF"})
			return
		}

		pdfBase64 := base64.StdEncoding.EncodeToString(pdfContent)

		c.JSON(http.StatusBadRequest, gin.H{
			"message":     "Booking failed",
			"invoice_id":  invoiceID,
			"payment_url": paymentURL,
			"pdf":         pdfBase64,
			"error":       "Payment failed",
		})
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

	seat.Status = "ON GOING"
	if err := s.storage.UpdateSeat(seat); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Booking ongoing",
		"invoice_id":  invoiceID,
		"payment_url": paymentURL,
	})
}

func (s *APIServer) paymentWebhook(c *gin.Context) {
	var webhookData struct {
		InvoiceID string `json:"invoice_id"`
		Status    string `json:"status"`
	}

	if err := c.BindJSON(&webhookData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	err := s.storage.UpdateBookingStatusByInvoiceID(webhookData.InvoiceID, webhookData.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update status"})
		return
	}

	booking, err := s.storage.GetBookingByInvoiceID(webhookData.InvoiceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve booking"})
		return
	}

	seat, err := s.storage.GetSeatByID(booking.SeatID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve seat"})
		return
	}

	if webhookData.Status == "True" {
		err := s.storage.UpdateSeatStatusByID(seat.ID, "BOOKED")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update seat status"})
			return
		} else {
			err := s.storage.UpdateSeatStatusByID(seat.ID, "OPEN")
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update seat status"})
				return
			}

		}
	}

	event, err := s.storage.GetEventByID(seat.EventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
		return
	}

	pdfContent, err := generatePDF(event.EventName, seat.SeatNumber, booking.InvoiceID, fmt.Sprintf("%d", booking.ID), booking.Status, "Unexpected Error")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate PDF"})
		return
	}

	err = sendPDFToClient(booking.InvoiceID, pdfContent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "PDF sent to ticket app"})
}

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

func callPaymentAPI() (string, string, error) {
	resp, err := http.Get("http://web:8000/payment/process-payment/")
	fmt.Println("OKI")

	if err != nil {
		fmt.Println(err)
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("API call failed with status code: %d", resp.StatusCode)
	}

	var paymentData struct {
		InvoiceID  string `json:"invoice_id"`
		PaymentURL string `json:"payment_url"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&paymentData); err != nil {
		return "", "", err
	}
	return paymentData.InvoiceID, paymentData.PaymentURL, nil
}

func sendPDFToClient(invoiceID string, pdfContent []byte) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	_ = writer.WriteField("invoice_id", invoiceID)

	part, err := writer.CreateFormFile("invoice_pdf", "invoice.pdf")
	if err != nil {
		return err
	}
	_, _ = part.Write(pdfContent)

	_ = writer.Close()

	apiURL := "http://client-web:8000/book/api/invoices/create/"

	req, err := http.NewRequest("POST", apiURL, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to upload invoice: %d", resp.StatusCode)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println("API Response:", string(respBody))

	return nil
}

func generateRandomInvoiceID() string {
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	rand.Seed(time.Now().UnixNano())

	invoiceID := make([]rune, 8)
	for i := range invoiceID {
		invoiceID[i] = chars[rand.Intn(len(chars))]
	}

	return string(invoiceID)
}
